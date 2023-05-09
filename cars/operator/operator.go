/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package operator

import (
	"context"
	"github.com/be-io/mesh/client/golang/plugin"
)

func init() {
	plugin.Provide(new(operator))
}

type operator struct {
}

func (that *operator) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.WorkloadCar, Flags: operator{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *operator) Start(ctx context.Context, runtime plugin.Runtime) {
}

func (that *operator) Stop(ctx context.Context, runtime plugin.Runtime) {
}
