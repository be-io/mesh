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
	Register(new(Hmset))
}

type Hmset struct {
}

func (that *Hmset) Length() int {
	return 1
}

func (that *Hmset) Name() string {
	return "HMSET"
}

// Serve
// Simple string reply
func (that *Hmset) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	var kv []interface{}
	for index, v := range cmd.Args {
		if index > 1 {
			kv = append(kv, v)
		}
	}
	if _, err := client.HMSet(ctx, string(cmd.Args[1]), kv...).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteString("OK")
	}
}
