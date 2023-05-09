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
	Register(new(PSubscribe))
}

// PSubscribe
// Subscribes the client to the given patterns.
// Supported glob-style patterns:
//
// h?llo subscribes to hello, hallo and hxllo
// h*llo subscribes to hllo and heeeello
// h[ae]llo subscribes to hello and hallo, but not hillo
// Use \ to escape special characters if you want to match them verbatim.
type PSubscribe struct {
}

func (that *PSubscribe) Length() int {
	return 2
}

func (that *PSubscribe) Name() string {
	return "PSUBSCRIBE"
}

func (that *PSubscribe) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	rc, ok := client.(redis.UniversalClient)
	if !ok {
		conn.WriteError("Unknown connection state.")
		return
	}
	var channels []string
	for index, v := range cmd.Args {
		if index > 0 {
			channels = append(channels, string(v))
		}
	}
	for index, channel := range channels {
		conn.WriteArray(3)
		conn.WriteBulkString("psubscribe")
		conn.WriteBulkString(channel)
		conn.WriteInt(index + 1)
	}
	pubsub := &PubSub{
		pmessage: true,
		client:   client,
		pubsub:   rc.PSubscribe(ctx, channels...),
		conn:     conn.Detach(),
	}
	pubsub.Start(ctx)
}
