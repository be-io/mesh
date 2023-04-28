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
	Register(new(Lrange))
}

// Lrange
// Returns the specified elements of the list stored at key. The offsets start and stop are zero-based indexes,
// with 0 being the first element of the list (the head of the list), 1 being the next element and so on.
// These offsets can also be negative numbers indicating offsets starting at the end of the list. For example,
// -1 is the last element of the list, -2 the penultimate, and so on.
//
// Note that if you have a list of numbers from 0 to 100, LRANGE list 0 10 will return 11 elements, that is,
// the rightmost item is included. This may or may not be consistent with behavior of range-related functions
// in your programming language of choice (think Ruby's Range.new, Array#slice or Python's range() function).
//
// Out of range indexes will not produce an error. If start is larger than the end of the list, an empty list is returned.
// If stop is larger than the actual end of the list, Redis will treat it like the last element of the list.
type Lrange struct {
}

func (that *Lrange) Length() int {
	return 4
}

func (that *Lrange) Name() string {
	return "LRANGE"
}

// Serve
// Array reply: list of elements in the specified range.
func (that *Lrange) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	start, err := ParseInt(cmd, 2, 0)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	stop, err := ParseInt(cmd, 3, 0)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	if r, err := client.LRange(ctx, string(cmd.Args[1]), start, stop).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteWriteBulkStrings(r)
	}
}
