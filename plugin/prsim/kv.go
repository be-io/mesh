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

var _ prsim.KV = new(PRSIKV)

// PRSIKV
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSIKV struct {
	cache prsim.Cache
}

func (that *PRSIKV) Get(ctx context.Context, key string) (*types.Entity, error) {
	return aware.KV.Get(ctx, key)
}

func (that *PRSIKV) Put(ctx context.Context, key string, value *types.Entity) error {
	return aware.KV.Put(ctx, key, value)
}

func (that *PRSIKV) Remove(ctx context.Context, key string) error {
	return aware.KV.Remove(ctx, key)
}

func (that *PRSIKV) Keys(ctx context.Context, pattern string) ([]string, error) {
	return aware.KV.Keys(ctx, pattern)
}

func (that *PRSIKV) Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Entry], error) {
	return aware.KV.Index(ctx, index)
}
