/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import (
	"bytes"
	"context"
	"github.com/dgraph-io/badger/v4"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/timshannon/badgerhold/v4"
	"time"
)

func init() {
	var _ prsim.KV = omegaKv
	macro.Provide(prsim.IKV, omegaKv)
}

var omegaKv = new(kv)

type kv struct {
}

func (that *kv) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *kv) Get(ctx context.Context, key string) (*types.Entity, error) {
	return TR(ctx, func(ctx context.Context, ss *badgerhold.Store) (*types.Entity, error) {
		set := new(KVEnt)
		err := ss.Get(key, set)
		if nil != err && !IsNotFound(err) {
			return nil, cause.Error(err)
		}
		if IsNotFound(err) {
			return new(types.Entity).AsEmpty(), nil
		}
		var entity types.Entity
		if _, err = aware.Codec.Decode(bytes.NewBufferString(set.Value), &entity); nil != err {
			return nil, cause.Error(err)
		}
		return &entity, nil
	})
}

func (that *kv) Put(ctx context.Context, key string, value *types.Entity) error {
	return TX(ctx, func(ctx context.Context, tx *badger.Txn, ss *badgerhold.Store) error {
		entity, err := aware.Codec.EncodeString(value)
		if nil != err {
			return cause.Error(err)
		}
		mtx := mpc.ContextWith(ctx)
		operator := tool.Anyone(mtx.GetAttachments()["omega.user.username"], mtx.GetAttachments()["omega.inst.name"])
		exist := new(KVEnt)
		err = ss.TxGet(tx, key, exist)
		if nil != err && !IsNotFound(err) {
			return cause.Error(err)
		}
		err = ss.TxUpsert(tx, key, &KVEnt{
			Key:      key,
			Value:    entity,
			CreateAt: tool.Anyone(exist.CreateAt, time.Now()),
			UpdateAt: time.Now(),
			CreateBy: tool.Anyone(exist.CreateBy, operator),
			UpdateBy: operator,
		})
		return cause.Error(err)
	})
}

func (that *kv) Remove(ctx context.Context, key string) error {
	return TX(ctx, func(ctx context.Context, tx *badger.Txn, ss *badgerhold.Store) error {
		return cause.Error(ss.TxDelete(tx, key, new(KVEnt)))
	})
}

func (that *kv) Keys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	return TR(ctx, func(ctx context.Context, ss *badgerhold.Store) ([]string, error) {
		err := ss.Badger().View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			opts.Prefix = []byte(pattern)
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				keys = append(keys, string(k))
			}
			return nil
		})
		return keys, cause.Error(err)
	})
}

func (that *kv) Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Entry], error) {
	return TR(ctx, func(ctx context.Context, ss *badgerhold.Store) (*types.Page[*types.Entry], error) {
		con := new(badgerhold.Query).SortBy("Key").Skip(int(index.Index * index.Limit)).Limit(int(index.Limit))
		key := index.Factor["key"]
		if nil != key && "" != key {
			con.And("Key").Contains(key)
		}
		total, err := ss.Count(new(KVEnt), con)
		if nil != err {
			return nil, cause.Error(err)
		}
		sid := index.SID
		if "" == sid {
			sid = tool.NextID()
		}
		var kes []*KVEnt
		if err = ss.Find(&kes, con); nil != err {
			return nil, cause.Error(err)
		}
		var kvs []*types.Entry
		for _, ke := range kes {
			var entity types.Entry
			if _, err = aware.Codec.Decode(bytes.NewBufferString(ke.Value), &entity); nil != err {
				return nil, cause.Error(err)
			}
			kvs = append(kvs, &entity)
		}
		return types.Reset(index, int64(total), kvs), nil
	})
}
