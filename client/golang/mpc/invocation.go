/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"bytes"
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ Invocation = new(ServiceInvocation)
}

type Invocation interface {

	// GetProxy the delegate target object.
	GetProxy() Invoker

	// GetInspector interface.
	GetInspector() macro.Inspector

	// GetParameters parameters.
	GetParameters() macro.Parameters

	// GetArguments parameters.
	GetArguments(ctx context.Context) []interface{}

	// GetAttachments the attachments. The attributes will be serialized.
	GetAttachments(ctx context.Context) map[string]string

	// GetBuffer get the input buffer.
	GetBuffer() *bytes.Buffer

	// GetExecution get the invocation execution.
	GetExecution() Execution

	// GetURN get the invoked urn.
	GetURN() *types.URN
}

type ServiceInvocation struct {
	Proxy      Invoker
	Inspector  macro.Inspector
	Parameters macro.Parameters
	Buffer     *bytes.Buffer
	Execution  Execution
	URN        *types.URN
}

func (that *ServiceInvocation) GetProxy() Invoker {
	return that.Proxy
}

func (that *ServiceInvocation) GetInspector() macro.Inspector {
	return that.Inspector
}

func (that *ServiceInvocation) GetParameters() macro.Parameters {
	return that.Parameters
}

func (that *ServiceInvocation) GetArguments(ctx context.Context) []interface{} {
	return that.Parameters.GetArguments(ctx)
}

func (that *ServiceInvocation) GetAttachments(ctx context.Context) map[string]string {
	return that.Parameters.GetAttachments(ctx)
}

func (that *ServiceInvocation) GetBuffer() *bytes.Buffer {
	return that.Buffer
}

func (that *ServiceInvocation) GetExecution() Execution {
	return that.Execution
}

func (that *ServiceInvocation) GetURN() *types.URN {
	return that.URN
}

type InvocationRuntime struct {
	invocation Invocation
}

func (that *InvocationRuntime) Parse(ctx context.Context, ptr interface{}) error {
	input := that.invocation.GetInspector().NewInbound()
	input.SetArguments(ctx, that.invocation.GetArguments(ctx)...)
	input.SetAttachments(ctx, that.invocation.GetAttachments(ctx))
	buff, err := aware.JSON.Encode(input)
	if nil != err {
		return cause.Error(err)
	}
	_, err = aware.JSON.Decode(buff, ptr)
	return cause.Error(err)
}

func (that *InvocationRuntime) ParseHeader(ctx context.Context, headers map[string]string) error {
	// not supported!
	return nil
}

func (that *InvocationRuntime) Dispatch(ctx context.Context, name string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
