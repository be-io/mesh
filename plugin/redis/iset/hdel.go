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
	Register(new(Hdel))
}

type Hdel struct {
}

func (that *Hdel) Length() int {
	return 2
}

func (that *Hdel) Name() string {
	return "HDEL"
}

// Serve
// Integer reply: the number of fields that were removed from the hash, not including specified but non existing fields.
func (that *Hdel) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	var kv []string
	for index, v := range cmd.Args {
		if index > 1 {
			kv = append(kv, string(v))
		}
	}
	if r, err := client.HDel(ctx, string(cmd.Args[1]), kv...).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt64(r)
	}
}
