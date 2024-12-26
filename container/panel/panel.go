/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package panel

import (
	"context"
	"github.com/opendatav/mesh/client/golang/plugin"
	// --
	_ "github.com/opendatav/mesh/plugin/panel"
	_ "github.com/opendatav/mesh/plugin/serve"
)

func init() {
	plugin.Provide(new(panel))

}

type panel struct {
}

func (that *panel) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.PanelCar, Flags: panel{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *panel) Start(ctx context.Context, runtime plugin.Runtime) {
	runtime.Load(plugin.Panel)
}

func (that *panel) Stop(ctx context.Context, runtime plugin.Runtime) {
}
