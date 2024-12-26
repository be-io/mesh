/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cache

import (
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/patrickmn/go-cache"
	"strings"
	"sync"
	"time"
)

func init() {
	var _ prsim.Cache = new(UnSharedCache)
	macro.Provide(prsim.ICache, new(UnSharedCache))
}

const Unshared = "unshared"

type UnSharedCache struct {
	cache *cache.Cache
	lock  sync.Mutex
}

func (that *UnSharedCache) Att() *macro.Att {
	return &macro.Att{Name: Unshared}
}

func (that *UnSharedCache) Store(ctx context.Context) *cache.Cache {
	if nil == that.cache {
		that.lock.Lock()
		defer that.lock.Unlock()
		that.cache = cache.New(time.Second*30, time.Minute*5)
	}
	return that.cache
}

func (that *UnSharedCache) Get(ctx context.Context, key string) (*types.CacheEntity, error) {
	value, exist := that.Store(ctx).Get(key)
	if !exist {
		return nil, nil
	}
	if v, ok := value.(*types.CacheEntity); ok {
		return v, nil
	}
	return nil, cause.Errorf("Cant recognized cached value of %s", key)
}

func (that *UnSharedCache) Put(ctx context.Context, cell *types.CacheEntity) error {
	that.Store(ctx).Set(cell.Key, cell, time.Duration(cell.Duration*time.Millisecond.Nanoseconds()))
	return nil
}

func (that *UnSharedCache) Remove(ctx context.Context, key string) error {
	that.Store(ctx).Delete(key)
	return nil
}

func (that *UnSharedCache) Incr(ctx context.Context, key string, value int64) (int64, error) {
	return that.Store(ctx).IncrementInt64(key, value)
}

func (that *UnSharedCache) Decr(ctx context.Context, key string, value int64) (int64, error) {
	return that.Store(ctx).DecrementInt64(key, value)
}

func (that *UnSharedCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	for key, _ := range that.Store(ctx).Items() {
		keys = append(keys, key)
	}
	return keys, nil
}

func (that *UnSharedCache) HGet(ctx context.Context, key string, name string) (*types.CacheEntity, error) {
	return that.Get(ctx, fmt.Sprintf("%s.%s", key, name))
}

func (that *UnSharedCache) HSet(ctx context.Context, key string, cell *types.CacheEntity) error {
	cell.Key = fmt.Sprintf("%s.%s", key, cell.Key)
	return cause.Error(that.Put(ctx, cell))
}

func (that *UnSharedCache) HDel(ctx context.Context, key string, name string) error {
	return cause.Error(that.Remove(ctx, fmt.Sprintf("%s.%s", key, name)))
}

func (that *UnSharedCache) HKeys(ctx context.Context, key string) ([]string, error) {
	dotKey := fmt.Sprintf("%s.", key)
	names, err := that.Keys(ctx, dotKey)
	if nil != err {
		return nil, cause.Error(err)
	}
	var keys []string
	for _, name := range names {
		if strings.Index(name, dotKey) == 0 {
			keys = append(keys, strings.Replace(name, dotKey, "", 1))
		}
	}
	return keys, nil
}
