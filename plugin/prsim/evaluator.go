/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
)

var _ prsim.Evaluator = new(PSRIEvaluator)

// PSRIEvaluator
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PSRIEvaluator struct {
}

func (that *PSRIEvaluator) Compile(ctx context.Context, script *types.Script) (string, error) {
	return aware.Evaluator.Compile(ctx, script)
}

func (that *PSRIEvaluator) Exec(ctx context.Context, code string, args map[string]string, dft string) (string, error) {
	return aware.Evaluator.Exec(ctx, code, args, dft)
}

func (that *PSRIEvaluator) Dump(ctx context.Context, feature map[string]string) ([]*types.Script, error) {
	return aware.Evaluator.Dump(ctx, feature)
}

func (that *PSRIEvaluator) Index(ctx context.Context, index *types.Paging) (*types.Page, error) {
	return aware.Evaluator.Index(ctx, index)
}
