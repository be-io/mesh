/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/types"
)

var IDataHouse = (*DataHouse)(nil)

// DataHouse spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type DataHouse interface {

	// Writes
	// @MPI("mesh.dh.writes")
	Writes(ctx context.Context, docs []*types.Document) error

	// Write
	// @MPI("mesh.dh.write")
	Write(ctx context.Context, doc *types.Document) error

	// Read
	// @MPI("mesh.dh.read")
	Read(ctx context.Context, index *types.Paging) (*types.Page[any], error)

	// Indies
	// @MPI("mesh.dh.indies")
	Indies(ctx context.Context, index *types.Paging) (*types.Page[any], error)

	// Tables
	// @MPI("mesh.dh.tables")
	Tables(ctx context.Context, index *types.Paging) (*types.Page[any], error)
}
