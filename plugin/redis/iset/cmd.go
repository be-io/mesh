/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
	"strconv"
	"strings"
	"time"
)

// Register
// tensor: ping brpop hmset expire lpush setex
// cube: set,lpush,ping,expire,publish,pubsub,del
// edge: redlock lpop lpush expire get set setnx incr hset/hget hmset/hmget hdel exists
func Register(command Command) {
	commands[command.Name()] = command
}

func Select(cmd string) Command {
	return commands[strings.ToUpper(cmd)]
}

var commands = map[string]Command{}

type Pool interface {
	// Cap returns the capacity of this pool.
	Cap() int
	// IsClosed indicates whether the pool is closed.
	IsClosed() bool
	// Tune changes the capacity of this pool, note that it is noneffective to the infinite or pre-allocation pool.
	Tune(size int)
	// Submit submits a task to this pool.
	//
	// Note that you are allowed to call Pool.Submit() from the current Pool.Submit(),
	// but what calls for special attention is that you will get blocked with the latest
	// Pool.Submit() call once the current Pool runs out of its capacity, and to avoid this,
	// you should instantiate a Pool with ants.WithNonblocking(true).
	Submit(task func()) error
	// Running returns the number of workers currently running.
	Running() int
	// Free returns the number of available goroutines to work, -1 indicates this pool is unlimited.
	Free() int
	// Reboot reboots a closed pool.
	Reboot()
	// Release closes this pool and releases the worker queue.
	Release()
	// ReleaseTimeout is like Release but with a timeout, it waits all workers to exit before timing out.
	ReleaseTimeout(timeout time.Duration) error
	// Waiting returns the number of tasks which are waiting be executed.
	Waiting() int
}

type Command interface {
	Length() int
	Name() string
	Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable)
}

func Exec(ctx context.Context, conn redcon.Conn, cmd redcon.Command, client redis.Cmdable) {
	rc := WithConn(conn)
	if len(cmd.Args) < 1 {
		rc.WriteError("ERR wrong number of arguments")
		return
	}
	sc := string(cmd.Args[0])
	if command := Select(sc); nil != command {
		if len(cmd.Args) != 0 && len(cmd.Args) < command.Length() {
			rc.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", sc))
			return
		}
		command.Serve(ctx, rc, cmd, client)
	} else {
		rc.WriteError(fmt.Sprintf("ERR unknown command '%s'", sc))
	}
}

func IsNil(err error) bool {
	return nil != err && "redis: nil" == err.Error()
}

func ParseInt(cmd redcon.Command, index int, dft int64) (int64, error) {
	if len(cmd.Args) < index+1 {
		return dft, nil
	}
	return strconv.ParseInt(string(cmd.Args[index]), 10, 64)
}

func ParseString(cmd redcon.Command, index int, dft string) (string, error) {
	if len(cmd.Args) < index+1 {
		return dft, nil
	}
	return string(cmd.Args[index]), nil
}

func ParseArgs(cmd redcon.Command, offset int) []string {
	var args []string
	for index, arg := range cmd.Args {
		if index >= offset {
			args = append(args, string(arg))
		}
	}
	return args
}
