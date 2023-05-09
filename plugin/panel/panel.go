/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package panel

import (
	"context"
	"embed"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	httpx "github.com/be-io/mesh/client/golang/http"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/traefik/traefik/v2/webui"
	"net/http"
	"strings"
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
}

//go:embed static
var home embed.FS
var nui = webui.WFS()

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
	if !strings.Contains(name, "/network") {
		return that.FS.Open(name)
	}
	paths := strings.Split(name, "/network")
	path := fmt.Sprintf("static%s", strings.Join(paths[1:], "/"))
	f, err := nui.Open(path)
	if nil != err {
		return nil, cause.Error(err)
	}
	info, err := f.Stat()
	if nil != err {
		return nil, cause.Error(err)
	}
	if info.IsDir() {
		path = fmt.Sprintf("%s/index.html", path)
		f, err = nui.Open(path)
		if nil != err {
			return nil, cause.Error(err)
		}
	}
	return httpx.StaticFile(nui, f, path), nil
}
