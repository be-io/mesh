/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package nsqio

import (
	"bytes"
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/nsqio/nsq/nsqd"
)

func init() {
	var _ prsim.Publisher = publisher
}

const Publisher = "nsq"

var (
	publisher     = new(nsqPublisher)
	exchangeTopic = &types.Topic{Topic: "mesh.queue.broker.exchange", Code: AllEventCode}
)

type nsqPublisher struct {
}

func (that *nsqPublisher) Att() *macro.Att {
	return &macro.Att{Name: Publisher}
}

func (that *nsqPublisher) Publish(ctx context.Context, events []*types.Event) ([]string, error) {
	environ, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	var eventIds []string
	for _, event := range events {
		if "" == event.Binding.Code || "*" == event.Binding.Code {
			event.Binding.Code = AllEventCode
		}
		if tool.IsLocalEnv(environ, event.Target.InstId, event.Target.NodeId) {
			if eventId, err := that.Push(ctx, event.Binding, event); nil != err {
				return eventIds, cause.Error(err)
			} else {
				log.Debug(ctx, "Queue %s:%s[%s] published in queue. ", event.Binding.Topic, event.Binding.Code, eventId)
				eventIds = append(eventIds, eventId)
			}
			continue
		}
		if eventId, err := that.Push(ctx, exchangeTopic, event); nil != err {
			return eventIds, cause.Error(err)
		} else {
			log.Debug(ctx, "Queue %s:%s[%s] published in broker. ", event.Binding.Topic, event.Binding.Code, eventId)
			eventIds = append(eventIds, eventId)
		}
	}
	return eventIds, nil
}

func (that *nsqPublisher) Broadcast(ctx context.Context, events []*types.Event) ([]string, error) {
	return that.Publish(ctx, events)
}

func (that *nsqPublisher) Push(ctx context.Context, binding *types.Topic, event *types.Event) (string, error) {
	if nil == daemon.core {
		return "", cause.Errorf("Queue dont startup. ")
	}
	topic := daemon.core.GetTopic(binding.Topic)
	channel := topic.GetChannel(binding.Code)
	if nil == event.Entity.Buffer {
		return "", cause.Errorf("Empty message buff %s:%s", binding.Topic, binding.Code)
	}
	buff, err := aware.Codec.Encode(event)
	if nil != err {
		return "", cause.Error(err)
	}
	copied := make([]byte, buff.Len())
	copy(copied, buff.Bytes())
	message := nsqd.NewMessage(topic.GenerateID(), copied)

	messageId := bytes.Buffer{}
	for _, bits := range message.ID {
		messageId.WriteByte(bits)
	}

	err = channel.PutMessage(message)
	if nil != err {
		return "", cause.Error(err)
	}
	return messageId.String(), nil
}
