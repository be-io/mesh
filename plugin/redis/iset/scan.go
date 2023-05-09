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
	Register(new(Scan))
}

type Scan struct {
}

func (that *Scan) Length() int {
	return 2
}

func (that *Scan) Name() string {
	return "SCAN"
}

// Serve
// SCAN, SSCAN, HSCAN and ZSCAN return a two elements multi-bulk reply, where the first element is a string representing an unsigned 64 bit number (the cursor), and the second element is a multi-bulk with an array of elements.
//
// SCAN array of elements is a list of keys.
// SSCAN array of elements is a list of Set members.
// HSCAN array of elements contain two elements, a field and a value, for every returned element of the Hash.
// ZSCAN array of elements contain two elements, a member and its associated score, for every returned element of the sorted set.
func (that *Scan) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	cursor, err := strconv.ParseUint(string(cmd.Args[1]), 10, 64)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	count, err := ParseInt(cmd, 5, 0)
	if nil != err {
		conn.WriteErr(err)
		return
	}
	match, err := ParseString(cmd, 3, "")
	if ks, c, err := client.Scan(ctx, cursor, match, count).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteArray(2)
		conn.WriteBulkString(strconv.FormatUint(c, 10))
		conn.WriteWriteBulkStrings(ks)
	}
}

func (that *Scan) ParseCount(cmd redcon.Command) (int64, error) {
	if len(cmd.Args) < 5 {
		return 0, nil
	}
	return strconv.ParseInt(string(cmd.Args[5]), 10, 64)
}
