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
	Register(new(Setex))
}

type Setex struct {
}

func (that *Setex) Length() int {
	return 4
}

func (that *Setex) Name() string {
	return "SETEX"
}

// Serve
// Simple string reply
func (that *Setex) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	expire, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	if r, err := client.SetEx(ctx, string(cmd.Args[1]), cmd.Args[3], time.Second*time.Duration(expire)).Result(); nil != err {
		conn.WriteError(err.Error())
	} else {
		conn.WriteString(r)
	}
}
