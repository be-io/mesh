/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	_ "github.com/opendatav/mesh/plugin/nsq"
	"regexp"
)

var TopicRegex = regexp.MustCompile(`^[\.a-zA-Z0-9_-]+(#ephemeral)?$`)

var _ prsim.Publisher = new(PSRIPublisher)

// PSRIPublisher
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PSRIPublisher struct {
}

func (that *PSRIPublisher) CheckEvents(ctx context.Context, events []*types.Event) error {
	if events == nil || len(events) == 0 {
		return cause.Errorf("Events cant be empty.")
	}
	for _, event := range events {
		if nil == event.Binding {
			return cause.Errorf("Topic binding cant be empty.")
		}
		if "" == event.Binding.Code || "*" == event.Binding.Code {
			event.Binding.Code = "-"
		}
		if len(event.Binding.Topic) > 64 || len(event.Binding.Topic) < 1 || len(event.Binding.Code) > 64 || len(event.Binding.Code) < 1 {
			return cause.Errorf("Invalid message topic or code %s:%s", event.Binding.Topic, event.Binding.Code)
		}
		if !TopicRegex.MatchString(event.Binding.Topic) || !TopicRegex.MatchString(event.Binding.Code) {
			return cause.Errorf("Invalid message topic or code %s:%s", event.Binding.Topic, event.Binding.Code)
		}
	}
	return nil
}

func (that *PSRIPublisher) Publish(ctx context.Context, events []*types.Event) ([]string, error) {
	if err := that.CheckEvents(ctx, events); nil != err {
		return nil, cause.Error(err)
	}
	return aware.NSQPublisher.Publish(ctx, events)
}

func (that *PSRIPublisher) Broadcast(ctx context.Context, events []*types.Event) ([]string, error) {
	if err := that.CheckEvents(ctx, events); nil != err {
		return nil, cause.Error(err)
	}
	environ, err := aware.LocalNet.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	edges, err := aware.LocalNet.GetRoutes(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	var es []*types.Event
	for _, event := range events {
		for _, edge := range edges {
			target := &types.Principal{
				NodeId: edge.NodeId,
				InstId: edge.InstId,
			}
			es = append(es, &types.Event{
				Version:   event.Version,
				Tid:       event.Tid,
				Eid:       event.Eid,
				Mid:       event.Mid,
				Timestamp: event.Timestamp,
				Source:    event.Source,
				Target:    target,
				Binding:   event.Binding,
				Entity:    event.Entity,
			})
		}
		es = append(es, &types.Event{
			Version:   event.Version,
			Tid:       event.Tid,
			Eid:       event.Eid,
			Mid:       event.Mid,
			Timestamp: event.Timestamp,
			Source:    event.Source,
			Target: &types.Principal{
				NodeId: environ.NodeId,
				InstId: environ.InstId,
			},
			Binding: event.Binding,
			Entity:  event.Entity,
		})
	}
	return aware.NSQPublisher.Publish(ctx, events)
}
