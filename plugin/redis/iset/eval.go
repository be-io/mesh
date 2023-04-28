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
	Register(new(Eval))
}

type Eval struct {
}

func (that *Eval) Length() int {
	return 3
}

func (that *Eval) Name() string {
	return "EVAL"
}

func (that *Eval) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	ns, err := strconv.Atoi(string(cmd.Args[2]))
	if nil != err {
		conn.WriteErr(err)
		return
	}
	var keys []string
	var args []interface{}
	for index, key := range cmd.Args {
		if index < 3 {
			continue
		}
		if index < 3+ns {
			keys = append(keys, string(key))
		} else {
			args = append(args, key)
		}
	}
	if r, err := client.Eval(ctx, string(cmd.Args[1]), keys, args...).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteAny(r)
	}
}
