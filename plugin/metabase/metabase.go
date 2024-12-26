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
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/plugin"
	"github.com/timshannon/badgerhold/v4"
)

func init() {
	plugin.Provide(metabase)
}

const Name = "metabase"

var metabase = new(meta)

type meta struct {
	DSN   string            `json:"dsn" dft:"" usage:"Data source name, Ex mysql://user:password@127.0.0.1:3306/dbname"`
	Home  string            `json:"meta.home" dft:"${MESH_HOME}/mesh/metadata/" usage:"Data storage home dir"`
	Store *badgerhold.Store `json:"-"`
}

func (that *meta) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: Name, Flags: meta{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *meta) Start(ctx context.Context, runtime plugin.Runtime) {
	log.Catch(runtime.Parse(that))
	store, err := badgerhold.Open(badgerhold.Options{
		Encoder:          badgerhold.DefaultEncode,
		Decoder:          badgerhold.DefaultDecode,
		SequenceBandwith: badgerhold.DefaultOptions.SequenceBandwith,
		Options:          badger.DefaultOptions(that.Home),
	})
	if nil != err {
		log.Error(ctx, err.Error())
		return
	}
	that.Store = store
}

func (that *meta) Stop(ctx context.Context, runtime plugin.Runtime) {
	if nil != that.Store {
		log.Catch(that.Store.Close())
	}
}

func (that *meta) GetStore(ctx context.Context) (*badgerhold.Store, error) {
	if nil == metabase.Store {
		return nil, cause.Errorable(cause.StartPending)
	}
	return metabase.Store, nil
}
