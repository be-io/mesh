/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/opendatav/mesh/plugin/redis/iset"
	"time"
)

func init() {
	var _ prsim.Cache = new(redisCache)
	macro.Provide(prsim.ICache, new(redisCache))
}

const Name = "redis"

type redisCache struct {
}

func (that *redisCache) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *redisCache) Decode(value string) (*types.CacheEntity, error) {
	var entity types.CacheEntity
	if _, err := aware.Codec.DecodeString(value, &entity); nil != err {
		return nil, cause.Error(err)
	}
	return &entity, nil
}

func (that *redisCache) Get(ctx context.Context, key string) (*types.CacheEntity, error) {
	client, err := Ref(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	r, err := client.Get(ctx, key).Result()
	if iset.IsNil(err) {
		return nil, nil
	}
	if nil != err {
		return nil, cause.Error(err)
	}
	return that.Decode(r)
}

func (that *redisCache) Put(ctx context.Context, cell *types.CacheEntity) error {
	client, err := Ref(ctx)
	if nil != err {
		return cause.Error(err)
	}
	txt, err := aware.Codec.EncodeString(cell)
	if nil != err {
		return cause.Error(err)
	}
	_, err = client.Set(ctx, cell.Key, txt, time.Duration(cell.Duration*time.Millisecond.Nanoseconds())).Result()
	return cause.Error(err)
}

func (that *redisCache) Remove(ctx context.Context, key string) error {
	client, err := Ref(ctx)
	if nil != err {
		return cause.Error(err)
	}
	_, err = client.Del(ctx, key).Result()
	return cause.Error(err)
}

func (that *redisCache) Incr(ctx context.Context, key string, value int64) (int64, error) {
	client, err := Ref(ctx)
	if nil != err {
		return 0, cause.Error(err)
	}
	r, err := client.IncrBy(ctx, key, value).Result()
	return r, cause.Error(err)
}

func (that *redisCache) Decr(ctx context.Context, key string, value int64) (int64, error) {
	client, err := Ref(ctx)
	if nil != err {
		return 0, cause.Error(err)
	}
	r, err := client.DecrBy(ctx, key, value).Result()
	return r, cause.Error(err)
}

func (that *redisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	client, err := Ref(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	return client.Keys(ctx, pattern).Result()
}

func (that *redisCache) HGet(ctx context.Context, key string, name string) (*types.CacheEntity, error) {
	client, err := Ref(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	r, err := client.HGet(ctx, key, name).Result()
	if iset.IsNil(err) {
		return nil, nil
	}
	if nil != err {
		return nil, cause.Error(err)
	}
	return that.Decode(r)
}

func (that *redisCache) HSet(ctx context.Context, key string, cell *types.CacheEntity) error {
	client, err := Ref(ctx)
	if nil != err {
		return cause.Error(err)
	}
	txt, err := aware.Codec.EncodeString(cell)
	if nil != err {
		return cause.Error(err)
	}
	_, err = client.HSet(ctx, key, cell.Key, txt).Result()
	return cause.Error(err)
}

func (that *redisCache) HDel(ctx context.Context, key string, name string) error {
	client, err := Ref(ctx)
	if nil != err {
		return cause.Error(err)
	}
	_, err = client.HDel(ctx, key, name).Result()
	return cause.Error(err)
}

func (that *redisCache) HKeys(ctx context.Context, key string) ([]string, error) {
	client, err := Ref(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	return client.HKeys(ctx, key).Result()
}
