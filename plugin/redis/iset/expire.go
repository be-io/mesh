/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
	"strconv"
	"time"
)

func init() {
	Register(new(Expire))
}

type Expire struct {
}

func (that *Expire) Length() int {
	return 2
}

func (that *Expire) Name() string {
	return "EXPIRE"
}

func (that *Expire) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	expire, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	if r, err := client.Expire(ctx, string(cmd.Args[1]), time.Second*time.Duration(expire)).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt(tool.Ternary(r, 1, 0))
	}
}
