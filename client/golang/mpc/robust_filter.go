/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"math"
	"strings"
)

func init() {
	var _ Filter = new(robustFilter)
	macro.Provide(IFilter, new(robustFilter))
}

type robustFilter struct {
}

func (that *robustFilter) Att() *macro.Att {
	return &macro.Att{Name: "robust", Pattern: CONSUMER, Priority: math.MaxInt}
}

func (that *robustFilter) Invoke(ctx context.Context, invoker Invoker, invocation Invocation) (interface{}, error) {
	mtx := ContextWith(ctx)
	schema := invocation.GetExecution().Schema()
	retries := tool.Ternary(schema.GetRetries() > 3, schema.GetRetries(), 3)
	for i := 0; i < retries-1; i++ {
		ret, err := invoker.Invoke(ctx, invocation)
		if nil == err {
			if that.isHealthCheckURN(invocation.GetURN()) {
				if server, ok := mtx.GetAttribute(AddressKey).(string); ok {
					tool.Address.Get().Available(server, true)
				}
			}
			return ret, err
		}
		if !that.isHealthCheckURN(invocation.GetURN()) {
			if that.retryCheck(err) {
				continue
			}
			return ret, err
		}
		if that.isAvailable(err) {
			if server, ok := mtx.GetAttribute(AddressKey).(string); ok {
				tool.Address.Get().Available(server, false)
			}
		}
		return ret, err
	}
	return invoker.Invoke(ctx, invocation)
}

func (that *robustFilter) isHealthCheckURN(urn *types.URN) bool {
	return strings.Index(urn.Name, "mesh.registry") == 0
}

func (that *robustFilter) isAvailable(err error) bool {
	if c, ok := err.(cause.Codeable); !ok {
		return false
	} else {
		return cause.Timeout.Code == c.GetCode() || cause.NetUnavailable.Code == c.GetCode()
	}
}

func (that *robustFilter) retryCheck(err error) bool {
	return false
}
