/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/types"
)

var ICache = (*Cache)(nil)

// Cache spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Cache interface {

	// Get the value from cache.
	// @MPI("mesh.cache.get")
	Get(ctx context.Context, key string) (*types.CacheEntity, error)

	// Put the value to cache.
	// @MPI("mesh.cache.put")
	Put(ctx context.Context, cell *types.CacheEntity) error

	// Remove the cache value.
	// @MPI("mesh.cache.remove")
	Remove(ctx context.Context, key string) error

	// Incr the cache of expire time.
	// @MPI("mesh.cache.incr")
	Incr(ctx context.Context, key string, value int64) (int64, error)

	// Decr the cache of expire time.
	// @MPI("mesh.cache.decr")
	Decr(ctx context.Context, key string, value int64) (int64, error)

	// Keys the cache key set.
	// @MPI("mesh.cache.keys")
	Keys(ctx context.Context, pattern string) ([]string, error)

	// HGet get value in hash
	// @MPI("mesh.cache.hget")
	HGet(ctx context.Context, key string, name string) (*types.CacheEntity, error)

	// HSet put value in hash
	// @MPI("mesh.cache.hset")
	HSet(ctx context.Context, key string, cell *types.CacheEntity) error

	// HDel put value in hash
	// @MPI("mesh.cache.hdel")
	HDel(ctx context.Context, key string, name string) error

	// HKeys get the hash keys
	// @MPI("mesh.cache.hkeys")
	HKeys(ctx context.Context, key string) ([]string, error)
}
