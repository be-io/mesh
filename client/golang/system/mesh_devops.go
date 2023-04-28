/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
)

func init() {
	var _ prsim.Devops = new(MeshDevops)
	macro.Provide(prsim.IDevops, new(MeshDevops))
}

// MeshDevops
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type MeshDevops struct {
}

func (that *MeshDevops) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *MeshDevops) Distribute(ctx context.Context, option *types.DistributeOption) (string, error) {
	return aware.Devops.Distribute(ctx, option)
}

func (that *MeshDevops) Transform(ctx context.Context, option *types.TransformOption) (string, error) {
	return aware.Devops.Transform(ctx, option)
}
