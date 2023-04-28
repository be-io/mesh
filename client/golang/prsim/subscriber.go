/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/types"
)

func init() {
	var _ Subscriber = new(SubscriberDecorator)
	var _ macro.SPI = new(SubscriberDecorator)
	var _ Listener = new(ListenerDecorator)
	var _ macro.SPI = new(ListenerDecorator)
}

var IListener = (*Listener)(nil)
var ISubscriber = (*Subscriber)(nil)

// Subscriber spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Subscriber interface {

	// Subscribe the event with {@link com.be.mesh.client.annotate.Bindings} or {@link com.be.mesh.client.annotate.Binding}
	// @MPI("mesh.queue.subscribe")
	Subscribe(ctx context.Context, event *types.Event) error
}

// Listener spi
// @SPI("mesh")
type Listener interface {

	// Bindings of the listener
	macro.Bindings

	// Listen the event.
	// Listen function can't be blocked.
	Listen(ctx context.Context, event *types.Event) error
}

type SubscriberDecorator struct {
	Name       string
	Subscriber Subscriber
}

func (that *SubscriberDecorator) Att() *macro.Att {
	return &macro.Att{Name: that.Name}
}

func (that *SubscriberDecorator) Subscribe(ctx context.Context, event *types.Event) error {
	return that.Subscriber.Subscribe(ctx, event)
}

type ListenerDecorator struct {
	Name     string
	Listener Listener
}

func (that *ListenerDecorator) Att() *macro.Att {
	return &macro.Att{Name: that.Name}
}

func (that *ListenerDecorator) Btt() []*macro.Btt {
	return that.Listener.Btt()
}

func (that *ListenerDecorator) Listen(ctx context.Context, event *types.Event) error {
	return that.Listener.Listen(ctx, event)
}
