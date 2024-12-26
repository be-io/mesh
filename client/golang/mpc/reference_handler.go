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
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"strconv"
	"strings"
)

func init() {
	var _ Invoker = new(ReferenceHandler)
	var _ macro.Caller = new(ReferenceHandler)
}

type Referencable interface {
	ReferExecution(ctx context.Context, inspector macro.Inspector) (Execution, error)
}

type ReferenceHandler struct {
	MPI      macro.MPI
	invoker  Invoker
	consumer Consumer
	refer    Referencable
}

func (that *ReferenceHandler) withFilters() Invoker {
	if nil == that.invoker {
		that.invoker = composite(that, CONSUMER)
	}
	return that.invoker
}

func (that *ReferenceHandler) withRefer() Referencable {
	if nil == that.refer {
		that.refer = that
	}
	return that.refer
}

func (that *ReferenceHandler) ReferExecution(ctx context.Context, inspector macro.Inspector) (Execution, error) {
	eden, ok := macro.Load(IEden).Get(macro.MeshSPI).(Eden)
	if !ok {
		return nil, cause.Errorf("No Eden provider named %s exist ", macro.MeshSPI)
	}
	execution, err := eden.Refer(ctx, that.MPI, inspector.GetDeclaredKind(), inspector)
	if nil != err {
		return nil, cause.Error(err)
	}
	if nil != execution {
		return execution, nil
	}
	return nil, cause.Errorf("Method %s cant be compatible", inspector.GetName())
}

func (that *ReferenceHandler) rewriteURN(ctx context.Context, execution Execution) string {
	mtx := ContextWith(ctx)
	uname, unameOk := mtx.GetAttribute(RemoteUname).(string)
	name, nameOk := mtx.GetAttribute(RemoteName).(string)
	principal := mtx.GetPrincipals().Peek()
	nop := nil == principal || ("" == principal.NodeId && "" == principal.InstId)
	if nop && (!unameOk || "" == uname) && (!nameOk || "" == name) {
		return execution.Schema().GetURN()
	}
	urn := types.FromURN(ctx, execution.Schema().GetURN())
	if nil != principal && "" != principal.NodeId {
		urn.NodeId = principal.NodeId
	}
	if nil != principal && "" != principal.InstId {
		urn.NodeId = principal.InstId
	}
	if unameOk && "" != uname {
		urn.Name = uname
	}
	if nameOk && "" != name {
		urn.Name = strings.ReplaceAll(urn.Name, "${mesh.name}", name)
	}
	return urn.String()
}

func (that *ReferenceHandler) rewriteAddress(ctx prsim.Context, uns string) string {
	if addr, ok := ctx.GetAttribute(AddressKey).(string); ok && "" != addr {
		return addr
	}
	urn := types.FromURN(ctx, uns)
	if strings.Index(urn.Name, "mesh.") == 0 {
		return tool.Address.Get().Any()
	}
	if "" != tool.Direct.Get() {
		names := strings.Split(tool.Direct.Get(), ",")
		for _, name := range names {
			pair := strings.Split(name, "=")
			if that.isDirect(urn, pair) {
				return pair[1]
			}
		}
	}
	address := strings.ReplaceAll(urn.Flag.Address, ".", "")
	if v, err := strconv.ParseFloat(address, 64); nil == err && v > 0 {
		return fmt.Sprintf("%s:%s", urn.Flag.Address, urn.Flag.Port)
	}
	return tool.Address.Get().Any()
}

func (that *ReferenceHandler) isDirect(urn *types.URN, pair []string) bool {
	if len(pair) < 2 || "" == pair[1] {
		return false
	}
	if !strings.Contains(pair[0], "@") {
		return strings.Index(urn.Name, pair[0]) == 0
	}
	nn := strings.Split(pair[0], "@")
	if len(nn) < 2 || "" == nn[1] {
		return false
	}
	return strings.ToLower(urn.NodeId) == strings.ToLower(nn[1]) && strings.Index(urn.Name, nn[0]) == 0
}

func (that *ReferenceHandler) withConsumer() (Consumer, error) {
	if nil == that.consumer {
		pv := macro.Load(IConsumer)
		co, ok := pv.Default().(Consumer)
		if !ok {
			return nil, cause.Errorf("No Consumer provider named %s exist ", pv.Name())
		}
		that.consumer = co
	}
	return that.consumer, nil
}

func (that *ReferenceHandler) Call(ctx context.Context, proxy interface{}, method macro.Inspector, args ...interface{}) (interface{}, error) {
	mtx := ContextWith(ctx)
	execution, err := that.withRefer().ReferExecution(mtx, method)
	if nil != err {
		return nil, cause.Error(err)
	}
	urn := that.rewriteURN(mtx, execution)
	mtx.RewriteURN(urn)
	mtx.SetAttribute(AddressKey, that.rewriteAddress(mtx, urn))

	parameters := method.NewInbound()
	parameters.SetArguments(mtx, args...)
	parameters.SetAttachments(mtx, map[string]string{})

	invocation := &ServiceInvocation{
		Proxy:      that,
		Inspector:  method,
		Parameters: parameters,
		Buffer:     nil,
		Execution:  execution,
		URN:        types.FromURN(ctx, urn),
	}

	return that.withFilters().Invoke(mtx, invocation)
}

func (that *ReferenceHandler) Invoke(ctx context.Context, invocation Invocation) (interface{}, error) {
	mtx := ContextWith(ctx)
	execution, err := that.withRefer().ReferExecution(mtx, invocation.GetInspector())
	if nil != err {
		return nil, cause.Error(err)
	}
	consume, err := that.withConsumer()
	if nil != err {
		return nil, cause.Error(err)
	}
	cdc, ok := macro.Load(codec.ICodec).Get(execution.Schema().GetCodec()).(codec.Codec)
	if !ok {
		return nil, cause.Errorf("No Codec provider named %s exist ", execution.Schema().GetCodec())
	}
	input, err := cdc.Encode(invocation.GetParameters())
	if nil != err {
		return nil, cause.Error(err)
	}
	output, err := consume.Consume(mtx, mtx.GetUrn(), execution, input)
	if nil != err {
		return nil, cause.Error(err)
	}
	returns := invocation.GetInspector().GetRetype()
	if _, err = cdc.Decode(output, &returns); nil != err {
		return nil, cause.Error(err)
	}
	if nil != returns.GetCause(ctx) {
		return nil, cause.Errorcf(returns, returns.GetCause(ctx).Text)
	}
	if cause.Success.Code != returns.GetCode() {
		return nil, cause.Errorcf(returns, returns.GetMessage())
	}
	return returns, cause.Error(err)
}
