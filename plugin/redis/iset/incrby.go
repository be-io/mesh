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
	Register(new(IncrBy))
}

// IncrBy
// Increments the number stored at key by increment. If the key does not exist, it is set to 0 before performing
// the operation. An error is returned if the key contains a value of the wrong type or contains a string that
// can not be represented as integer. This operation is limited to 64 bit signed integers.
//
// See INCR for extra information on increment/decrement operations.
type IncrBy struct {
}

func (that *IncrBy) Length() int {
	return 3
}

func (that *IncrBy) Name() string {
	return "INCRBY"
}

// Serve Integer reply: the value of key after the increment
func (that *IncrBy) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	incr, err := ParseInt(cmd, 2, 0)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	if r, err := client.IncrBy(ctx, string(cmd.Args[1]), incr).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt64(r)
	}
}
