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
	"github.com/be-io/mesh/client/golang/log"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
	"runtime/debug"
	"strings"
)

// PubSub
// Once the client enters the subscribed state it is not supposed to issue any other commands, except for additional
// SUBSCRIBE, SSUBSCRIBE, PSUBSCRIBE, UNSUBSCRIBE, SUNSUBSCRIBE, PUNSUBSCRIBE, PING, RESET and QUIT commands.
type PubSub struct {
	pmessage bool
	client   redis.Cmdable
	pubsub   *redis.PubSub
	conn     redcon.DetachedConn
}

func (that *PubSub) ParseChannels(cmd redcon.Command) []string {
	var channels []string
	for index, v := range cmd.Args {
		if index > 0 {
			channels = append(channels, string(v))
		}
	}
	return channels
}

func (that *PubSub) Close(ctx context.Context) {
	if err := that.pubsub.Close(); nil != err {
		log.Error(ctx, "Redis proxy close pubsub, %s", err.Error())
	}
	if err := that.conn.Close(); nil != err {
		log.Error(ctx, "Redis proxy close conn, %s", err.Error())
	}
}

func (that *PubSub) Read(ctx context.Context) {
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, string(debug.Stack()))
			that.conn.WriteError(fmt.Sprintf("%v", err))
		}
	}()
	for {
		cmd, err := that.conn.ReadCommand()
		if nil != err {
			that.conn.WriteError(err.Error())
			return
		}
		if len(cmd.Args) < 1 {
			continue
		}
		switch strings.ToUpper(string(cmd.Args[0])) {
		case "SUBSCRIBE":
			channels := that.ParseChannels(cmd)
			for index, channel := range channels {
				that.conn.WriteArray(3)
				that.conn.WriteBulkString("subscribe")
				that.conn.WriteBulkString(channel)
				that.conn.WriteInt(index + 1)
			}
			if err = that.pubsub.Subscribe(ctx, channels...); nil != err {
				log.Error(ctx, "Redis proxy subscribe % with error, %s", strings.Join(channels, ","), err.Error())
			}
		case "PSUBSCRIBE":
			channels := that.ParseChannels(cmd)
			for index, channel := range channels {
				that.conn.WriteArray(3)
				that.conn.WriteBulkString("psubscribe")
				that.conn.WriteBulkString(channel)
				that.conn.WriteInt(index + 1)
			}
			if err = that.pubsub.PSubscribe(ctx, channels...); nil != err {
				log.Error(ctx, "Redis proxy psubscribe % with error, %s", strings.Join(channels, ","), err.Error())
			}
		case "UNSUBSCRIBE":
			channels := that.ParseChannels(cmd)
			if err = that.pubsub.Unsubscribe(ctx, channels...); nil != err {
				log.Error(ctx, "Redis proxy unsubscribe % with error, %s", strings.Join(channels, ","), err.Error())
			}
		case "PUNSUBSCRIBE":
			channels := that.ParseChannels(cmd)
			if err = that.pubsub.PUnsubscribe(ctx, channels...); nil != err {
				log.Error(ctx, "Redis proxy punsubscribe % with error, %s", strings.Join(channels, ","), err.Error())
			}
		case "PING":
			that.conn.WriteString("PONG")
		case "QUIT":
			that.conn.WriteString("OK")
			that.Close(ctx)
			return
		default:
			that.conn.WriteError(fmt.Sprintf(
				"ERR Can't execute '%s': only (P)SUBSCRIBE / (P)UNSUBSCRIBE / PING / QUIT are allowed in this context",
				cmd.Args[0]))
		}
	}
}

func (that *PubSub) Write(ctx context.Context) {
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, string(debug.Stack()))
			that.conn.WriteError(fmt.Sprintf("%v", err))
		}
	}()
	for {
		message, err := that.pubsub.ReceiveMessage(ctx)
		if nil != err {
			that.conn.WriteError(err.Error())
			return
		}
		if nil == message {
			continue
		}
		if that.pmessage {
			that.conn.WriteArray(4)
			that.conn.WriteBulkString("pmessage")
			that.conn.WriteBulkString(message.Pattern)
			that.conn.WriteBulkString(message.Channel)
			that.conn.WriteBulkString(message.Payload)
		} else {
			that.conn.WriteArray(3)
			that.conn.WriteBulkString("message")
			that.conn.WriteBulkString(message.Channel)
			that.conn.WriteBulkString(message.Payload)
		}
		if err = that.conn.Flush(); nil != err {
			log.Error(ctx, "Redis proxy flush conn, %s", err.Error())
			that.Close(ctx)
			that.client.Publish(ctx, message.Channel, message.Payload)
			return
		}
	}
}

func (that *PubSub) Start(ctx context.Context) {
	go that.Read(ctx)
	go that.Write(ctx)
}
