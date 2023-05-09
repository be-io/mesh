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
	Register(new(Hget))
}

type Hget struct {
}

func (that *Hget) Length() int {
	return 3
}

func (that *Hget) Name() string {
	return "HGET"
}

// Serve
// Bulk string reply: the value associated with field, or nil when field is not present in the hash or key does not exist.
func (that *Hget) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.HGet(ctx, string(cmd.Args[1]), string(cmd.Args[2])).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteBulkString(r)
	}
}
