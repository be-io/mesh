/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
	"time"
)

func init() {
	Register(new(Setnx))
}

type Setnx struct {
}

func (that *Setnx) Length() int {
	return 2
}

func (that *Setnx) Name() string {
	return "SETNX"
}

// Serve
// Integer reply, specifically:
//
// 1 if the key was set
// 0 if the key was not set
func (that *Setnx) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.SetNX(ctx, string(cmd.Args[1]), cmd.Args[2], time.Hour*24).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt(tool.Ternary(r, 1, 0))
	}
}
