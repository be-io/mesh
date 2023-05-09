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
	Register(new(Publish))
}

type Publish struct {
}

func (that *Publish) Length() int {
	return 3
}

func (that *Publish) Name() string {
	return "PUBLISH"
}

func (that *Publish) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.Publish(ctx, string(cmd.Args[1]), cmd.Args[2]).Result(); nil != err {
		conn.WriteError(err.Error())
	} else {
		conn.WriteInt64(r)
	}
}
