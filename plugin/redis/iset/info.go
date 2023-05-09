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
	Register(new(Info))
}

type Info struct {
}

func (that *Info) Length() int {
	return 1
}

func (that *Info) Name() string {
	return "INFO"
}

// Serve
// Bulk string reply: as a collection of text lines.
//
// Lines can contain a section name (starting with a # character) or a property. All the properties are in the form of field:value terminated by \r\n.
func (that *Info) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	var sections []string
	for index, section := range cmd.Args {
		if index > 0 {
			sections = append(sections, string(section))
		}
	}
	if r, err := client.Info(ctx, sections...).Result(); nil != err {
		conn.WriteErr(err)
	} else {
		conn.WriteBulkString(r)
	}
}
