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
	"github.com/golang/groupcache"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"sync"
)

func init() {
	var _ prsim.Cache = new(SharedCache)
	macro.Provide(prsim.ICache, new(SharedCache))
}

const Shared = "shared"

type SharedCache struct {
	group *groupcache.Group
	lock  sync.Mutex
}

func (that *SharedCache) Att() *macro.Att {
	return &macro.Att{Name: Shared}
}

func (that *SharedCache) Store(ctx context.Context, name string) *groupcache.Group {
	if nil == that.group {
		that.lock.Lock()
		defer that.lock.Unlock()
		getter := groupcache.GetterFunc(func(ctx context.Context, key string, dest groupcache.Sink) error {
			return cause.Errorf("")
		})
		that.group = groupcache.NewGroup(fmt.Sprintf("mesh-shared-%s-cache", name), 1<<20, getter)
	}
	return that.group
}

func (that *SharedCache) Get(ctx context.Context, key string) (*types.CacheEntity, error) {
	view := &groupcache.ByteView{}
	err := that.group.Get(ctx, key, groupcache.ByteViewSink(view))
	return nil, err
}

func (that *SharedCache) Put(ctx context.Context, cell *types.CacheEntity) error {
	return cause.NotImplementError()
}

func (that *SharedCache) Remove(ctx context.Context, key string) error {
	return cause.NotImplementError()
}

func (that *SharedCache) Incr(ctx context.Context, key string, value int64) (int64, error) {
	return 0, cause.NotImplementError()
}

func (that *SharedCache) Decr(ctx context.Context, key string, value int64) (int64, error) {
	return 0, cause.NotImplementError()
}

func (that *SharedCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return nil, cause.NotImplementError()
}

func (that *SharedCache) HGet(ctx context.Context, key string, name string) (*types.CacheEntity, error) {
	return nil, cause.NotImplementError()
}

func (that *SharedCache) HSet(ctx context.Context, key string, cell *types.CacheEntity) error {
	return cause.NotImplementError()
}

func (that *SharedCache) HDel(ctx context.Context, key string, name string) error {
	return cause.NotImplementError()
}

func (that *SharedCache) HKeys(ctx context.Context, key string) ([]string, error) {
	return nil, cause.NotImplementError()
}
