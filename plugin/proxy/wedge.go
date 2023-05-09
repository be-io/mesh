/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
	"net/http"
	"sort"
)

func init() {
	var _ http.Handler = new(wedgeBox)
	middleware.Provide(dslBoxes)
}

var dslBoxes = &wedgeBoxMiddleware{wedges: map[string]wedges{}}

type wedgeBox struct {
	name string
	next http.Handler
}

func (that *wedgeBox) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	dslBoxes.ServeHTTP(that.next, writer, request)
}

type wedgeBoxMiddleware struct {
	wedges map[string]wedges
}

func (that *wedgeBoxMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginWedge, ProviderName)
}

func (that *wedgeBoxMiddleware) Priority() int {
	return 0
}

func (that *wedgeBoxMiddleware) Scope() int {
	return 0
}

func (that *wedgeBoxMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &wedgeBox{next: next, name: name}, nil
}

func (that *wedgeBoxMiddleware) ServeHTTP(next http.Handler, writer http.ResponseWriter, request *http.Request) {
	for p, wedge := range that.wedges {
		ok, err := wedge.Match(p, writer, request)
		if nil != err {
			log.Warn0(err.Error())
			continue
		}
		if ok {

		}
	}
	next.ServeHTTP(writer, request)
}

type Wedge interface {
	macro.SPI
	// Tap
	//
	//	proto==1.1 && path:startWith:/web-proxy {
	//		strip_prefix /web-proxy
	//	}
	Tap(next http.Handler, writer http.ResponseWriter, request *http.Request)
}

var IWedge = (*Wedge)(nil)

type wedges []Wedge

func (that wedges) Len() int {
	return len(that)
}

func (that wedges) Less(x, y int) bool {
	return that[x].Att().Priority < that[y].Att().Priority
}

func (that wedges) Swap(x, y int) {
	tmp := that[y]
	that[y] = that[x]
	that[x] = tmp
}

func (that wedges) Match(pattern string, writer http.ResponseWriter, request *http.Request) (bool, error) {
	return true, nil
}

func (that wedges) Serve(h http.Handler, pattern string) http.Handler {
	var ws wedges
	ps := macro.Load(IWedge).List()
	for _, p := range ps {
		w, ok := p.(Wedge)
		if !ok || w.Att().Pattern != pattern {
			continue
		}
		ws = append(ws, w)
	}
	sort.Sort(ws)
	last := h
	for i := len(ws) - 1; i >= 0; i-- {
		w := ws[i]
		next := last
		var x http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
			w.Tap(next, writer, request)
		}
		last = x
	}
	return last
}

var _ Wedge = new(stripePrefix)

type stripePrefix struct {
}

func (that *stripePrefix) Att() *macro.Att {
	return &macro.Att{Name: "strip"}
}

func (that *stripePrefix) Tap(next http.Handler, writer http.ResponseWriter, request *http.Request) {
	request.URL.Path = that.ensureLeadingSlash("gaia/v1/jupyter" + request.URL.Path)
	if request.URL.RawPath != "" {
		request.URL.RawPath = that.ensureLeadingSlash("gaia/v1/jupyter" + request.URL.RawPath)
	}
	request.RequestURI = request.URL.RequestURI()
}

func (that *stripePrefix) ensureLeadingSlash(str string) string {
	if str == "" {
		return str
	}
	if str[0] == '/' {
		return str
	}
	return "/" + str
}
