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
	var _ prsim.KV = new(MeshKV)
	macro.Provide(prsim.IKV, new(MeshKV))
}

type MeshKV struct {
}

func (that *MeshKV) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *MeshKV) Get(ctx context.Context, key string) (*types.Entity, error) {
	return aware.KV.Get(ctx, key)
}

func (that *MeshKV) Put(ctx context.Context, key string, value *types.Entity) error {
	return aware.KV.Put(ctx, key, value)
}

func (that *MeshKV) Remove(ctx context.Context, key string) error {
	return aware.KV.Remove(ctx, key)
}

func (that *MeshKV) Keys(ctx context.Context, pattern string) ([]string, error) {
	return aware.KV.Keys(ctx, pattern)
}

func (that *MeshKV) Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Entry], error) {
	return aware.KV.Index(ctx, index)
}
