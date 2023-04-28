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
)

func init() {
	Register(new(Hincrbyfloat))
}

type Hincrbyfloat struct {
}

func (that *Hincrbyfloat) Length() int {
	return 2
}

func (that *Hincrbyfloat) Name() string {
	return "HINCRBYFLOAT"
}

// Serve
// Bulk string reply: the value of field after the increment.
func (that *Hincrbyfloat) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	incr, err := strconv.ParseFloat(string(cmd.Args[3]), 64)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	if r, err := client.HIncrByFloat(ctx, string(cmd.Args[1]), string(cmd.Args[2]), incr).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteBulkString(strconv.FormatFloat(r, 'f', -1, 64))
	}
}
