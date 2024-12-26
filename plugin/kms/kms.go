/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package kms

import (
	"context"
	"github.com/opendatav/mesh/client/golang/plugin"
)

func init() {
	plugin.Provide(new(kms))
}

type kms struct {
}

func (that *kms) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.KMS, Flags: kms{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *kms) Start(ctx context.Context, runtime plugin.Runtime) {
}

func (that *kms) Stop(ctx context.Context, runtime plugin.Runtime) {
}
