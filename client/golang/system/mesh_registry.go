/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ prsim.Registry = new(systemRegistry)
	macro.Provide(prsim.IRegistry, new(systemRegistry))
}

type systemRegistry struct {
}

func (that *systemRegistry) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *systemRegistry) Register(ctx context.Context, registration *types.Registration[any]) error {
	return aware.Registry.Register(ctx, registration)
}

func (that *systemRegistry) Registers(ctx context.Context, registrations []*types.Registration[any]) error {
	return aware.Registry.Registers(ctx, registrations)
}

func (that *systemRegistry) Unregister(ctx context.Context, registration *types.Registration[any]) error {
	return aware.Registry.Unregister(ctx, registration)
}

func (that *systemRegistry) Export(ctx context.Context, kind string) ([]*types.Registration[any], error) {
	return aware.Registry.Export(ctx, kind)
}
