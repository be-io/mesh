/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/types"
)

var IDevops = (*Devops)(nil)

// Devops spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Devops interface {

	// Distribute the service sdk.
	// @MPI("mesh.devops.distribute")
	Distribute(ctx context.Context, option *types.DistributeOption) (string, error)

	// Transform standard schema to another schema.
	// @MPI("mesh.devops.transform")
	Transform(ctx context.Context, option *types.TransformOption) (string, error)
}
