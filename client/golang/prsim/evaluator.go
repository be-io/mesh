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

var IEvaluator = (*Evaluator)(nil)

// Evaluator
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Evaluator interface {

	// Compile the named rule.
	// @MPI("mesh.eval.compile")
	Compile(ctx context.Context, script *types.Script) (string, error)

	// Exec the script with name.
	// @MPI("mesh.eval.exec")
	Exec(ctx context.Context, code string, args map[string]string, dft string) (string, error)

	// Dump the scripts.
	// @MPI("mesh.eval.dump")
	Dump(ctx context.Context, feature map[string]string) ([]*types.Script, error)

	// Index the scripts.
	// @MPI("mesh.eval.index")
	Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Script], error)
}
