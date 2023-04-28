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
	Register(new(Hmget))
}

type Hmget struct {
}

func (that *Hmget) Length() int {
	return 2
}

func (that *Hmget) Name() string {
	return "HMGET"
}

// Serve
// Array reply: list of values associated with the given fields, in the same order as they are requested.
func (that *Hmget) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	var kv []string
	for index, v := range cmd.Args {
		if index > 1 {
			kv = append(kv, string(v))
		}
	}
	if r, err := client.HMGet(ctx, string(cmd.Args[1]), kv...).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteAnyArray(r)
	}
}
