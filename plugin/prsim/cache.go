/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	_ "github.com/be-io/mesh/plugin/cache"
)

var _ prsim.Cache = new(PRSICache)

// PRSICache
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSICache struct {
}

func (that *PRSICache) Ref() prsim.Cache {
	if macro.RCache.Enable() {
		return aware.Redis
	}
	return aware.Local
}

func (that *PRSICache) Get(ctx context.Context, key string) (*types.CacheEntity, error) {
	return that.Ref().Get(ctx, key)
}

func (that *PRSICache) Put(ctx context.Context, cell *types.CacheEntity) error {
	return that.Ref().Put(ctx, cell)
}

func (that *PRSICache) Remove(ctx context.Context, key string) error {
	return that.Ref().Remove(ctx, key)
}

func (that *PRSICache) Incr(ctx context.Context, key string, value int64) (int64, error) {
	return that.Ref().Incr(ctx, key, value)
}

func (that *PRSICache) Decr(ctx context.Context, key string, value int64) (int64, error) {
	return that.Ref().Decr(ctx, key, value)
}

func (that *PRSICache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return that.Ref().Keys(ctx, pattern)
}

func (that *PRSICache) HGet(ctx context.Context, key string, name string) (*types.CacheEntity, error) {
	return that.Ref().HGet(ctx, key, name)
}

func (that *PRSICache) HSet(ctx context.Context, key string, cell *types.CacheEntity) error {
	return that.Ref().HSet(ctx, key, cell)
}

func (that *PRSICache) HDel(ctx context.Context, key string, name string) error {
	return that.Ref().HDel(ctx, key, name)
}

func (that *PRSICache) HKeys(ctx context.Context, key string) ([]string, error) {
	return that.Ref().HKeys(ctx, key)
}
