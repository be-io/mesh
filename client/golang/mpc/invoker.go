/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
)

func init() {
	var _ Invoker = new(InvokerDecorator)
	var _ macro.SPI = new(InvokerDecorator)
}

var IInvoker = (*Invoker)(nil)

type Invoker interface {

	// Invoke the next invoker.
	Invoke(ctx context.Context, invocation Invocation) (interface{}, error)
}

type InvokerFn func(ctx context.Context, invocation Invocation) (interface{}, error)

func (that InvokerFn) Invoke(ctx context.Context, invocation Invocation) (interface{}, error) {
	return that(ctx, invocation)
}

type InvokerDecorator struct {
	Name    string
	Invoker Invoker
}

func (that *InvokerDecorator) Att() *macro.Att {
	return &macro.Att{Name: that.Name}
}

func (that *InvokerDecorator) Invoke(ctx context.Context, invocation Invocation) (interface{}, error) {
	return that.Invoker.Invoke(ctx, invocation)
}
