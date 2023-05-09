/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"io"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Server struct {
	mux       http.ServeMux
	handlers  HandlerMap
	basePath  string
	unaryInt  grpc.UnaryServerInterceptor
	streamInt grpc.StreamServerInterceptor
	opts      handlerOpts
}

var _ http.Handler = (*Server)(nil)
var _ grpc.ServiceRegistrar = (*Server)(nil)

type ServerOption interface {
	apply(*Server)
}

type serverOptFunc func(*Server)

func (that serverOptFunc) apply(s *Server) {
	that(s)
}

func WithBasePath(path string) ServerOption {
	return serverOptFunc(func(s *Server) {
		s.basePath = path
	})
}

func WithServerUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) ServerOption {
	return serverOptFunc(func(s *Server) {
		s.unaryInt = interceptor
	})
}

func WithServerStreamInterceptor(interceptor grpc.StreamServerInterceptor) ServerOption {
	return serverOptFunc(func(s *Server) {
		s.streamInt = interceptor
	})
}

func NewServer(opts ...ServerOption) *Server {
	var s Server
	s.basePath = "/"
	s.handlers = HandlerMap{}
	for _, o := range opts {
		o.apply(&s)
	}
	return &s
}

func (that *Server) RegisterService(desc *grpc.ServiceDesc, svr interface{}) {
	that.handlers.RegisterService(desc, svr)
	for i := range desc.Methods {
		md := desc.Methods[i]
		h := handleMethod(svr, desc.ServiceName, &md, that.unaryInt, &that.opts)
		that.mux.HandleFunc(path.Join(that.basePath, fmt.Sprintf("%s/%s", desc.ServiceName, md.MethodName)), h)
	}
	for i := range desc.Streams {
		sd := desc.Streams[i]
		h := handleStream(svr, desc.ServiceName, &sd, that.streamInt, &that.opts)
		that.mux.HandleFunc(path.Join(that.basePath, fmt.Sprintf("%s/%s", desc.ServiceName, sd.StreamName)), h)
	}
}
func (that *Server) GetServiceInfo() map[string]grpc.ServiceInfo {
	return that.handlers.GetServiceInfo()
}

func (that *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	that.mux.ServeHTTP(w, r)
}

type Mux func(pattern string, handler func(http.ResponseWriter, *http.Request))

type HandlerOption func(*handlerOpts)

func (ho HandlerOption) apply(s *Server) {
	ho(&s.opts)
}

type handlerOpts struct {
	errFunc func(context.Context, *status.Status, http.ResponseWriter)
}

func ErrorRenderer(errFunc func(reqCtx context.Context, st *status.Status, response http.ResponseWriter)) HandlerOption {
	return func(h *handlerOpts) {
		h.errFunc = errFunc
	}
}

// DefaultErrorRenderer
//
//	Canceled:         * 502 Bad Gateway
//	Unknown:            500 Internal Server Error
//	InvalidArgument:    400 Bad Request
//	DeadlineExceeded: * 504 Gateway Timeout
//	NotFound:           404 Not Found
//	AlreadyExists:      409 Conflict
//	PermissionDenied:   403 Forbidden
//	Unauthenticated:    401 Unauthorized
//	ResourceExhausted:  429 Too Many Requests
//	FailedPrecondition: 412 Precondition Failed
//	Aborted:            409 Conflict
//	OutOfRange:         422 Unprocessable Entity
//	Unimplemented:      501 Not Implemented
//	Internal:           500 Internal Server Error
//	Unavailable:        503 Service Unavailable
//	DataLoss:           500 Internal Server Error
func DefaultErrorRenderer(ctx context.Context, st *status.Status, w http.ResponseWriter) {
	if (st.Code() == codes.Canceled || st.Code() == codes.DeadlineExceeded) && ctx.Err() != nil {
		http.Error(w, "Client Closed Request", 499)
		return
	}
	code := httpStatusFromCode(st.Code())
	msg := http.StatusText(code)
	if msg == "" {
		msg = st.Code().String()
	}
	http.Error(w, msg, code)
}

func HandleServices(mux Mux, basePath string, reg HandlerMap, unaryInt grpc.UnaryServerInterceptor, streamInt grpc.StreamServerInterceptor, opts ...HandlerOption) {
	var hOpts handlerOpts
	for _, opt := range opts {
		opt(&hOpts)
	}

	reg.ForEach(func(desc *grpc.ServiceDesc, svr interface{}) {
		for i := range desc.Methods {
			md := desc.Methods[i]
			h := handleMethod(svr, desc.ServiceName, &md, unaryInt, &hOpts)
			mux(path.Join(basePath, fmt.Sprintf("%s/%s", desc.ServiceName, md.MethodName)), h)
		}
		for i := range desc.Streams {
			sd := desc.Streams[i]
			h := handleStream(svr, desc.ServiceName, &sd, streamInt, &hOpts)
			mux(path.Join(basePath, fmt.Sprintf("%s/%s", desc.ServiceName, sd.StreamName)), h)
		}
	})
}

func HandleMethod(svr interface{}, serviceName string, desc *grpc.MethodDesc, unaryInt grpc.UnaryServerInterceptor, opts ...HandlerOption) http.HandlerFunc {
	var hOpts handlerOpts
	for _, opt := range opts {
		opt(&hOpts)
	}
	return handleMethod(svr, serviceName, desc, unaryInt, &hOpts)
}

func handleMethod(svr interface{}, serviceName string, desc *grpc.MethodDesc, unaryInt grpc.UnaryServerInterceptor, opts *handlerOpts) http.HandlerFunc {
	errHandler := opts.errFunc
	if errHandler == nil {
		errHandler = DefaultErrorRenderer
	}
	fullMethod := fmt.Sprintf("/%s/%s", serviceName, desc.MethodName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if p := peerFromRequest(r); p != nil {
			ctx = peer.NewContext(ctx, p)
		}
		defer func() { _ = drainAndClose(r.Body) }()
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			writeError(w, http.StatusMethodNotAllowed)
			return
		}

		contentType := r.Header.Get("Content-Type")
		codec := getUnaryCodec(contentType)
		if codec == nil {
			writeError(w, http.StatusUnsupportedMediaType)
			return
		}

		ctx, cancel, err := contextFromHeaders(ctx, r.Header)
		if nil != err {
			writeError(w, http.StatusBadRequest)
			return
		}
		defer cancel()

		req, err := io.ReadAll(r.Body)
		if nil != err {
			writeError(w, 499)
			return
		}

		dec := func(msg interface{}) error {
			if err := codec.Unmarshal(req, msg); nil != err {
				return status.Error(codes.InvalidArgument, err.Error())
			}
			return nil
		}
		sts := UnaryServerTransportStream{Name: fullMethod}
		resp, err := desc.Handler(svr, grpc.NewContextWithServerTransportStream(ctx, &sts), dec, unaryInt)
		toHeaders(sts.GetHeaders(), w.Header(), "")
		toHeaders(sts.GetTrailers(), w.Header(), "X-GRPC-Trailer-")
		if nil != err {
			st, _ := status.FromError(err)
			if st.Code() == codes.OK {
				// preserve all error details, but rewrite the code since we don't want
				// to send back a non-error status when we know an error occured
				stpb := st.Proto()
				stpb.Code = int32(codes.Internal)
				st = status.FromProto(stpb)
			}
			statProto := st.Proto()
			w.Header().Set("X-GRPC-Status", fmt.Sprintf("%d:%s", statProto.Code, statProto.Message))
			for _, d := range statProto.Details {
				b, err := codec.Marshal(d)
				if nil != err {
					continue
				}
				str := base64.RawURLEncoding.EncodeToString(b)
				w.Header().Add(grpcDetailsHeader, str)
			}
			errHandler(r.Context(), st, w)
			return
		}

		b, err := codec.Marshal(resp)
		if nil != err {
			writeError(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
		w.Write(b)
	}
}

func HandleStream(svr interface{}, serviceName string, desc *grpc.StreamDesc, streamInt grpc.StreamServerInterceptor, opts ...HandlerOption) http.HandlerFunc {
	var hOpts handlerOpts
	for _, opt := range opts {
		opt(&hOpts)
	}
	return handleStream(svr, serviceName, desc, streamInt, &hOpts)
}

func handleStream(svr interface{}, serviceName string, desc *grpc.StreamDesc, streamInt grpc.StreamServerInterceptor, opts *handlerOpts) http.HandlerFunc {
	info := &grpc.StreamServerInfo{
		FullMethod:     fmt.Sprintf("/%s/%s", serviceName, desc.StreamName),
		IsClientStream: desc.ClientStreams,
		IsServerStream: desc.ServerStreams,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if p := peerFromRequest(r); p != nil {
			ctx = peer.NewContext(ctx, p)
		}
		defer func() { _ = drainAndClose(r.Body) }()
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			writeError(w, http.StatusMethodNotAllowed)
			return
		}

		contentType := r.Header.Get("Content-Type")
		codec := getStreamingCodec(contentType)
		if codec == nil {
			writeError(w, http.StatusUnsupportedMediaType)
			return
		}

		ctx, cancel, err := contextFromHeaders(ctx, r.Header)
		if nil != err {
			writeError(w, http.StatusBadRequest)
			return
		}
		defer cancel()

		w.Header().Set("Content-Type", contentType)

		str := &serverStream{r: r, w: w, respStream: desc.ClientStreams, codec: codec}
		sts := ServerTransportStream{Name: info.FullMethod, Stream: str}
		str.ctx = grpc.NewContextWithServerTransportStream(ctx, &sts)
		if streamInt != nil {
			err = streamInt(svr, str, info, desc.Handler)
		} else {
			err = desc.Handler(svr, str)
		}
		if str.writeFailed {
			// nothing else we can do
			return
		}

		tr := HttpTrailer{
			Code:     int32(codes.OK),
			Message:  codes.OK.String(),
			Metadata: asTrailerProto(metadata.Join(str.tr...)),
		}
		if nil != err {
			st, _ := status.FromError(err)
			if st.Code() == codes.OK {
				// preserve all error details, but rewrite the code since we don't want
				// to send back a non-error status when we know an error occured
				stpb := st.Proto()
				stpb.Code = int32(codes.Internal)
				st = status.FromProto(stpb)
			}
			statProto := st.Proto()
			tr.Code = statProto.Code
			tr.Message = statProto.Message
			tr.Details = statProto.Details
		}

		log.Catch(writeMessage(w, codec, &tr, true))
	}
}

func peerFromRequest(r *http.Request) *peer.Peer {
	pr := peer.Peer{Addr: strAddr(r.RemoteAddr)}
	if r.TLS != nil {
		pr.AuthInfo = credentials.TLSInfo{State: *r.TLS}
	}
	return &pr
}

func drainAndClose(r io.ReadCloser) error {
	_, copyErr := io.Copy(io.Discard, r)
	closeErr := r.Close()
	// error from io.Copy likely more useful than the one from Close
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func writeError(w http.ResponseWriter, code int) {
	msg := http.StatusText(code)
	if msg == "" {
		if code == 499 {
			msg = "Client Closed Request"
		} else {
			msg = "Unknown"
		}
	}
	http.Error(w, msg, code)
}

func asTrailerProto(md metadata.MD) map[string]*TrailerValues {
	result := map[string]*TrailerValues{}
	for k, vs := range md {
		tvs := TrailerValues{}
		tvs.Values = append(tvs.Values, vs...)
		result[k] = &tvs
	}
	return result
}

type serverStream struct {
	ctx         context.Context
	respStream  bool
	codec       encoding.Codec
	rmu         sync.Mutex
	r           *http.Request
	recvd       int
	wmu         sync.Mutex
	w           http.ResponseWriter
	headersSent bool
	writeFailed bool
	tr          []metadata.MD
}

func (that *serverStream) SetHeader(md metadata.MD) error {
	return that.setHeader(md, false)
}

func (that *serverStream) SendHeader(md metadata.MD) error {
	return that.setHeader(md, true)
}

func (that *serverStream) setHeader(md metadata.MD, send bool) error {
	that.wmu.Lock()
	defer that.wmu.Unlock()

	if that.headersSent {
		return errors.New("headers already sent")
	}

	h := that.w.Header()
	toHeaders(md, h, "")

	if send {
		that.w.WriteHeader(http.StatusOK)
		that.headersSent = true
	}

	return nil
}

func (that *serverStream) SetTrailer(md metadata.MD) {
	that.wmu.Lock()
	defer that.wmu.Unlock()

	that.tr = append(that.tr, md)
}

func (that *serverStream) Context() context.Context {
	return that.ctx
}

func (that *serverStream) SendMsg(m interface{}) error {
	that.wmu.Lock()
	defer that.wmu.Unlock()

	if that.writeFailed {
		return io.EOF
	}

	that.headersSent = true // sent implicitly
	err := writeMessage(that.w, that.codec, m, false)
	if nil != err {
		that.writeFailed = true
	}
	return err
}

func (that *serverStream) RecvMsg(m interface{}) error {
	that.rmu.Lock()
	defer that.rmu.Unlock()

	if !that.respStream && that.recvd > 0 {
		return cause.Error(io.EOF)
	}

	that.recvd++

	size, err := readSizePreface(that.r.Body)
	if nil != err {
		return cause.Error(err)
	}

	err = readMessage(that.r.Body, that.codec, size, m)
	if err == io.EOF {
		return cause.Error(io.ErrUnexpectedEOF)
	} else if nil != err {
		return cause.Error(err)
	}

	if !that.respStream {
		_, err = readSizePreface(that.r.Body)
		if err != io.EOF {
			// client tried to send >1 message!
			return cause.Error(status.Error(codes.InvalidArgument, "method accepts 1 request message but client sent >1"))
		}
	}

	return nil
}

func contextFromHeaders(parent context.Context, h http.Header) (context.Context, context.CancelFunc, error) {
	cancel := func() {} // default to no-op
	md, err := asMetadata(h)
	if nil != err {
		return parent, cancel, err
	}
	ctx := metadata.NewIncomingContext(parent, md)

	// deadline propagation
	timeout := h.Get("GRPC-Timeout")
	if timeout != "" {
		// See GRPC wire format, "Timeout" component of request: https://grpc.io/docs/guides/wire.html#requests
		suffix := timeout[len(timeout)-1]
		if timeoutVal, err := strconv.ParseInt(timeout[:len(timeout)-1], 10, 64); err == nil {
			var unit time.Duration
			switch suffix {
			case 'H':
				unit = time.Hour
			case 'M':
				unit = time.Minute
			case 'S':
				unit = time.Second
			case 'm':
				unit = time.Millisecond
			case 'u':
				unit = time.Microsecond
			case 'n':
				unit = time.Nanosecond
			}
			if unit != 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutVal)*unit)
			}
		}
	}
	return ctx, cancel, nil
}
