/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
)

func init() {
	Register(new(Del))
}

type Del struct {
}

func (that *Del) Length() int {
	return 2
}

func (that *Del) Name() string {
	return "DEL"
}

func (that *Del) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	switch client.(type) {
	case redis.Pipeliner:
		that.MasterDel(ctx, conn, cmd, client)
	case *redis.ClusterClient:
		that.EachMasterDel(ctx, conn, cmd, client)
	default:
		that.MasterDel(ctx, conn, cmd, client)
	}
}

func (that *Del) MasterDel(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.Del(ctx, ParseArgs(cmd, 1)...).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt64(r)
	}
}

func (that *Del) EachMasterDel(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	var count int64
	var cause error
	for _, key := range ParseArgs(cmd, 1) {
		if r, err := client.Del(ctx, key).Result(); nil != err {
			log.Error(ctx, err.Error())
			cause = err
		} else {
			count += r
		}
	}
	if 0 == count && nil != cause {
		conn.WriteErr(cause)
	} else {
		conn.WriteInt64(count)
	}
}
