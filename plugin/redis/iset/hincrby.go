/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
	"strconv"
)

func init() {
	Register(new(Hincrby))
}

type Hincrby struct {
}

func (that *Hincrby) Length() int {
	return 4
}

func (that *Hincrby) Name() string {
	return "HINCRBY"
}

// Serve
// Integer reply: the value at field after the increment operation.
func (that *Hincrby) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	incr, err := strconv.ParseInt(string(cmd.Args[3]), 10, 64)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	if r, err := client.HIncrBy(ctx, string(cmd.Args[1]), string(cmd.Args[2]), incr).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt64(r)
	}
}
