/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cache

import (
	"context"
	"github.com/buraksezer/olric"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var cc = new(clusterCache)
	var _ prsim.Cache = cc
	var _ prsim.RuntimeHook = cc
	macro.Provide(prsim.ICache, cc)
	macro.Provide(prsim.IRuntimeHook, cc)
}

const Cluster = "cluster"

type clusterCache struct {
	cache *olric.Olric
}

func (that *clusterCache) Att() *macro.Att {
	return &macro.Att{Name: Cluster}
}

func (that *clusterCache) Start(ctx context.Context, runtime prsim.Runtime) error {
	return that.Refresh(ctx, runtime)
}

func (that *clusterCache) Stop(ctx context.Context, runtime prsim.Runtime) error {
	if nil != that.cache {
		if err := that.cache.Shutdown(ctx); nil != err {
			log.Error(ctx, "Stop cluster cache, %s", err.Error())
		}
	}
	return nil
}

func (that *clusterCache) Refresh(ctx context.Context, runtime prsim.Runtime) error {
	//if nil != that.cache {
	//	return nil
	//}
	//// config.New returns a new config.Config with sane defaults. Available values for env:
	//// local, lan, wan
	//conf := config.New("wan")
	//cc, err := olric.New(conf)
	//if nil != err {
	//	log.Error(ctx, "Init cluster cache, %s", err.Error())
	//	return nil
	//}
	//that.cache = cc
	//runtime.Submit(func() {
	//	// Call Start at background. It's a blocker call.
	//	if err = that.cache.Start(); nil != err {
	//		log.Error(ctx, "Start cluster cache, %s", err.Error())
	//	}
	//})
	return nil
}

func (that *clusterCache) Get(ctx context.Context, key string) (*types.CacheEntity, error) {
	return nil, nil
}

func (that *clusterCache) Put(ctx context.Context, cell *types.CacheEntity) error {
	return nil
}

func (that *clusterCache) Remove(ctx context.Context, key string) error {
	return nil
}

func (that *clusterCache) Incr(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (that *clusterCache) Decr(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (that *clusterCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return nil, nil
}

func (that *clusterCache) HGet(ctx context.Context, key string, name string) (*types.CacheEntity, error) {
	return nil, cause.NotImplementError()
}

func (that *clusterCache) HSet(ctx context.Context, key string, cell *types.CacheEntity) error {
	return cause.NotImplementError()
}

func (that *clusterCache) HDel(ctx context.Context, key string, name string) error {
	return cause.NotImplementError()
}

func (that *clusterCache) HKeys(ctx context.Context, key string) ([]string, error) {
	return nil, cause.NotImplementError()
}
