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
	Register(new(Hgetall))
}

type Hgetall struct {
}

func (that *Hgetall) Length() int {
	return 2
}

func (that *Hgetall) Name() string {
	return "HGETALL"
}

// Serve
// Array reply: list of fields and their values stored in the hash, or an empty list when key does not exist.
func (that *Hgetall) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.HGetAll(ctx, string(cmd.Args[1])).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteAny(r)
	}
}
