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
	"sort"
)

const (
	// PROVIDER Flag.
	PROVIDER = "PROVIDER"
	// CONSUMER Flag.
	CONSUMER = "CONSUMER"
)

var IFilter = (*Filter)(nil)

type Filter interface {

	// Invoke the next filter.
	// invoker    Service invoker.
	// invocation Service invocation
	Invoke(ctx context.Context, invoker Invoker, invocation Invocation) (interface{}, error)

	macro.SPI
}

type filters []Filter

func (that filters) Len() int {
	return len(that)
}

func (that filters) Less(x, y int) bool {
	return that[x].Att().Priority < that[y].Att().Priority
}

func (that filters) Swap(x, y int) {
	tmp := that[y]
	that[y] = that[x]
	that[x] = tmp
}

// Composite the filter spi providers as an invoker.
// invoker Delegate invoker.
// pattern Filter pattern.
func composite(invoker Invoker, pattern string) Invoker {
	var filters filters
	x := macro.Load(IFilter).List()
	for _, filter := range x {
		filter, ok := filter.(Filter)
		if !ok || filter.Att().Pattern != pattern {
			continue
		}
		filters = append(filters, filter)
	}
	sort.Sort(filters)
	return composite0(invoker, filters)
}

// Composite the filter spi providers as a invoker.
// invoker Delegate invoker.
// filters Filter list.
func composite0(invoker Invoker, filters []Filter) Invoker {
	last := invoker
	for i := len(filters) - 1; i >= 0; i-- {
		filter := filters[i]
		next := last
		var x InvokerFn = func(ctx context.Context, invocation Invocation) (interface{}, error) {
			return filter.Invoke(ctx, next, invocation)
		}
		last = x
	}
	return last
}
