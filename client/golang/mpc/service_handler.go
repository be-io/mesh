/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ Invoker = new(ServiceHandler)
}

type ServiceHandler struct {
	inspector macro.Inspector
}

func (that *ServiceHandler) Invoke(ctx context.Context, invocation Invocation) (interface{}, error) {
	if service := that.determineMPS(invocation); nil != service {
		return service.Invoke(ctx, invocation)
	}
	if service := that.determineBinding(invocation); nil != service {
		return service.Invoke(ctx, invocation)
	}
	if service, ok := macro.Load(IInvoker).StartWith(invocation.GetURN().Name).(Invoker); ok && nil != service {
		return service.Invoke(ctx, invocation)
	}
	return nil, NoService(ctx, invocation.GetURN())
}

func (that *ServiceHandler) determineMPS(invocation Invocation) Invoker {
	if nil != invocation.GetInspector().GetMPI() {
		if rtt := invocation.GetInspector().GetMPI().Rtt(); nil != rtt {
			if service, ok := macro.Load(IInvoker).Get(rtt.Name).(Invoker); ok {
				return service
			}
		}
	}
	if nil == invocation.GetInspector().GetMPS() {
		return nil
	}
	stt := invocation.GetInspector().GetMPS().Stt()
	if nil == stt {
		return nil
	}
	if service, ok := macro.Load(IInvoker).Get(stt.Name).(Invoker); ok {
		return service
	}
	return nil
}

func (that *ServiceHandler) determineBinding(invocation Invocation) Invoker {
	if len(invocation.GetInspector().GetBinding()) < 1 {
		return nil
	}
	for _, binding := range invocation.GetInspector().GetBinding() {
		if nil == binding {
			continue
		}
		btt := binding.Btt()
		if nil == btt {
			continue
		}
		if service, ok := macro.Load(IInvoker).Get(fmt.Sprintf("%s.%s", btt.Topic, btt.Code)).(Invoker); ok {
			return service
		}
	}
	return nil
}

func NoService(ctx context.Context, urn *types.URN) error {
	if nil != urn && "" != urn.Name {
		return cause.Errorcf(cause.NoService, "No service named %s", urn.Name)
	}
	mtx := ContextWith(ctx)
	return cause.Errorcf(cause.NoService, "No service named %s", mtx.GetUrn())
}
