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
	"strconv"
	"time"
)

func init() {
	Register(new(Blpop))
}

type Blpop struct {
}

func (that *Blpop) Length() int {
	return 2
}

func (that *Blpop) Name() string {
	return "BLPOP"
}

// Serve
// Array reply: specifically:
//
// A nil multi-bulk when no element could be popped and the timeout expired.
// A two-element multi-bulk with the first element being the name of the key where an element was popped and the second element being the value of the popped element.
func (that *Blpop) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	expire, err := strconv.ParseFloat(string(cmd.Args[len(cmd.Args)-1]), 64)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	var rs []string
	for index, key := range cmd.Args {
		if index < 1 || index >= len(cmd.Args)-1 {
			continue
		}
		r, err := client.BLPop(ctx, time.Second*time.Duration(expire), string(key)).Result()

		if IsNil(err) {
			continue
		}
		if nil != err {
			conn.WriteErr(err)
			return
		}
		rs = append(rs, r...)
	}
	conn.WriteWriteBulkStrings(rs)
}
