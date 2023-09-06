/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/fs"
	"net/http"
	"net/http/pprof"
	"path"
	"runtime/debug"
	"strings"
)

func init() {
	var _ macro.SPI = new(SPIDecorator)
	var _ http.Handler = new(SPIDecorator)
	var _ mpc.Provider = new(httpProvider)
	macro.Provide(mpc.IProvider, new(httpProvider))
}

const (
	Name = "http"
)

type SPIDecorator struct {
	Handler   http.Handler
	Attribute *macro.Att
}

func (that *SPIDecorator) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	that.Handler.ServeHTTP(writer, request)
}

func (that *SPIDecorator) Att() *macro.Att {
	return that.Attribute
}

type httpProvider struct {
	server *http.Server
}

func (that *httpProvider) Att() *macro.Att {
	return &macro.Att{Name: Name, Prototype: true, Constructor: func() macro.SPI {
		return new(httpProvider)
	}}
}

func (that *httpProvider) Start(ctx context.Context, address string, tc *tls.Config) error {
	router := httprouter.New()
	router.GET("/debug/", that.Handle(pprof.Index))
	router.GET("/debug/cmdline", that.Handle(pprof.Cmdline))
	router.GET("/debug/profile", that.Handle(pprof.Profile))
	router.POST("/debug/symbol", that.Handle(pprof.Symbol))
	router.GET("/debug/symbol", that.Handle(pprof.Symbol))
	router.GET("/debug/trace", that.Handle(pprof.Trace))
	router.GET("/debug/allocs", that.Handler(pprof.Handler("allocs")))
	router.GET("/debug/block", that.Handler(pprof.Handler("block")))
	router.GET("/debug/goroutine", that.Handler(pprof.Handler("goroutine")))
	router.GET("/debug/heap", that.Handler(pprof.Handler("heap")))
	router.GET("/debug/mutex", that.Handler(pprof.Handler("mutex")))
	router.GET("/debug/threadcreate", that.Handler(pprof.Handler("threadcreate")))
	router.GET("/stats", that.Stats)
	for _, pdd := range macro.Load((*http.Handler)(nil)).List() {
		if han, ok := pdd.(http.Handler); ok {
			that.Any(router, pdd.Att().Pattern, that.Handler(han))
		}
	}
	for _, pdd := range macro.Load((*http.FileSystem)(nil)).List() {
		if han, ok := pdd.(http.FileSystem); ok {
			that.StaticFS(router, pdd.Att().Pattern, han)
		}
	}
	that.Any(router, "/", that.Forward)
	that.Any(router, "/mesh/invoke", that.Forward)
	that.Back(router, that.Forbidden)
	that.server = &http.Server{
		Addr:    address,
		Handler: router,
	}
	log.Info(ctx, "Listening and serving HTTP 1.x on %s", address)
	if err := that.server.ListenAndServe(); nil != err {
		log.Error(ctx, "Http broker up with %s fail, %s", address, err.Error())
	}
	return nil
}

func (that *httpProvider) Close() error {
	return that.server.Close()
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (that *httpProvider) StaticFS(router *httprouter.Router, relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := http.StripPrefix(relativePath, &staticHttpProvider{fs: fs})
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET and HEAD handlers
	router.GET(urlPattern, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handler.ServeHTTP(writer, request)
	})
	router.HEAD(urlPattern, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handler.ServeHTTP(writer, request)
	})
}

func (that *httpProvider) Back(router *httprouter.Router, handler httprouter.Handle) {

}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (that *httpProvider) Any(router *httprouter.Router, relativePath string, handler httprouter.Handle) {
	router.Handle(http.MethodGet, relativePath, handler)
	router.Handle(http.MethodPost, relativePath, handler)
	router.Handle(http.MethodPut, relativePath, handler)
	router.Handle(http.MethodPatch, relativePath, handler)
	router.Handle(http.MethodHead, relativePath, handler)
	router.Handle(http.MethodOptions, relativePath, handler)
	router.Handle(http.MethodDelete, relativePath, handler)
	router.Handle(http.MethodConnect, relativePath, handler)
	router.Handle(http.MethodTrace, relativePath, handler)
}

func (that *httpProvider) Handle(handle func(w http.ResponseWriter, r *http.Request)) httprouter.Handle {
	return func(writer http.ResponseWriter, h *http.Request, params httprouter.Params) {
		handle(writer, h)
	}
}

func (that *httpProvider) Handler(handler http.Handler) httprouter.Handle {
	return func(writer http.ResponseWriter, h *http.Request, params httprouter.Params) {
		handler.ServeHTTP(writer, h)
	}
}

func FormatContentType(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func (that *httpProvider) Forward(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ctx := mpc.Context()
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, string(debug.Stack()))
			log.Error(ctx, "%v", err)
			that.Failure(ctx, writer, err)
		}
	}()
	urn, input, err := func() (*types.URN, *types.Request, error) {
		binding := Default(request.Method, FormatContentType(request.Header.Get("Content-Type")))
		mu := prsim.MeshUrn.GetHeader(request.Header)
		if "" != mu {
			content := map[string]any{}
			if err := binding.Bind(request, &content); nil != err {
				return nil, nil, cause.Error(err)
			}
			urn := types.FromURN(ctx, mu)
			input := &types.Request{Version: "", Method: urn.Name, Content: content}
			return urn, input, nil
		}
		input := &types.Request{}
		if err := binding.Bind(request, input); nil != err {
			return nil, nil, cause.Error(err)
		}
		return &types.URN{
			Domain: types.MeshDomain,
			NodeId: types.LocalNodeId,
			Flag:   &types.URNFlag{Proto: mpc.MeshFlag.GRPC.Code(), Codec: mpc.MeshFlag.JSON.Code()},
			Name:   input.Method,
		}, input, nil
	}()
	if nil != err {
		log.Error(ctx, "Forward %s, %s", request.RequestURI, err.Error())
		that.Failure(ctx, writer, err)
		return
	}
	if text, err := aware.Codec.EncodeString(input); nil != err {
		log.Info(ctx, err.Error())
	} else {
		log.Info(ctx, "%s", text)
	}
	for name, value := range request.Header {
		ctx.GetAttachments()[strings.ToLower(name)] = tool.Anyone(value...)
	}
	ctx.RewriteURN(urn.String())
	execution, err := aware.Eden.Infer(ctx, urn.String())
	if nil != err {
		that.Failure(ctx, writer, err)
		return
	}
	if nil == execution {
		that.Failure(ctx, writer, cause.Errorf("Service %s not exist.", input.Method))
		return
	}
	buff, err := aware.Codec.Encode(input.Content)
	if nil != err {
		that.Failure(ctx, writer, err)
		return
	}
	parameters := execution.Inspect().GetIntype()
	if nil != err {
		that.Failure(ctx, writer, err)
		return
	}
	if _, err = aware.Codec.Decode(buff, parameters); nil != err {
		that.Failure(ctx, writer, err)
		return
	}
	parameters.SetAttachments(ctx, ctx.GetAttachments())
	invocation := &mpc.ServiceInvocation{
		Proxy:      execution,
		Inspector:  execution.Inspect(),
		Parameters: parameters,
		Buffer:     buff,
		Execution:  execution,
		URN:        urn,
	}
	ret, err := execution.Invoke(ctx, invocation)
	if nil != err {
		log.Error(ctx, "Forward %s, %s", input.Method, err.Error())
		that.Failure(ctx, writer, err)
		return
	}
	returns := execution.Inspect().NewOutbound()
	returns.SetCode(cause.Success.Code)
	returns.SetMessage(cause.Success.Message)
	returns.SetContent(ctx, ret)
	that.Success(ctx, writer, returns)
}

func (that *httpProvider) Failure(ctx prsim.Context, writer http.ResponseWriter, except any) {
	outbound := &types.Outbound{}
	if cable, ok := except.(cause.Codeable); ok {
		outbound.Code = cable.GetCode()
		outbound.Message = cable.GetMessage()
	} else {
		outbound.Code = cause.SystemError.GetCode()
		outbound.Message = fmt.Sprintf("%v", except)
	}
	buff, err := aware.Codec.Encode(outbound)
	if nil != err {
		log.Error(ctx, err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(buff.Bytes())
	if nil != err {
		log.Error(ctx, err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (that *httpProvider) Success(ctx prsim.Context, writer http.ResponseWriter, output interface{}) {
	buff, err := aware.Codec.Encode(output)
	if nil != err {
		that.Failure(ctx, writer, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(buff.Bytes())
	if nil != err {
		that.Failure(ctx, writer, err)
		return
	}
}

func (that *httpProvider) Forbidden(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ctx := mpc.Context()
	writer.WriteHeader(http.StatusForbidden)
	_, err := writer.Write([]byte("Forbidden"))
	if nil != err {
		that.Failure(ctx, writer, err)
		return
	}
}

func (that *httpProvider) Stats(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ctx := mpc.Context()
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(cause.Success.Code))
	if nil != err {
		that.Failure(ctx, writer, err)
		return
	}
}

type staticHttpProvider struct {
	fs http.FileSystem
}

func (that *staticHttpProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	that.serveFile(w, r, that.fs, path.Clean(upath))
}

// name is '/'-separated, not filepath.Separator.
func (that *staticHttpProvider) serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string) {
	f, err := fs.Open(name)
	if nil != err {
		msg, code := that.toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer func() { log.Catch(f.Close()) }()
	d, err := f.Stat()
	if nil != err {
		msg, code := that.toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

// toHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that toHTTPError does not
// actually return err.Error(), since msg and httpStatus are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func (that *staticHttpProvider) toHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

type VFS interface {
	macro.SPI
	http.FileSystem
}
type FS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}
type StaticFileSystem struct {
	FS      FS
	Home    string
	Name    string
	Pattern string
}

func (that *StaticFileSystem) Att() *macro.Att {
	return &macro.Att{Name: that.Name, Pattern: that.Pattern}
}

func (that *StaticFileSystem) Open(name string) (http.File, error) {
	relativePath := fmt.Sprintf("%s%s", that.Home, tool.Ternary("/" == name, "/index.html", name))
	file, err := that.FS.Open(relativePath)
	if nil != err {
		return nil, cause.Error(err)
	}
	return &staticFile{fs: that.FS, file: file, name: relativePath, ctx: macro.Context()}, nil
}

func StaticFile(fs FS, f fs.File, path string) http.File {
	return &staticFile{fs: fs, file: f, name: path, ctx: macro.Context()}
}

type staticFile struct {
	ctx  context.Context
	name string
	file fs.File
	fs   FS
}

func (that *staticFile) Stat() (fs.FileInfo, error) {
	return that.file.Stat()
}

func (that *staticFile) Read(buffer []byte) (int, error) {
	return that.file.Read(buffer)
}
func (that *staticFile) Close() error {
	return that.file.Close()
}

func (that *staticFile) Seek(offset int64, whence int) (int64, error) {
	fb, err := that.fs.ReadFile(that.name)
	if nil != err {
		log.Error(that.ctx, "Seek file %s failed, %s", that.name, err.Error())
		return -1, cause.Error(err)
	}
	switch whence {
	case io.SeekStart:
		return offset, nil
	case io.SeekCurrent:
		return offset, nil
	case io.SeekEnd:
		return int64(len(fb)), nil
	default:
		return offset, nil
	}
}

func (that *staticFile) Readdir(count int) ([]fs.FileInfo, error) {
	entries, err := that.fs.ReadDir(that.name)
	if nil != err {
		log.Error(that.ctx, "Read dir %s failed, %s", that.name, err.Error())
		return nil, cause.Error(err)
	}
	var infos []fs.FileInfo
	for index, entry := range entries {
		if index > count {
			continue
		}
		if info, err := entry.Info(); nil != err {
			log.Error(that.ctx, "Read dir %s failed, %s", entry.Name(), err.Error())
			return nil, cause.Error(err)
		} else {
			infos = append(infos, info)
		}
	}
	return infos, nil
}
