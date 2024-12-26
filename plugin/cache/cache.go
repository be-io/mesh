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
	"github.com/buraksezer/olric/config"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/plugin"
)

func init() {
	plugin.Provide(new(mache))
}

type mache struct {
	Address string       `json:"plugin.redis.address" dft:"0.0.0.0:6379" usage:"Redis access layer listen address. "`
	Home    string       `json:"plugin.redis.home" dft:"${MESH_HOME}/mesh/redis" usage:"Redis work home dir. "`
	Servers string       `json:"plugin.redis.servers" dft:"" usage:"Redis proxy backend servers. "`
	CDB     *olric.Olric `json:"-"`
}

func (that *mache) Ptt() *plugin.Ptt {
	return &plugin.Ptt{
		Name:  plugin.Cache,
		Flags: mache{},
		Create: func() plugin.Plugin {
			return that
		},
		Priority: 0,
		WaitAny:  false,
	}
}

func (that *mache) Start(ctx context.Context, runtime plugin.Runtime) {
	// config.New returns a new config.Config with sane defaults. Available values for env:
	// local, lan, wan
	conf := config.New("local")
	// Callback function. It's called when this node is ready to accept connections.
	conf.Started = func() {
		log.Info(ctx, "Cluster cache is ready to accept connections.")
	}
	cdb, err := olric.New(conf)
	if err != nil {
		log.Error(ctx, "Failed to create cluster cache object: %v", err)
		return
	}
	that.CDB = cdb
	runtime.Submit(func() {
		// Call Start at background. It's a blocker call.
		if err := that.CDB.Start(); nil != err {
			log.Error(ctx, "Failed to start cluster cache: %v", err)
		}
	})
}

func (that *mache) Stop(ctx context.Context, runtime plugin.Runtime) {
	if nil != that.CDB {
		log.Catch(that.CDB.Shutdown(ctx))
	}
}
