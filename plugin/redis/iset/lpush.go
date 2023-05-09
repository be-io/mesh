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
	Register(new(Lpush))
}

type Lpush struct {
}

func (that *Lpush) Length() int {
	return 2
}

func (that *Lpush) Name() string {
	return "LPUSH"
}

// Serve
// Integer reply: the length of the list after the push operations.
func (that *Lpush) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	var kv []interface{}
	for index, v := range cmd.Args {
		if index > 1 {
			kv = append(kv, v)
		}
	}
	if r, err := client.LPush(ctx, string(cmd.Args[1]), kv...).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt64(r)
	}
}
