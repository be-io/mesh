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
	Register(new(Ping))
}

type Ping struct {
}

func (that *Ping) Length() int {
	return 1
}

func (that *Ping) Name() string {
	return "PING"
}

func (that *Ping) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.Ping(ctx).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteString(r)
	}

}
