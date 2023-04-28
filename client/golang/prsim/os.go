/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
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

	// Cluster definition.
	// @MPI("mesh.os.cluster")
	Cluster(ctx context.Context, cluster *types.Cluster) error

	// Clusters definition.
	// @MPI("mesh.os.cluster")
	Clusters(ctx context.Context, index *types.Paging) (*types.Page, error)

	// Workspace definition.
	// @MPI("mesh.os.workspace")
	Workspace(ctx context.Context, cluster *types.Workspace) error

	// Workspaces definition.
	// @MPI("mesh.os.workspaces")
	Workspaces(ctx context.Context, index *types.Paging) (*types.Page, error)

	// Install apps.
	// @MPI("mesh.os.install")
	Install(ctx context.Context, chart *types.OSCharts) error

	// Uninstall apps.
	// @MPI("mesh.os.uninstall")
	Uninstall(ctx context.Context, chart *types.OSCharts) error

	// Index apps.
	// @MPI("mesh.os.index")
	Index(ctx context.Context, index *types.Paging) (*types.Page, error)

	// Operations index.
	// @MPI("mesh.os.operations")
	Operations(ctx context.Context, index *types.Paging) (*types.Page, error)
}
