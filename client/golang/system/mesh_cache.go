/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/schema"
	"github.com/be-io/mesh/client/golang/types"
	"time"
)

func init() {
	var _ prsim.Cache = new(MeshCache)
	macro.Provide(prsim.ICache, new(MeshCache))
}

type MeshCache struct {
}

func (that *MeshCache) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *MeshCache) Get(ctx context.Context, key string) (*types.CacheEntity, error) {
	return aware.Cache.Get(ctx, key)
}

func (that *MeshCache) Put(ctx context.Context, cell *types.CacheEntity) error {
	return aware.Cache.Put(ctx, cell)
}

func (that *MeshCache) Remove(ctx context.Context, key string) error {
	return aware.Cache.Remove(ctx, key)
}

func (that *MeshCache) Incr(ctx context.Context, key string, value int64) (int64, error) {
	return aware.Cache.Incr(ctx, key, value)
}

func (that *MeshCache) Decr(ctx context.Context, key string, value int64) (int64, error) {
	return aware.Cache.Decr(ctx, key, value)
}

func (that *MeshCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return aware.Cache.Keys(ctx, pattern)
}

func (that *MeshCache) HGet(ctx context.Context, key string, name string) (*types.CacheEntity, error) {
	return aware.Cache.HGet(ctx, key, name)
}

func (that *MeshCache) HSet(ctx context.Context, key string, cell *types.CacheEntity) error {
	return aware.Cache.HSet(ctx, key, cell)
}

func (that *MeshCache) HDel(ctx context.Context, key string, name string) error {
	return aware.Cache.HDel(ctx, key, name)
}

func (that *MeshCache) HKeys(ctx context.Context, key string) ([]string, error) {
	return aware.Cache.HKeys(ctx, key)
}

func GetWithCache(ctx context.Context, cache prsim.Cache, key string, kind interface{}) error {
	cell, err := cache.Get(ctx, key)
	if nil != err {
		return cause.Error(err)
	}
	if nil == cell || nil == cell.Entity {
		return nil
	}
	return cell.Entity.TryReadObject(kind)
}

func PutWithCache(ctx context.Context, cache prsim.Cache, key string, value interface{}, duration time.Duration) error {
	buffer, err := aware.JSON.Encode(value)
	if nil != err {
		return cause.Error(err)
	}
	definition, err := schema.Runtime.Define(ctx, value)
	if nil != err {
		return cause.Error(err)
	}
	cell := &types.CacheEntity{
		Version: types.CacheVersion,
		Entity: &types.Entity{
			Codec:  codec.JSON,
			Schema: definition,
			Buffer: buffer.Bytes(),
		},
		Timestamp: time.Now().UnixMilli(),
		Duration:  duration.Milliseconds(),
		Key:       key,
	}
	return cache.Put(ctx, cell)
}
