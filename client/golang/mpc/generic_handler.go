/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
)

func init() {
	var _ Invoker = new(GenericHandler)
	var _ macro.Caller = new(GenericHandler)
}

var executionKey = &prsim.Key{Name: "mesh.generic.execution", Dft: func() interface{} { return &GenericExecution{} }}

type GenericHandler struct {
	ReferenceHandler
}

func (that *GenericHandler) Invoke00(ctx context.Context, urn string, arguments map[string]interface{}) ([]interface{}, error) {
	that.refer = that
	mtx := ContextWith(ctx)
	uname := types.FromURN(ctx, urn)
	execution := new(GenericExecution)
	execution.Init(uname, arguments)
	mtx.SetAttribute(executionKey, execution)

	parameters := GenericParameters{}
	parameters.SetAttachments(mtx, map[string]string{})
	parameters.PutAll(arguments)

	invocation := &ServiceInvocation{
		Proxy:      that,
		Inspector:  execution.Inspect(),
		Parameters: parameters,
		Buffer:     nil,
		Execution:  execution,
		URN:        uname,
	}
	mtx.RewriteURN(that.rewriteURN(mtx, execution))
	mtx.SetAttribute(AddressKey, that.rewriteAddress(mtx, urn))
	ret, err := that.withFilters().Invoke(mtx, invocation)
	if nil != err {
		return nil, cause.Error(err)
	}
	x, ok := ret.(*GenericReturns)
	if ok {
		return x.GetContent(ctx), err
	}
	return nil, cause.Errorf("Cant resolve response ")
}

func (that *GenericHandler) ReferExecution(ctx context.Context, inspector macro.Inspector) (Execution, error) {
	if execution, ok := ContextWith(ctx).GetAttribute(executionKey).(Execution); ok {
		return execution, nil
	} else {
		return nil, cause.Errorf("Compatible error, execution not exist in context ")
	}
}
