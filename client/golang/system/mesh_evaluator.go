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
	var _ prsim.Evaluator = new(MeshEvaluator)
	macro.Provide(prsim.IEvaluator, new(MeshEvaluator))
}

type MeshEvaluator struct {
}

func (that *MeshEvaluator) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *MeshEvaluator) Compile(ctx context.Context, script *types.Script) (string, error) {
	return aware.Evaluator.Compile(ctx, script)
}

func (that *MeshEvaluator) Exec(ctx context.Context, code string, args map[string]string, dft string) (string, error) {
	return aware.Evaluator.Exec(ctx, code, args, dft)
}

func (that *MeshEvaluator) Dump(ctx context.Context, feature map[string]string) ([]*types.Script, error) {
	return aware.Evaluator.Dump(ctx, feature)
}

func (that *MeshEvaluator) Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Script], error) {
	return aware.Evaluator.Index(ctx, index)
}
