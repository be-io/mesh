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
	"time"
)

func init() {
	Register(new(Set))
}

type Set struct {
}

func (that *Set) Length() int {
	return 3
}

func (that *Set) Name() string {
	return "SET"
}

// Serve
// Simple string reply: OK if SET was executed correctly.
//
// Null reply: (nil) if the SET operation was not performed because the user specified the NX or XX option but the condition was not met.
//
// If the command is issued with the GET option, the above does not apply. It will instead reply as follows, regardless if the SET was actually performed:
//
// Bulk string reply: the old string value stored at key.
//
// Null reply: (nil) if the key did not exist.
func (that *Set) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	if r, err := client.Set(ctx, string(cmd.Args[1]), cmd.Args[2], time.Hour*24).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteString(r)
	}
}
