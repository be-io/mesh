/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

var _ prsim.KMS = new(PRSIKMS)

// PRSIKMS
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSIKMS struct {
}

func (that *PRSIKMS) Reset(ctx context.Context, env *types.Environ) error {
	return aware.KMS.Reset(ctx, env)
}

func (that *PRSIKMS) Environ(ctx context.Context) (*types.Environ, error) {
	return aware.KMS.Environ(ctx)
}

func (that *PRSIKMS) List(ctx context.Context, cno string) ([]*types.Keys, error) {
	return aware.KMS.List(ctx, cno)
}

func (that *PRSIKMS) ApplyRoot(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error) {
	return aware.KMS.ApplyRoot(ctx, csr)
}

func (that *PRSIKMS) ApplyIssue(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error) {
	return aware.KMS.ApplyIssue(ctx, csr)
}
