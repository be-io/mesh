/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"math/rand"
)

func init() {
	var _ prsim.Listener = subscribeListener
	macro.Provide(prsim.IListener, &prsim.ListenerDecorator{Name: "mesh.message.graph", Listener: subscribeListener})
}

var subscribeListener = new(SubscriberListener)

type SubscriberListener struct {
	subscribers map[string]map[string][]*types.Service
}

func (that *SubscriberListener) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.RegistryEventRefresh}
}

func (that *SubscriberListener) Listen(ctx context.Context, event *types.Event) error {
	var registrations types.MetadataRegistrations
	if err := event.TryGetObject(&registrations); nil != err {
		return cause.Error(err)
	}
	subscribers := map[string]map[string][]*types.Service{}
	for _, binding := range registrations.Of(types.METADATA).InferService() {
		if binding.Kind != types.Binding {
			continue
		}
		if nil == subscribers[binding.Namespace] {
			subscribers[binding.Namespace] = map[string][]*types.Service{}
		}
		subscribers[binding.Namespace][binding.Name] = append(subscribers[binding.Namespace][binding.Name], binding)
	}
	that.subscribers = subscribers
	return nil
}

func (that *SubscriberListener) ConsumerAddress(topic *types.Topic) string {
	if nil == that.subscribers[topic.Topic] || nil == that.subscribers[topic.Topic][topic.Code] {
		return ""
	}
	services := that.subscribers[topic.Topic][topic.Code]
	if len(services) < 1 {
		return ""
	}
	if "" == topic.Sets {
		return services[rand.Intn(len(services))].Address
	}
	for _, subscriber := range services {
		if subscriber.Sets == topic.Sets {
			return subscriber.Address
		}
	}
	return ""
}

var _ prsim.Subscriber = new(PRSISubscriber)

// PRSISubscriber
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSISubscriber struct {
}

func (that *PRSISubscriber) Subscribe(ctx context.Context, event *types.Event) error {
	name := fmt.Sprintf("%s.%s", event.Binding.Topic, event.Binding.Code)
	if subscriber, ok := macro.Load(prsim.ISubscriber).Get(name).(prsim.Subscriber); ok {
		return subscriber.Subscribe(ctx, event)
	}
	address := subscribeListener.ConsumerAddress(event.Binding)
	if "" == address {
		return cause.Errorf("No subscriber %s:%s exist. ", event.Binding.Topic, event.Binding.Code)
	}
	mtx := mpc.ContextWith(ctx)
	mtx.SetAttribute(mpc.AddressKey, address)
	mtx.SetAttribute(mpc.RemoteUname, name)
	mtx.GetPrincipals().Push(event.Target)
	defer func() {
		mtx.GetPrincipals().Pop()
	}()
	return aware.RemoteSubscriber.Subscribe(mtx, event)
}
