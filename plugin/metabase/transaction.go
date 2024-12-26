/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/opendatav/mesh/client/golang/cause"
	bh "github.com/timshannon/badgerhold/v4"
)

func TXR[T any](ctx context.Context, fn func(ctx context.Context, tx *badger.Txn, ss *bh.Store) (T, error)) (r T, err error) {
	s, err := metabase.GetStore(ctx)
	if nil != err {
		return r, cause.Error(err)
	}
	err = s.Badger().Update(func(txn *badger.Txn) error {
		r, err = fn(ctx, txn, metabase.Store)
		return err
	})
	return r, err
}

func TX(ctx context.Context, fn func(ctx context.Context, tx *badger.Txn, ss *bh.Store) error) (err error) {
	s, err := metabase.GetStore(ctx)
	if nil != err {
		return cause.Error(err)
	}
	return s.Badger().Update(func(txn *badger.Txn) error {
		return fn(ctx, txn, metabase.Store)
	})
}

func TR[T any](ctx context.Context, fn func(ctx context.Context, ss *bh.Store) (T, error)) (r T, err error) {
	s, err := metabase.GetStore(ctx)
	if nil != err {
		return r, cause.Error(err)
	}
	return fn(ctx, s)
}

func IsDuplicateKeyError(err error) bool {
	if errMySQL, ok := err.(*mysql.MySQLError); ok {
		switch errMySQL.Number {
		case 1062:
			return true
		}
		return false
	}
	if errPQ, ok := err.(*pq.Error); ok {
		return "23505" == string(errPQ.Code)
	}
	return errors.Is(err, bh.ErrKeyExists)
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, bh.ErrNotFound) || errors.Is(err, badger.ErrKeyNotFound)
}
