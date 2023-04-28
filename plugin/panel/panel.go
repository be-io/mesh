/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package panel

import (
	"context"
	"embed"
	"github.com/be-io/mesh/client/golang/cause"
	httpx "github.com/be-io/mesh/client/golang/http"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/plugin"
	"net/http"
	"os"
)

func init() {
	mfs := &panel{
		FS: &httpx.StaticFileSystem{
			FS:      home,
			Home:    "static",
			Name:    "mesh.plugin.panel.home",
			Pattern: "/mesh",
		},
	}
	plugin.Provide(mfs)
	macro.Provide((*http.FileSystem)(nil), mfs)
	macro.Provide((*http.FileSystem)(nil), &favicon{
		FS: &httpx.StaticFileSystem{
			FS:      home,
			Home:    "static",
			Name:    "mesh.plugin.panel.favicon",
			Pattern: "/favicon.ico",
		},
	})
}

//go:embed static
var home embed.FS

type panel struct {
	FS httpx.VFS
}

func (that *panel) Att() *macro.Att {
	return that.FS.Att()
}

func (that *panel) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.Panel, Flags: panel{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *panel) Start(ctx context.Context, runtime plugin.Runtime) {

}

func (that *panel) Stop(ctx context.Context, runtime plugin.Runtime) {

}

func (that *panel) Open(name string) (http.File, error) {
	if f, err := that.FS.Open(name); os.IsNotExist(cause.DeError(err)) {
		return that.FS.Open("/")
	} else {
		return f, err
	}
}

type favicon struct {
	FS httpx.VFS
}

func (that *favicon) Att() *macro.Att {
	return that.FS.Att()
}

func (that *favicon) Open(name string) (http.File, error) {
	return that.FS.Open(name)
}
