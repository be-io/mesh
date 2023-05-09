/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	grpcproto "google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var dialer = &net.Dialer{
	KeepAlive: 30 * time.Second,
}

var DefaultTransport http.RoundTripper = &http.Transport{
	DialContext:           dialer.DialContext,
	ForceAttemptHTTP2:     false,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	uri, err := url.Parse(target)
	if nil != err {
		return nil, cause.Error(err)
	}

	return &Channel{Transport: DefaultTransport, BaseURL: uri}, nil
}

var _ grpc.ClientConnInterface = (*Channel)(nil)

var grpcDetailsHeader = textproto.CanonicalMIMEHeaderKey("X-GRPC-Details")

type Channel struct {
	Transport http.RoundTripper
	BaseURL   *url.URL
}

func (that *Channel) Invoke(ctx context.Context, methodName string, req, resp interface{}, opts ...grpc.CallOption) error {
	h := headersFromContext(ctx)
	prsim.SetMetadata(mpc.ContextWith(ctx), h)
	h.Set("Content-Type", UnaryContentTypeV1)

	copts := GetCallOptions(opts)

	uri := *that.BaseURL
	uri.Path = path.Join(uri.Path, methodName)
	uriStr := uri.String()
	ctx, err := ApplyPerRPCCreds(ctx, copts, uriStr, uri.Scheme == "https")
	if nil != err {
		return err
	}

	codec := encoding.GetCodec(grpcproto.Name)
	b, err := codec.Marshal(req)
	if nil != err {
		return err
	}

	// TODO: enforce max send and receive size in call options

	r, err := http.NewRequest("POST", uriStr, bytes.NewReader(b))
	if nil != err {
		return err
	}
	r.Header = h
	reply, err := that.Transport.RoundTrip(r.WithContext(ctx))
	if nil != err {
		return statusFromContextError(err)
	}

	// we fire up a goroutine to read the response so that we can properly
	// respect any context deadline (e.g. don't want to be blocked, reading
	// from socket, long past requested timeout).
	respCh := make(chan struct{})
	go func() {
		defer close(respCh)
		b, err = io.ReadAll(reply.Body)
		log.Catch(reply.Body.Close())
	}()

	if len(copts.Peer) > 0 {
		copts.SetPeer(getPeer(that.BaseURL, r.TLS))
	}

	// gather headers and trailers
	if len(copts.Headers) > 0 || len(copts.Trailers) > 0 {
		if err := setMetadata(reply.Header, copts); nil != err {
			return err
		}
	}

	if stat := statFromResponse(reply); stat.Code() != codes.OK {
		return stat.Err()
	}

	select {
	case <-ctx.Done():
		return statusFromContextError(ctx.Err())
	case <-respCh:
	}
	if nil != err {
		return err
	}
	return codec.Unmarshal(b, resp)
}

func (that *Channel) NewStream(ctx context.Context, desc *grpc.StreamDesc, methodName string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	h := headersFromContext(ctx)
	prsim.SetMetadata(mpc.ContextWith(ctx), h)
	h.Set("Content-Type", StreamContentTypeV1)

	copts := GetCallOptions(opts)

	uri := *that.BaseURL
	uri.Path = path.Join(uri.Path, methodName)
	uriStr := uri.String()
	ctx, err := ApplyPerRPCCreds(ctx, copts, uriStr, uri.Scheme == "https")
	if nil != err {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	// Intercept r.Close() so we can control the error sent across to the writer thread.
	r, w := io.Pipe()
	req, err := http.NewRequest("POST", uriStr, io.NopCloser(r))
	if nil != err {
		cancel()
		return nil, err
	}
	req.Header = h

	cs := newClientStream(ctx, cancel, w, desc.ServerStreams, copts, that.BaseURL)
	go cs.doHttpCall(that.Transport, req, r)

	// ensure that context is cancelled, even if caller
	// fails to fully consume or cancel the stream
	ret := &clientStreamWrapper{cs}
	runtime.SetFinalizer(ret, func(*clientStreamWrapper) { cancel() })

	return ret, nil
}

type clientStreamWrapper struct {
	grpc.ClientStream
}

func getPeer(baseUrl *url.URL, tls *tls.ConnectionState) *peer.Peer {
	hostPort := baseUrl.Host
	if !strings.Contains(hostPort, ":") {
		if baseUrl.Scheme == "https" {
			hostPort = hostPort + ":443"
		} else if baseUrl.Scheme == "http" {
			hostPort = hostPort + ":80"
		}
	}
	pr := peer.Peer{Addr: strAddr(hostPort)}
	if tls != nil {
		pr.AuthInfo = credentials.TLSInfo{State: *tls}
	}
	return &pr
}

func setMetadata(h http.Header, copts *CallOptions) error {
	hdr, err := asMetadata(h)
	if nil != err {
		return err
	}
	tlr := metadata.MD{}

	const trailerPrefix = "x-grpc-trailer-"

	for k, v := range hdr {
		if strings.HasPrefix(strings.ToLower(k), trailerPrefix) {
			trailerName := k[len(trailerPrefix):]
			if trailerName != "" {
				tlr[trailerName] = v
				delete(hdr, k)
			}
		}
	}

	copts.SetHeaders(hdr)
	copts.SetTrailers(tlr)
	return nil
}

type clientStream struct {
	ctx     context.Context
	cancel  context.CancelFunc
	copts   *CallOptions
	baseUrl *url.URL
	codec   encoding.Codec

	// respStream is set to indicate whether client expects stream response; unary if false
	respStream bool

	// hd and hdErr are populated when ready is done
	ready sync.WaitGroup
	hdErr error
	hd    metadata.MD

	// rCh is used to deliver messages from doHttpCall goroutine
	// to callers of RecvMsg.
	// done must be set to true before it is closed
	rCh chan []byte

	// rMu protects done, rErr, and tr
	rMu  sync.RWMutex
	done bool
	rErr error
	tr   HttpTrailer

	// wMu protects w and wErr
	wMu  sync.Mutex
	w    io.WriteCloser
	wErr error
}

func newClientStream(ctx context.Context, cancel context.CancelFunc, w io.WriteCloser, recvStream bool, copts *CallOptions, baseUrl *url.URL) *clientStream {
	cs := &clientStream{
		ctx:        ctx,
		cancel:     cancel,
		copts:      copts,
		baseUrl:    baseUrl,
		codec:      getStreamingCodec(StreamContentTypeV1),
		w:          w,
		respStream: recvStream,
		rCh:        make(chan []byte),
	}
	cs.ready.Add(1)
	return cs
}

func (that *clientStream) Header() (metadata.MD, error) {
	that.ready.Wait()
	return that.hd, that.hdErr
}

func (that *clientStream) Trailer() metadata.MD {
	// only safe to read trailers after stream has completed
	that.rMu.RLock()
	defer that.rMu.RUnlock()
	if that.done {
		return metadataFromProto(that.tr.Metadata)
	}
	return nil
}

func metadataFromProto(trailers map[string]*TrailerValues) metadata.MD {
	md := metadata.MD{}
	for k, vs := range trailers {
		md[k] = vs.Values
	}
	return md
}

func (that *clientStream) CloseSend() error {
	that.wMu.Lock()
	defer that.wMu.Unlock()
	return that.w.Close()
}

func (that *clientStream) Context() context.Context {
	return that.ctx
}

func (that *clientStream) readErrorIfDone() (bool, error) {
	that.rMu.RLock()
	defer that.rMu.RUnlock()
	if !that.done {
		return false, nil
	}
	if that.rErr != nil {
		return true, that.rErr
	}
	if that.tr.Code == int32(codes.OK) {
		return true, io.EOF
	}
	statProto := spb.Status{
		Code:    that.tr.Code,
		Message: that.tr.Message,
		Details: that.tr.Details,
	}
	return true, status.FromProto(&statProto).Err()
}

func (that *clientStream) SendMsg(m interface{}) error {
	// GRPC streams return EOF error for attempts to send on closed stream
	if done, _ := that.readErrorIfDone(); done {
		return io.EOF
	}

	that.wMu.Lock()
	defer that.wMu.Unlock()
	if that.wErr != nil {
		// earlier write error means stream is effectively closed
		return io.EOF
	}

	that.wErr = writeMessage(that.w, that.codec, m, false)
	return that.wErr
}

func (that *clientStream) RecvMsg(m interface{}) error {
	if done, err := that.readErrorIfDone(); done {
		return err
	}

	select {
	case <-that.ctx.Done():
		return statusFromContextError(that.ctx.Err())
	case msg, ok := <-that.rCh:
		if !ok {
			done, err := that.readErrorIfDone()
			if !done {
				// sanity check: this shouldn't be possible
				panic("cs.rCh was closed but cs.done == false!")
			}
			return err
		}
		err := that.codec.Unmarshal(msg, m)
		if nil != err {
			return status.Error(codes.Internal, fmt.Sprintf("server sent invalid message: %v", err))
		}
		if !that.respStream {
			// We need to query the channel for a second message. If there *is* a
			// second message, the server tried to send too many, and that's an
			// error. And if there isn't a second message, we still need to see the
			// channel close (e.g. end-of-stream) so we know that tr is set (so that
			// it's available for a subsequent call to Trailer)
			select {
			case <-that.ctx.Done():
				return statusFromContextError(that.ctx.Err())
			case _, ok := <-that.rCh:
				if ok {
					// server tried to send >1 message!
					that.rMu.Lock()
					defer that.rMu.Unlock()
					if that.rErr == nil {
						that.rErr = status.Error(codes.Internal, "method should return 1 response message but server sent >1")
						that.done = true
						// we won't be reading from the channel anymore, so we must
						// cancel the context so that doHttpCall doesn't hang trying
						// to write to channel
						that.cancel()
					}
					return that.rErr
				}
				// if server sent a failure after the single message, the failure takes precedence
				done, err := that.readErrorIfDone()
				if !done {
					// sanity check: this shouldn't be possible
					panic("cs.rCh was closed but cs.done == false!")
				}
				if err != io.EOF {
					return err
				}
			}
		}
		return nil
	}
}

func (that *clientStream) doHttpCall(transport http.RoundTripper, req *http.Request, readPipe *io.PipeReader) {
	// On completion, we must fill in cs.tr or cs.rErr and then close channel,
	// which signals to client code that we've reached end-of-stream.

	var rErr error
	rMuHeld := false

	defer func() {
		if !rMuHeld {
			that.rMu.Lock()
		}
		defer that.rMu.Unlock()

		if rErr != nil && that.rErr == nil {
			that.rErr = rErr
		}
		that.done = true
		log.Catch(readPipe.CloseWithError(rErr))
		close(that.rCh)
	}()

	onReady := func(err error, headers metadata.MD) {
		that.hdErr = err
		that.hd = headers
		if len(headers) > 0 && len(that.copts.Headers) > 0 {
			that.copts.SetHeaders(headers)
		}
		rErr = err
		that.ready.Done()
	}

	reply, err := transport.RoundTrip(req.WithContext(that.ctx))
	if nil != err {
		onReady(statusFromContextError(err), nil)
		return
	}
	defer func() {
		_, ex := io.ReadAll(reply.Body)
		log.Catch(ex)
		log.Catch(reply.Body.Close())
	}()

	if len(that.copts.Peer) > 0 {
		that.copts.SetPeer(getPeer(that.baseUrl, reply.TLS))
	}
	md, err := asMetadata(reply.Header)
	if nil != err {
		onReady(err, nil)
		return
	}

	onReady(nil, md)

	stat := statFromResponse(reply)
	if stat.Code() != codes.OK {
		statProto := stat.Proto()
		that.tr.Code = statProto.Code
		that.tr.Message = statProto.Message
		that.tr.Details = statProto.Details
		return
	}

	counter := 0
	for {
		// TODO: enforce max send and receive size in call options

		counter++
		var sz int32
		sz, rErr = readSizePreface(reply.Body)
		if rErr != nil {
			return
		}
		if sz < 0 {
			// final message is a trailer (need lock to write to cs.tr)
			that.rMu.Lock()
			rMuHeld = true // defer above will unlock for us
			that.rErr = readMessage(reply.Body, that.codec, int32(-sz), &that.tr)
			if that.rErr != nil {
				if that.rErr == io.EOF {
					that.rErr = io.ErrUnexpectedEOF
				}
			}
			if len(that.tr.Metadata) > 0 && len(that.copts.Trailers) > 0 {
				that.copts.SetTrailers(metadataFromProto(that.tr.Metadata))
			}
			return
		}
		msg := make([]byte, sz)
		_, rErr = io.ReadAtLeast(reply.Body, msg, int(sz))
		if rErr != nil {
			if rErr == io.EOF {
				rErr = io.ErrUnexpectedEOF
			}
			return
		}

		select {
		case <-that.ctx.Done():
			// operation timed out or was cancelled before we could
			// successfully send this message to client code
			rErr = statusFromContextError(that.ctx.Err())
			return
		case that.rCh <- msg:
		}
	}
}

func statusFromContextError(err error) error {
	if err == context.DeadlineExceeded {
		return status.Error(codes.DeadlineExceeded, err.Error())
	} else if err == context.Canceled {
		return status.Error(codes.Canceled, err.Error())
	}
	return err
}

func headersFromContext(ctx context.Context) http.Header {
	h := http.Header{}
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		toHeaders(md, h, "")
	}
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		millis := int64(timeout / time.Millisecond)
		if millis <= 0 {
			millis = 1
		}
		h.Set("GRPC-Timeout", fmt.Sprintf("%dm", millis))
	}
	return h
}

func statFromResponse(reply *http.Response) *status.Status {
	code := codeFromHttpStatus(reply.StatusCode)
	msg := reply.Status
	codeStrs := strings.SplitN(reply.Header.Get("X-GRPC-Status"), ":", 2)
	if len(codeStrs) > 0 && codeStrs[0] != "" {
		if c, err := strconv.ParseInt(codeStrs[0], 10, 32); err == nil {
			code = codes.Code(c)
		}
		if len(codeStrs) > 1 {
			msg = codeStrs[1]
		}
	}
	if code != codes.OK {
		var details []*anypb.Any
		if detailHeaders := reply.Header[grpcDetailsHeader]; len(detailHeaders) > 0 {
			details = make([]*anypb.Any, 0, len(detailHeaders))
			for _, d := range detailHeaders {
				b, err := base64.RawURLEncoding.DecodeString(d)
				if nil != err {
					continue
				}
				var msg anypb.Any
				if err := proto.Unmarshal(b, &msg); nil != err {
					continue
				}
				details = append(details, &msg)
			}
		}
		if len(details) > 0 {
			statProto := spb.Status{
				Code:    int32(code),
				Message: msg,
				Details: details,
			}
			return status.FromProto(&statProto)
		}
		return status.New(code, msg)
	}
	return nil
}
