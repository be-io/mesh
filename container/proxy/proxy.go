/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"github.com/be-io/mesh/client/golang/plugin"
	// --
	_ "github.com/be-io/mesh/plugin/panel"
	_ "github.com/be-io/mesh/plugin/proxy"
	_ "github.com/be-io/mesh/plugin/prsim"
	_ "github.com/be-io/mesh/plugin/serve"
)

func init() {
	plugin.Provide(new(proxy))
}

type proxy struct {
}

func (that *proxy) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.ProxyCar, Flags: proxy{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *proxy) Start(ctx context.Context, runtime plugin.Runtime) {
	runtime.Load(plugin.Proxy)
	runtime.Load(plugin.Panel)
	runtime.Load(plugin.PRSIM)
	runtime.Load(plugin.Http)
}

func (that *proxy) Stop(ctx context.Context, runtime plugin.Runtime) {

}
