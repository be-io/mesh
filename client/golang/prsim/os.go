/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"github.com/be-io/mesh/client/golang/types"
	"golang.org/x/net/context"
)

var IOperateSystem = (*OperateSystem)(nil)

// OperateSystem spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type OperateSystem interface {

	// Install apps.
	// @MPI("mesh.os.install")
	Install(ctx context.Context, chart *types.OSCharts) error

	// Uninstall apps.
	// @MPI("mesh.os.uninstall")
	Uninstall(ctx context.Context, chart *types.OSCharts) error

	// Index apps.
	// @MPI("mesh.os.index")
	Index(ctx context.Context) (*types.Page[*types.OSCharts], error)
}
