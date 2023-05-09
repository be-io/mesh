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
	Register(new(Decr))
}

type Decr struct {
}

func (that *Decr) Length() int {
	return 2
}

func (that *Decr) Name() string {
	return "DECR"
}

func (that *Decr) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.Decr(ctx, string(cmd.Args[1])).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt64(r)
	}
}
