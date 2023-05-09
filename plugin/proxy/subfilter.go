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
	"compress/gzip"
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	mtypes "github.com/be-io/mesh/client/golang/types"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/middlewares/stripprefix"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
	"github.com/traefik/traefik/v2/pkg/tracing"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func init() {
	middleware.Provide(filter)
}

const subfilterRuler = "subfilter"

var filter = &subfilterMiddleware{rewrites: dsa.NewStringMap[[]*Rewrite]()}

type subfilter struct {
	name string
	next http.Handler
}

func (that *subfilter) ParseRuler(request *http.Request) string {
	uname := prsim.MeshUrn.GetHeader(request.Header)
	if "" == uname {
		return subfilterRuler
	}
	urn := mtypes.FromURN(macro.Context(), uname)
	return urn.Name
}

func (that *subfilter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ruler := that.ParseRuler(request)
	rewrites, ok := filter.rewrites.Get(ruler)
	if !ok || len(rewrites) < 1 {
		that.next.ServeHTTP(writer, request)
		return
	}
	for _, rewrite := range rewrites {
		if strings.HasPrefix(request.URL.Path, rewrite.prefix) {
			request.URL.Path = that.getPrefixStripped(request.URL.Path, rewrite.prefix)
			if request.URL.RawPath != "" {
				request.URL.RawPath = that.getPrefixStripped(request.URL.RawPath, rewrite.prefix)
			}
			//*  - x-forwarded-proto://x-forwarded-host:x-forwarded-port/HttpServletRequest.getRequestURI()
			//*  - x-forwarded-proto://x-forwarded-host/HttpServletRequest.getRequestURI()
			//*  - x-forwarded-proto://host:x-forwarded-port/HttpServletRequest.getRequestURI()
			//*  - x-forwarded-proto://host/HttpServletRequest.getRequestURI() request.getRequestURL()
			proto := request.Header.Get("X-Forwarded-Proto")
			if strings.Contains(proto, "ws") {
				schema := tool.Ternary("wss" == proto, "https", "http")
				request.Header.Set("X-Forwarded-Proto", schema)
				request.Header.Set("X-API-request-url", fmt.Sprintf("%s://%s%s", schema, request.Host, request.URL.Path))
			}
			that.subfilter(writer, request, rewrite)
			return
		}
	}
	that.next.ServeHTTP(writer, request)
}

func (that *subfilter) subfilter(writer http.ResponseWriter, request *http.Request, pattern *Rewrite) {
	request.Header.Add(stripprefix.ForwardedPrefixHeader, pattern.prefix)
	request.RequestURI = request.URL.RequestURI()

	delegate := &subfilterResponseWriter{
		buffer:  &bytes.Buffer{},
		code:    http.StatusOK,
		headers: make(http.Header),
		writer:  writer,
	}

	that.next.ServeHTTP(delegate, request)

	buff := delegate.buffer.Bytes()
	contentEncoding := delegate.Header().Get("Content-Encoding")

	switch strings.ToLower(contentEncoding) {
	case "identity":
	case "":
		rep := []byte(fmt.Sprintf("${0}%s/", strings.TrimSuffix(strings.TrimPrefix(pattern.prefix, "/"), "/")))
		that.writeResponse(writer, delegate, pattern.ReplaceAll(buff, rep))
	case "gzip":
		explain, err := gzip.NewReader(bytes.NewReader(buff))
		if err != nil {
			log.Error0("Unable to create gzip reader: %s", err.Error())
			that.writeResponse(writer, delegate, buff)
			return
		}

		eff, err := io.ReadAll(explain)
		if err != nil {
			log.Error0("Unable to read gzipped response: %s", err.Error())
			that.writeResponse(writer, delegate, buff)
			return
		}

		var cff bytes.Buffer
		gz := gzip.NewWriter(&cff)

		rep := []byte(fmt.Sprintf("${0}%s/", strings.TrimSuffix(strings.TrimPrefix(pattern.prefix, "/"), "/")))
		_, err = gz.Write(pattern.ReplaceAll(eff, rep))
		if nil != err {
			log.Error0("Unable to write gzipped modified response: %s", err.Error())
			that.writeResponse(writer, delegate, buff)
			return
		}
		if err = gz.Close(); nil != err {
			log.Error0("Unable to close gzip writer: %s", err.Error())
			that.writeResponse(writer, delegate, buff)
			return
		}
		that.writeResponse(writer, delegate, cff.Bytes())
	default:
		that.writeResponse(writer, delegate, buff)
	}
}

func (that *subfilter) writeResponse(writer http.ResponseWriter, delegate *subfilterResponseWriter, buff []byte) {
	writer.Header().Del("Content-Length")
	writer.WriteHeader(delegate.code)
	if _, err := writer.Write(buff); nil != err {
		log.Error0(err.Error())
	}
}

func (that *subfilter) getPrefixStripped(urlPath, prefix string) string {
	rep := strings.TrimPrefix(urlPath, prefix)
	if rep == "" {
		return rep
	}
	if rep[0] != '/' {
		rep = "/" + rep
	}
	if !strings.HasPrefix(rep, prefix) {
		return rep
	}
	return that.getPrefixStripped(rep, prefix)
}

func (that *subfilter) GetTracingInformation() (string, ext.SpanKindEnum) {
	return that.name, tracing.SpanKindNoneEnum
}

type subfilterMiddleware struct {
	rewrites dsa.Map[string, []*Rewrite]
}

func (that *subfilterMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginSubfilter, ProviderName)
}

func (that *subfilterMiddleware) Priority() int {
	return 0
}

func (that *subfilterMiddleware) Scope() int {
	return 1
}

func (that *subfilterMiddleware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.subfilter"}
}

func (that *subfilterMiddleware) Refresh(rewrites []*Rewrite) {
	that.rewrites.Put(subfilterRuler, rewrites)
}

// New creates and returns a new rewrite body plugin instance.
func (that *subfilterMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &subfilter{name: name, next: next}, nil
}

type Rewrite struct {
	regex   *regexp.Regexp
	prefix  string
	filter  *regexp.Regexp
	replace string
}

func (that *Rewrite) ReplaceAll(src, repl []byte) []byte {
	if nil == that.regex && nil == that.filter {
		return src
	}
	if nil == that.filter {
		return that.regex.ReplaceAll(src, repl)
	}
	if nil == that.regex {
		return that.filter.ReplaceAll(src, []byte(that.replace))
	}
	return that.filter.ReplaceAll(that.regex.ReplaceAll(src, repl), []byte(that.replace))
}

type subfilterResponseWriter struct {
	buffer  *bytes.Buffer
	code    int
	headers http.Header
	writer  http.ResponseWriter
}

func (that *subfilterResponseWriter) Header() http.Header {
	return that.writer.Header()
}

func (that *subfilterResponseWriter) WriteHeader(statusCode int) {
	that.code = statusCode
}

func (that *subfilterResponseWriter) Write(p []byte) (int, error) {
	return that.buffer.Write(p)
}

func (that *subfilterResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := that.writer.(http.Hijacker)
	if !ok {
		return nil, nil, cause.Errorf("%T is not a http.Hijacker", that.writer)
	}

	return hijacker.Hijack()
}

func (that *subfilterResponseWriter) Flush() {
	if flusher, ok := that.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
