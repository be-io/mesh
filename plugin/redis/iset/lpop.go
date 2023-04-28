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
	Register(new(Lpop))
}

type Lpop struct {
}

func (that *Lpop) Length() int {
	return 2
}

func (that *Lpop) Name() string {
	return "LPOP"
}

// Serve
// When called without the count argument:
//
// Bulk string reply: the value of the first element, or nil when key does not exist.
//
// When called with the count argument:
//
// Array reply: list of popped elements, or nil when key does not exist.
func (that *Lpop) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.LPop(ctx, string(cmd.Args[1])).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteBulkString(r)
	}
}
