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
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/redis/go-redis/v9"
	"strings"
)

func (that *redisAccessLayer) NewMasterSalves(ctx context.Context, master string, slaves string) (*MasterSlave, error) {
	masterServer, err := that.NewClient(ctx, master)
	if nil != err {
		return nil, cause.Error(err)
	}
	var slaveServers []*StatefulServer
	for _, slave := range strings.Split(slaves, ";") {
		if "" == strings.TrimSpace(slave) {
			continue
		}
		slaveServer, err := that.NewClient(ctx, slave)
		if nil != err {
			log.Error(ctx, err.Error())
			continue
		}
		slaveServers = append(slaveServers, &StatefulServer{
			addr:      slave,
			reference: slaveServer,
			available: true,
		})
	}
	return &MasterSlave{
		master: &StatefulServer{
			addr:      master,
			reference: masterServer,
			available: true,
		},
		salves: slaveServers}, nil
}

type StatefulServer struct {
	addr      string
	reference redis.UniversalClient
	available bool
}

type MasterSlave struct {
	master          *StatefulServer
	salves          []*StatefulServer
	availableSlaves []*StatefulServer
}

func (that *MasterSlave) Any(ctx context.Context) redis.UniversalClient {
	if nil != that.master && that.master.available {
		return that.master.reference
	}
	slaves := that.availableSlaves
	if len(slaves) > 0 {
		slave := slaves[0]
		log.Warn(ctx, "Redis proxy switch to slave backend")
		return slave.reference
	}
	if nil == that.master {
		return nil
	}
	return that.master.reference
}

func (that *MasterSlave) Pings(ctx context.Context) {
	if nil != that.master {
		if _, err := that.master.reference.Ping(ctx).Result(); nil != err {
			log.Warn(ctx, "Redis proxy ping check master %s, %s", that.master.addr, err.Error())
			that.master.available = false
		} else {
			that.master.available = true
		}
	}
	for _, slave := range that.salves {
		if _, err := slave.reference.Ping(ctx).Result(); nil != err {
			log.Warn(ctx, "Redis proxy ping check slave %s, %s", slave.addr, err.Error())
			slave.available = false
		} else {
			slave.available = true
		}
	}
	var availableSlaves []*StatefulServer
	for _, slave := range that.salves {
		if slave.available {
			availableSlaves = append(availableSlaves, slave)
		}
	}
	that.availableSlaves = availableSlaves

}
