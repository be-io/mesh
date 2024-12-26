/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package node

import (
	"context"
	"github.com/opendatav/mesh/client/golang/plugin"
	// --
	_ "github.com/opendatav/mesh/plugin/metabase"
	_ "github.com/opendatav/mesh/plugin/proxy"
	_ "github.com/opendatav/mesh/plugin/prsim"
	_ "github.com/opendatav/mesh/plugin/redis"
	_ "github.com/opendatav/mesh/plugin/serve"
)

func init() {
	plugin.Provide(new(node))
}

type node struct {
}

func (that *node) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.NodeCar, Flags: node{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *node) Start(ctx context.Context, runtime plugin.Runtime) {
	runtime.Load(plugin.Http)
	runtime.Load(plugin.Proxy)
	runtime.Load(plugin.Redis)
}

func (that *node) Stop(ctx context.Context, runtime plugin.Runtime) {

}
