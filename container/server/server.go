/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package server

import (
	"context"
	"github.com/opendatav/mesh/client/golang/plugin"
	// --
	_ "github.com/opendatav/mesh/plugin/cache"
	_ "github.com/opendatav/mesh/plugin/metabase"
	_ "github.com/opendatav/mesh/plugin/nsq"
	_ "github.com/opendatav/mesh/plugin/panel"
	_ "github.com/opendatav/mesh/plugin/proxy"
	_ "github.com/opendatav/mesh/plugin/prsim"
	_ "github.com/opendatav/mesh/plugin/redis"
)

func init() {
	plugin.Provide(new(server))
}

type server struct {
}

func (that *server) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.ServerCar, Flags: server{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *server) Start(ctx context.Context, runtime plugin.Runtime) {
	runtime.Load(plugin.Proxy)
	runtime.Load(plugin.Http)
	runtime.Load(plugin.Metabase)
	runtime.Load(plugin.NSQ)
	runtime.Load(plugin.PRSIM)
	runtime.Load(plugin.Panel)
	runtime.Load(plugin.Redis)
	runtime.Load(plugin.Cache)
}

func (that *server) Stop(ctx context.Context, runtime plugin.Runtime) {

}
