/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/types"
)

var IWorkflow = (*Workflow)(nil)

// Workflow
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Workflow interface {

	// Mass workflow work in group.
	// Return workflow code
	// @MPI("mesh.workflow.mass")
	Mass(ctx context.Context, group *types.WorkGroup) (string, error)

	// Groups page workflow review groups.
	// @MPI("mesh.workflow.groups")
	Groups(ctx context.Context, index *types.Paging) (*types.Page[*types.WorkGroup], error)

	// Compile workflow in engine.
	// Return workflow code
	// @MPI("mesh.workflow.compile")
	Compile(ctx context.Context, chart *types.WorkChart) (string, error)

	// Index workflows.
	// @MPI("mesh.workflow.index")
	Index(ctx context.Context, index *types.Paging) (*types.Page[*types.WorkChart], error)

	// Submit workflow.
	// Return workflow instance code
	// @MPI("mesh.workflow.submit")
	Submit(ctx context.Context, intent *types.WorkIntent) (string, error)

	// Take action on workflow instance.
	// @MPI("mesh.workflow.take")
	Take(ctx context.Context, vertex *types.WorkVertex) error

	// Routines infer workflow instance as routine.
	// @MPI("mesh.workflow.routines")
	Routines(ctx context.Context, index *types.Paging) (*types.Page[*types.WorkVertex], error)
}
