/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"bytes"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"google.golang.org/grpc"
	"runtime/debug"
)

func init() {
	var _ grpc.StreamHandler = new(bindableService).OnNext
	var _ StreamService = new(bindableService)
}

var Service = new(bindableService)

type StreamService interface {
	OnNext(srv interface{}, stream grpc.ServerStream) error
}

type bindableService struct {
}

func (that *bindableService) Start() error {
	return nil
}

func (that *bindableService) OnNext(srv interface{}, stream grpc.ServerStream) (ex error) {
	mtx := mpc.ContextWith(Interceptors.ServerContext(stream.Context()))
	defer func() {
		if err := recover(); nil != err {
			log.Error(mtx, "%v", err)
			log.Error(mtx, string(debug.Stack()))
			ex = that.OnError(mtx, aware.JSON, stream, cause.Errorf("%v", err))
		}
	}()
	cdc, err := that.OnService(mtx, stream)
	if nil == err {
		return
	}
	log.Warn(mtx, "%s, %s", mtx.GetUrn(), err.Error())
	return that.OnError(mtx, cdc, stream, err)
}

func (that *bindableService) OnError(ctx prsim.Context, cdc codec.Codec, stream grpc.ServerStream, except error) error {
	code, message := cause.Parse(except)
	outbound := &types.Outbound{
		Code:    code,
		Message: message,
	}
	buff, err := cdc.Encode(outbound)
	if nil != err {
		log.Error(ctx, "%s serialize outbound body, %s", ctx.GetUrn(), err.Error())
		return cause.Error(err)
	}
	if err = stream.SendMsg(buff); nil != err {
		log.Error(ctx, "%s send outbound body, %s", ctx.GetUrn(), err.Error())
		return cause.Error(err)
	}
	return nil
}

func (that *bindableService) OnFallback(srv interface{}, stream grpc.ServerStream) (ex error) {
	ctx := mpc.Context()
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, "%v", err)
			log.Error(ctx, string(debug.Stack()))
			ex = that.OnError(ctx, aware.JSON, stream, cause.Errorf("%v", err))
		}
	}()
	getCodec := func() codec.Codec {
		if "" == ctx.GetUrn() {
			return aware.JSON
		}
		log.Info(ctx, "Fallback %s", ctx.GetUrn())
		urn := types.FromURN(ctx, ctx.GetUrn())
		name := mpc.MeshFlag.OfCodec(urn.Flag.Codec).Name()
		cdc, ok := macro.Load(codec.ICodec).Get(name).(codec.Codec)
		if !ok {
			log.Error(ctx, "No Codec provider named %s exist ", name)
			return aware.JSON
		}
		return cdc
	}
	outbound := &types.Outbound{Code: cause.NotFound.Code, Message: fmt.Sprintf("%s %s", cause.NotFound.Message, ctx.GetUrn())}
	buff, err := getCodec().Encode(outbound)
	if nil != err {
		log.Error(ctx, "%s serialize outbound body fail, %s", ctx.GetUrn(), err.Error())
		return cause.Error(err)
	}
	if err = stream.SendMsg(buff); nil != err {
		log.Error(ctx, "%s send outbound body fail, %s", ctx.GetUrn(), err.Error())
		return cause.Error(err)
	}
	return nil
}

func (that *bindableService) OnService(ctx prsim.Context, stream grpc.ServerStream) (codec.Codec, error) {
	if "" == ctx.GetUrn() {
		return aware.JSON, cause.Errorf("Cant resolve grpc authority.")
	}
	log.Debug(ctx, "Receive %s", ctx.GetUrn())
	urn := types.FromURN(ctx, ctx.GetUrn())
	name := mpc.MeshFlag.OfCodec(urn.Flag.Codec).Name()
	cdc, ok := macro.Load(codec.ICodec).Get(name).(codec.Codec)
	if !ok {
		return aware.JSON, cause.Errorf("No Codec provider named %s exist ", name)
	}
	var inbound bytes.Buffer
	if err := stream.RecvMsg(&inbound); nil != err {
		return cdc, cause.Error(err)
	}
	log.Debug(ctx, "input=%s", inbound.String())
	execution, err := aware.EDEN.Infer(ctx, ctx.GetUrn())
	if nil != err {
		return cdc, cause.Error(err)
	}
	if nil == execution {
		return cdc, cause.Errorcf(cause.NotFound, "No mpi found for %s.", ctx.GetUrn())
	}
	parameters := execution.Inspect().GetIntype()
	if nil != err {
		return cdc, cause.Error(err)
	}
	if _, err = cdc.Decode(&inbound, parameters); nil != err {
		return cdc, cause.Error(err)
	}
	invocation := &mpc.ServiceInvocation{
		Proxy:      execution,
		Inspector:  execution.Inspect(),
		Parameters: parameters,
		Buffer:     &inbound,
		Execution:  execution,
		URN:        urn,
	}
	output, err := execution.Invoke(ctx, invocation)
	if nil != err {
		return cdc, cause.Error(err)
	}
	returns := execution.Inspect().NewOutbound()
	returns.SetCode(cause.Success.Code)
	returns.SetMessage(cause.Success.Message)
	returns.SetContent(ctx, output)
	outbound, err := cdc.Encode(returns)
	if nil != err {
		return cdc, cause.Error(err)
	}
	return cdc, cause.Error(stream.SendMsg(outbound.Bytes()))
}
