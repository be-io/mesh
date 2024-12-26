/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import (
	"context"
	"github.com/dgraph-io/badger/v4"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/timshannon/badgerhold/v4"
	"time"
)

func init() {
	macro.Provide(prsim.ISequence, ss)
}

var ss = &prsim.SyncSequence[*session]{
	Syncer:   new(sequence),
	Macro:    &macro.Att{Name: Name},
	Sections: map[string]chan string{}}

type session struct {
	tx *badger.Txn
	ss *badgerhold.Store
}
type sequence struct {
}

func (that *sequence) Tx(ctx context.Context, tx func(session *session) ([]string, error)) ([]string, error) {
	return TXR(ctx, func(ctx context.Context, x *badger.Txn, ss *badgerhold.Store) ([]string, error) {
		return tx(&session{tx: x, ss: ss})
	})
}

func (that *sequence) Incr(ctx context.Context, sx *session, kind string) (*prsim.Section, error) {
	seq := new(SequenceEnt)
	err := sx.ss.TxGet(sx.tx, kind, seq)
	if nil != err && !IsNotFound(err) {
		return nil, cause.Error(err)
	}
	if IsNotFound(err) {
		return nil, nil
	}
	return &prsim.Section{
		Kind:    seq.Kind,
		Min:     seq.Min,
		Max:     seq.Max,
		Size:    seq.Size,
		Length:  seq.Length,
		Version: seq.Version,
	}, nil
}

func (that *sequence) Init(ctx context.Context, sx *session, section *prsim.Section) error {
	err := sx.ss.TxInsert(sx.tx, section.Kind, &SequenceEnt{
		Kind:     section.Kind,
		Min:      section.Min,
		Max:      section.Max,
		Size:     section.Size,
		Length:   section.Length,
		Status:   0,
		Version:  section.Version,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	})
	return cause.Error(err)
}

func (that *sequence) Sync(ctx context.Context, sx *session, section *prsim.Section) error {
	err := sx.ss.TxUpdate(sx.tx, section.Kind, &SequenceEnt{
		Kind:     section.Kind,
		Min:      section.Min + int64(section.Size),
		Max:      section.Max,
		Size:     section.Size,
		Length:   section.Length,
		Status:   0,
		Version:  section.Version,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	})
	return cause.Error(err)
}
