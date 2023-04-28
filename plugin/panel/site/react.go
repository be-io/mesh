/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package site

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/tool"
	"html/template"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/fatih/structs"
	echo "github.com/labstack/echo/v4"
	"github.com/olebedev/gojax/fetch"
)

// React struct is contains JS vms
// Runtime to serve HTTP requests and
// separates some domain specific
// resources.
type React struct {
	Runtime
	debug bool
	path  string
}

// NewReact initialized React struct
func NewReact(ctx context.Context, filePath string, debug bool, proxy http.Handler) *React {
	r := &React{
		debug: debug,
		path:  filePath,
	}
	if !debug {
		r.Runtime = newEnginePool(ctx, filePath, runtime.NumCPU(), proxy)
	} else {
		// Use daemon runtime to load full react
		// app each time for any http requests.
		// Useful to debug the app.
		r.Runtime = &daemonRuntime{
			ctx:   ctx,
			path:  filePath,
			proxy: proxy,
		}
	}
	return r
}

// Handle handles all HTTP requests which
// have not been caught via static file
// handler or other middlewares.
func (that *React) Handle(c echo.Context) error {
	uuid := tool.NextID()
	defer func() {
		if r := recover(); r != nil {
			c.Render(http.StatusInternalServerError, "react.html", Resp{
				UUID:  uuid,
				Error: r.(string),
			})
		}
	}()

	vm := that.Get()

	start := time.Now()
	select {
	case re := <-vm.Handle(map[string]interface{}{
		"url":     c.Request().URL.String(),
		"headers": map[string][]string(c.Request().Header),
		"uuid":    uuid,
	}):
		// Return vm back to the Runtime
		that.Put(vm)

		re.RenderTime = time.Since(start)

		// Handle the Response
		if len(re.Redirect) == 0 && len(re.Error) == 0 {
			// If no redirection and no errors
			c.Response().Header().Set("X-React-Render-Time", re.RenderTime.String())
			return c.Render(http.StatusOK, "react.html", re)
			// If redirect
		} else if len(re.Redirect) != 0 {
			return c.Redirect(http.StatusMovedPermanently, re.Redirect)
			// If internal error
		} else if len(re.Error) != 0 {
			c.Response().Header().Set("X-React-Render-Time", re.RenderTime.String())
			return c.Render(http.StatusInternalServerError, "react.html", re)
		}
	case <-time.After(2 * time.Second):
		// release the context
		that.Release(vm)
		return c.Render(http.StatusInternalServerError, "react.html", Resp{
			UUID:  tool.NextID(),
			Error: "timeout",
		})
	}
	return nil
}

// Resp is a struct for convinient
// react app Response parsing.
// Feel free to add any other keys to this struct
// and return value for this key at ecmascript side.
// Keep it sync with: src/app/client/router/toString.js:23
type Resp struct {
	UUID       string        `json:"uuid"`
	Error      string        `json:"error"`
	Redirect   string        `json:"redirect"`
	App        string        `json:"app"`
	Title      string        `json:"title"`
	Meta       string        `json:"meta"`
	Initial    string        `json:"initial"`
	RenderTime time.Duration `json:"-"`
}

// HTMLApp returns a application template
func (r Resp) HTMLApp() template.HTML {
	return template.HTML(r.App)
}

// HTMLTitle returns a title data
func (r Resp) HTMLTitle() template.HTML {
	return template.HTML(r.Title)
}

// HTMLMeta returns a meta data
func (r Resp) HTMLMeta() template.HTML {
	return template.HTML(r.Meta)
}

// Runtime Interface to serve React app on demand or from prepared Runtime.
type Runtime interface {
	Get() *v8
	Put(v8 *v8)
	Release(v8 *v8)
}

// newEnginePool return Runtime of JS vms.
func newEnginePool(ctx context.Context, filePath string, size int, proxy http.Handler) *debugRuntime {
	es := &debugRuntime{
		ctx:   mpc.Context(),
		path:  filePath,
		ch:    make(chan *v8, size),
		proxy: proxy,
	}
	go func() {
		for i := 0; i < size; i++ {
			es.ch <- newV8(ctx, filePath, proxy)
		}
	}()

	return es
}

type debugRuntime struct {
	ctx   context.Context
	ch    chan *v8
	path  string
	proxy http.Handler
}

func (that *debugRuntime) Get() *v8 {
	return <-that.ch
}

func (that *debugRuntime) Put(ot *v8) {
	that.ch <- ot
}

func (that *debugRuntime) Release(ot *v8) {
	ot.Stop()
	ot = nil
	that.ch <- newV8(that.ctx, that.path, that.proxy)
}

type daemonRuntime struct {
	ctx   context.Context
	path  string
	proxy http.Handler
}

func (that *daemonRuntime) Get() *v8 {
	return newV8(that.ctx, that.path, that.proxy)
}

func (that *daemonRuntime) Put(c *v8) {
	c.Stop()
}

func (that *daemonRuntime) Release(c *v8) {
	that.Put(c)
}

func newV8(ctx context.Context, filePath string, proxy http.Handler) *v8 {
	log.Info(ctx, "V8 run with %s", filePath)
	vm := &v8{
		EventLoop: eventloop.NewEventLoop(),
		ch:        make(chan Resp, 1),
	}
	vm.EventLoop.Start()
	if err := fetch.Enable(vm.EventLoop, proxy); nil != err {

	}
	vm.EventLoop.RunOnLoop(func(_vm *goja.Runtime) {
		var seed int64
		if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
			panic(fmt.Errorf("Could not read random bytes: %v", err))
		}
		_vm.SetRandSource(rand.New(rand.NewSource(seed)).Float64)

		_, err := _vm.RunScript("bundle.js", filePath)
		if err != nil {
			panic(err)
		}

		if fn, ok := goja.AssertFunction(_vm.Get("main")); ok {
			vm.fn = fn
		} else {
			fmt.Println("fn assert failed")
		}

		_vm.Set("__goServerCallback__", func(call goja.FunctionCall) goja.Value {
			obj := call.Argument(0).Export().(map[string]interface{})
			re := &Resp{}
			for _, field := range structs.Fields(re) {
				if n := field.Tag("json"); len(n) > 1 {
					field.Set(obj[n])
				}
			}
			vm.ch <- *re
			return nil
		})
	})

	return vm
}

// v8 wraps goja EventLoop
type v8 struct {
	*eventloop.EventLoop
	ch chan Resp
	fn goja.Callable
}

// Handle handles http requests
func (r *v8) Handle(req map[string]interface{}) <-chan Resp {
	r.EventLoop.RunOnLoop(func(vm *goja.Runtime) {
		r.fn(nil, vm.ToValue(req), vm.ToValue("__goServerCallback__"))
	})
	return r.ch
}
