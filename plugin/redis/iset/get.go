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
)

func init() {
	Register(new(Get))
}

type Get struct {
}

func (that *Get) Length() int {
	return 2
}

func (that *Get) Name() string {
	return "GET"
}

// Serve
// Bulk string reply: the value of key, or nil when key does not exist.
func (that *Get) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.Get(ctx, string(cmd.Args[1])).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteBulkString(r)
	}
}
