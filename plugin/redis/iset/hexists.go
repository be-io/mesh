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
)

func init() {
	Register(new(HExists))
}

type HExists struct {
}

func (that *HExists) Length() int {
	return 2
}

func (that *HExists) Name() string {
	return "HEXISTS"
}

// Serve
// Integer reply, specifically:
//
// 1 if the hash contains field.
// 0 if the hash does not contain field, or key does not exist.
func (that *HExists) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.HExists(ctx, string(cmd.Args[1]), string(cmd.Args[2])).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteInt(tool.Ternary(r, 1, 0))
	}
}
