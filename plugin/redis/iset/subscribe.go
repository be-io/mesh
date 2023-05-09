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
	Register(new(Subscribe))
}

// Subscribe
// Once the client enters the subscribed state it is not supposed to issue any other commands, except for additional
// SUBSCRIBE, SSUBSCRIBE, PSUBSCRIBE, UNSUBSCRIBE, SUNSUBSCRIBE, PUNSUBSCRIBE, PING, RESET and QUIT commands.
type Subscribe struct {
}

func (that *Subscribe) Length() int {
	return 2
}

func (that *Subscribe) Name() string {
	return "SUBSCRIBE"
}

func (that *Subscribe) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
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
		conn.WriteBulkString("subscribe")
		conn.WriteBulkString(channel)
		conn.WriteInt(index + 1)
	}
	pubsub := &PubSub{
		pmessage: false,
		client:   client,
		pubsub:   rc.Subscribe(ctx, channels...),
		conn:     conn.Detach(),
	}
	pubsub.Start(ctx)
}
