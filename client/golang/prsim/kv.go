/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/types"
)

var IKV = (*KV)(nil)

var ManProxyRouteMetadata = &macro.Btt{Topic: "mesh.kv.proxy.route.metadata", Code: "*"}

// KV spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type KV interface {

	// Get the value from kv store.
	// @MPI("mesh.kv.get")
	Get(ctx context.Context, key string) (*types.Entity, error)

	// Put the value to kv store.
	// @MPI("mesh.kv.put")
	Put(ctx context.Context, key string, value *types.Entity) error

	// Remove the kv store.
	// @MPI("mesh.kv.remove")
	Remove(ctx context.Context, key string) error

	// Keys with the pattern of kv store.
	// @MPI("mesh.kv.keys")
	Keys(ctx context.Context, pattern string) ([]string, error)

	// Index the kv for webui
	// @MPI("mesh.kv.index")
	Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Entry], error)
}

func GetKV(ctx context.Context, kv KV, key string, ptr any) error {
	ent, err := kv.Get(ctx, key)
	if nil != err {
		return cause.Error(err)
	}
	if nil == ent || !ent.Present() {
		return nil
	}
	return ent.TryReadObject(ptr)
}

func PutKV(ctx context.Context, kv KV, key string, value any) error {
	ent, err := new(types.Entity).Wrap(value)
	if nil != err {
		return cause.Error(err)
	}
	return kv.Put(ctx, key, ent)
}
