/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/plugin"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/opendatav/mesh/plugin/redis/iset"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
	"runtime/debug"
	"strings"
	"time"
)

func init() {
	var _ plugin.Plugin = ral
	var _ prsim.Listener = ral
	plugin.Provide(ral)
	macro.Provide(prsim.IListener, ral)
}

func Ref(ctx context.Context) (redis.UniversalClient, error) {
	if nil == ral.Backends || nil == ral.Backends.Any(ctx) {
		return nil, cause.Errorable(cause.NoProvider)
	}
	return ral.Backends.Any(ctx), nil
}

var checkPing = []*macro.Btt{{Topic: "mesh.plugin.redis.health", Code: "ping"}}
var ral = new(redisAccessLayer)

type redisOption struct {
	Address string        `json:"plugin.redis.address" dft:"0.0.0.0:6379" usage:"Redis access layer listen address. "`
	Home    string        `json:"plugin.redis.home" dft:"${MESH_HOME}/mesh/redis" usage:"Redis work home dir. "`
	Servers string        `json:"plugin.redis.servers" dft:"redis://127.0.0.1:6379" usage:"Redis proxy backend servers. "`
	Slaves  string        `json:"plugin.redis.slaves" dft:"" usage:"Redis proxy backend slave servers. "`
	Timeout time.Duration `json:"plugin.redis.timeout" dft:"24m" usage:"Redis proxy idle connection timeout in seconds. "`
}

type redisAccessLayer struct {
	Server   *redcon.Server       `json:"-"`
	Backends *MasterSlave         `json:"-"`
	Mini     *miniredis.Miniredis `json:"-"`
}

func (that *redisAccessLayer) Close() error {
	return nil
}

func (that *redisAccessLayer) Ptt() *plugin.Ptt {
	return &plugin.Ptt{
		Name:  plugin.Redis,
		Flags: redisOption{},
		Create: func() plugin.Plugin {
			return that
		},
	}
}

func (that *redisAccessLayer) IsInternal(ctx context.Context, option *redisOption) bool {
	u, err := types.FormatURL(option.Servers)
	if nil != err {
		log.Error(ctx, err.Error())
		return true
	}
	return strings.Contains(option.Address, fmt.Sprintf(":%s", u.Port())) && strings.Contains(u.Hostname(), "127.0.0.1")
}

func (that *redisAccessLayer) Start(ctx context.Context, runtime plugin.Runtime) {
	option := new(redisOption)
	err := runtime.Parse(option)
	if nil != err {
		log.Error(ctx, "Redis proxy dont startup because options is invalid, %s. ", option.Servers, err.Error())
		return
	}
	that.Backends, err = that.NewMasterSalves(ctx, option.Servers, option.Slaves)
	if nil != err {
		log.Error(ctx, "Redis proxy dont startup, %s", err.Error())
		return
	}
	if "" == option.Servers || that.IsInternal(ctx, option) {
		log.Info(ctx, "Redis proxy dont startup because has no backend servers. ")
		runtime.Submit(func() {
			that.Mini = miniredis.NewMiniRedis()
			if err = that.Mini.StartAddr(option.Address); nil != err {
				log.Error(ctx, "Redis server dont startup because %s. ", err.Error())
			}
		})
		return
	}
	that.Server = redcon.NewServerNetwork("tcp", option.Address, that.Serve, that.Accept, that.Closed)
	that.Server.SetIdleClose(option.Timeout)
	runtime.Submit(func() {
		if err = that.Server.ListenAndServe(); nil != err {
			log.Error(ctx, "Redis proxy with unexpected error, %s", err.Error())
		}
	})
	topic := &types.Topic{Topic: checkPing[0].Topic, Code: checkPing[0].Code}
	if _, err = aware.Scheduler.Period(ctx, time.Second*10, topic); nil != err {
		log.Error(ctx, "Redis proxy active ping ticker, %s", err.Error())
	}
	log.Warn(ctx, "Redis proxy with %s. ", option.Servers)
}

func (that *redisAccessLayer) Stop(ctx context.Context, runtime plugin.Runtime) {
	if nil != that.Server {
		log.Catch(that.Server.Close())
	}
	if nil != that.Backends {
		if nil != that.Backends.master {
			log.Catch(that.Backends.master.reference.Close())
		}
		for _, slave := range that.Backends.salves {
			log.Catch(slave.reference.Close())
		}
	}
	if nil != that.Mini {
		that.Mini.Close()
	}
}

func (that *redisAccessLayer) Serve(rc redcon.Conn, cmd redcon.Command) {
	ctx := mpc.Context()
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, "%v", err)
			log.Error(ctx, string(debug.Stack()))
			rc.WriteError(fmt.Sprintf("ERR %v", err))
		}
	}()
	if nil == that.Backends {
		rc.WriteError("ERR no backend service")
		return
	}
	service := that.Backends.Any(ctx)
	if nil == service {
		rc.WriteError("ERR no available backend service")
		return
	}
	iset.Exec(ctx, rc, cmd, service)
}

func (that *redisAccessLayer) Accept(conn redcon.Conn) bool {
	// Use this function to accept or deny the connection.
	// log.Info0("accept: %s", conn.RemoteAddr())
	return true
}
func (that *redisAccessLayer) Closed(conn redcon.Conn, err error) {
	// This is called when the connection has been closed
	// log.Info0("closed: %s, err: %v", conn.RemoteAddr(), err)
}

func (that *redisAccessLayer) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.redis.health"}
}

func (that *redisAccessLayer) Btt() []*macro.Btt {
	return checkPing
}

func (that *redisAccessLayer) Listen(ctx context.Context, event *types.Event) error {
	if nil != that.Backends {
		that.Backends.Pings(ctx)
	}
	return nil
}
