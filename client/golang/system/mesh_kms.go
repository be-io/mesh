/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ prsim.KMS = new(systemKMS)
	macro.Provide(prsim.IKMS, new(systemKMS))
}

type systemKMS struct {
}

func (that *systemKMS) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *systemKMS) Reset(ctx context.Context, env *types.Environ) error {
	return aware.KMS.Reset(ctx, env)
}

func (that *systemKMS) Environ(ctx context.Context) (*types.Environ, error) {
	return aware.KMS.Environ(ctx)
}

func (that *systemKMS) List(ctx context.Context, cno string) ([]*types.Keys, error) {
	return aware.KMS.List(ctx, cno)
}

func (that *systemKMS) ApplyRoot(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error) {
	return aware.KMS.ApplyRoot(ctx, csr)
}

func (that *systemKMS) ApplyIssue(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error) {
	return aware.KMS.ApplyIssue(ctx, csr)
}
