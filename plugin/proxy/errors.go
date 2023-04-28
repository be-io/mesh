/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	mtypes "github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/server/middleware"
	"github.com/traefik/traefik/v3/pkg/types"
	"net"
	"net/http"
	"strings"
)

func init() {
	var _ http.Handler = new(errors)
	middleware.Provide(errs)
	var _ prsim.Listener = errs
	macro.Provide(prsim.IListener, errs)
}

var errs = new(errorsMiddleware)

type errors struct {
	name           string
	ctx            context.Context
	httpCodeRanges types.HTTPCodeRanges
	next           http.Handler
}

func (that *errors) copyRequest(request *http.Request, code int) (*http.Request, error) {
	newRequest := request.Clone(request.Context())
	urn := prsim.MeshUrn.GetHeader(newRequest.Header)
	uname := mtypes.FromURN(mpc.Context(), urn)
	if nil != errs.env && !tool.IsLocalEnv(errs.env, uname.NodeId) && !tool.Contains(errs.nodeIds, strings.ToUpper(uname.NodeId)) {
		uname.NodeId = mtypes.LocalNodeId
		uname.Name = fmt.Sprintf("%s.%d", "mesh.builtin.fallback", 618)
	} else {
		uname.Name = fmt.Sprintf("%s.%d", "mesh.builtin.fallback", code)
	}
	newRequest.Host = uname.String()
	prsim.MeshUrn.SetHeader(newRequest.Header, newRequest.Host)
	return newRequest, nil
}

func (that *errors) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if nil == proxy.TCPRouters || nil == proxy.TCPRouters[TransportY] {
		that.next.ServeHTTP(writer, request)
		return
	}
	catcher := newCodeCatcher(writer, that.httpCodeRanges)
	that.next.ServeHTTP(catcher, request)
	if !catcher.CodeMatches() {
		return
	}

	// check the recorder code against the configured http status code ranges
	code := catcher.getCode()
	log.Debug(that.ctx, "Caught HTTP Status Code %d, returning error page", code)

	if request.ProtoMajor < 2 {
		catcher.FlushWithCode(tool.Ternary(code == http.StatusNotFound, http.StatusForbidden, code))
		return
	}
	newRequest, err := that.copyRequest(request, code)
	if nil != err {
		log.Error(that.ctx, err.Error())
		writer.WriteHeader(code)
		if _, err = writer.Write([]byte(err.Error())); nil != err {
			log.Error(that.ctx, err.Error())
		}
		return
	}
	proxy.TCPRouters[TransportY].GetHTTPHandler().ServeHTTP(writer, newRequest)
}

type errorsMiddleware struct {
	nodeIds []string
	env     *mtypes.Environ
}

func (that *errorsMiddleware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.errors"}
}

func (that *errorsMiddleware) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.NetworkRouteRefresh}
}

func (that *errorsMiddleware) Listen(ctx context.Context, event *mtypes.Event) error {
	var routes []*mtypes.Route
	if err := event.TryGetObject(&routes); nil != err {
		return cause.Error(err)
	}
	var nodeIds []string
	for _, route := range routes {
		nodeIds = append(nodeIds, route.NodeId, route.InstId)
	}
	env, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	that.nodeIds = nodeIds
	that.env = env
	return nil
}

func (that *errorsMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginErrors, ProviderName)
}

func (that *errorsMiddleware) Priority() int {
	return 0
}

func (that *errorsMiddleware) Scope() int {
	return 0
}

func (that *errorsMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	httpCodeRanges, err := types.NewHTTPCodeRanges([]string{"404"})
	if err != nil {
		return nil, err
	}
	return &errors{name: name, ctx: mpc.Context(), httpCodeRanges: httpCodeRanges, next: next}, nil
}

type responseInterceptor interface {
	http.ResponseWriter
	http.Flusher
	getCode() int
	CodeMatches() bool
	FlushWithCode(code int)
}

// codeCatcher is a response writer that detects as soon as possible whether the
// response is a code within the ranges of codes it watches for. If it is, it
// simply drops the data from the response. Otherwise, it forwards it directly to
// the original client (its responseWriter) without any buffering.
type codeCatcher struct {
	httpCodeRanges types.HTTPCodeRanges
	rw             http.ResponseWriter
	headers        http.Header
	code           int
	buffer         *bytes.Buffer
}

type codeCatcherWithCloseNotify struct {
	*codeCatcher
}

// CloseNotify returns a channel that receives at most a
// single value (true) when the client connection has gone away.
func (that *codeCatcherWithCloseNotify) CloseNotify() <-chan bool {
	return that.rw.(http.CloseNotifier).CloseNotify()
}

func newCodeCatcher(rw http.ResponseWriter, httpCodeRanges types.HTTPCodeRanges) responseInterceptor {
	catcher := &codeCatcher{
		httpCodeRanges: httpCodeRanges,
		rw:             rw,
		headers:        make(http.Header),
		code:           http.StatusOK, // If backend does not call WriteHeader on us, we consider it's a 200.
		buffer:         &bytes.Buffer{},
	}
	if _, ok := rw.(http.CloseNotifier); ok {
		return &codeCatcherWithCloseNotify{catcher}
	}
	return catcher
}

func (that *codeCatcher) Header() http.Header {
	if !that.CodeMatches() {
		return that.rw.Header()
	}
	if nil == that.headers {
		that.headers = make(http.Header)
	}
	return that.headers
}

func (that *codeCatcher) getCode() int {
	return that.code
}

// CodeMatches returns whether the codeCatcher received a response code among the ones it is watching,
// and for which the response should be deferred to the error handler.
func (that *codeCatcher) CodeMatches() bool {
	return that.httpCodeRanges.Contains(that.code)
}

func (that *codeCatcher) Write(buf []byte) (int, error) {
	if !that.CodeMatches() {
		return that.rw.Write(buf)
	}
	that.buffer.Write(buf)
	return len(buf), nil
}

func (that *codeCatcher) WriteHeader(code int) {
	if !that.httpCodeRanges.Contains(code) {
		that.rw.WriteHeader(code)
		return
	}
	that.code = code
}

// Hijack hijacks the connection.
func (that *codeCatcher) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := that.rw.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("%T is not a http.Hijacker", that.rw)
}

// FlushWithCode sends any buffered data to the client.
func (that *codeCatcher) FlushWithCode(code int) {
	if that.CodeMatches() {
		prsim.CopyHeadersOverride(that.rw.Header(), that.Header())
		that.rw.WriteHeader(code)
		if _, err := that.rw.Write(that.buffer.Bytes()); nil != err {
			log.Error0(err.Error())
		}
	}
	that.Flush()
}

// Flush sends any buffered data to the client.
func (that *codeCatcher) Flush() {
	// If WriteHeader was already called from the caller, this is a NOOP.
	// Otherwise, cc.code is actually a 200 here.
	if flusher, ok := that.rw.(http.Flusher); ok {
		flusher.Flush()
	}
}
