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
	Register(new(Execs))
}

type Execs struct {
}

func (that *Execs) Length() int {
	return 1
}

func (that *Execs) Name() string {
	return "EXEC"
}

// Serve
// Executes all previously queued commands in a transaction and restores the connection state to normal.
//
// When using WATCH, EXEC will execute commands only if the watched keys were not modified, allowing for a check-and-set mechanism.
//
// Return
// Array reply: each element being the reply to each of the commands in the atomic transaction.
//
// When using WATCH, EXEC can return a Null reply if the execution was aborted.
func (that *Execs) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	pip, ok := client.(redis.Pipeliner)
	if !ok {
		conn.WriteError("Unknown connection state.")
		return
	}
	rs, err := pip.Exec(ctx)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	conn.WriteArray(len(rs))
	for _, r := range rs {
		if nil != r.Err() {
			conn.WriteErr(r.Err())
			continue
		}
		switch d := r.(type) {
		case *redis.IntCmd:
			conn.WriteInt64(d.Val())
		case *redis.StringCmd:
			conn.WriteBulkString(d.Val())
		case *redis.StatusCmd:
			conn.WriteBulkString(d.Val())
		}
	}
}
