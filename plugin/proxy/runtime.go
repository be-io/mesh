/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/grpc"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	RemoteNet      prsim.Network
	Network        prsim.Network
	Cache          prsim.Cache
	Codec          codec.Codec
	Dispatcher     prsim.Dispatcher
	Scheduler      prsim.Scheduler
	RemoteRegistry prsim.Registry
	LocalRegistry  prsim.Registry
	KV             prsim.KV
	Builtin        prsim.Builtin
	RemoteBuiltin  prsim.Builtin
	Channel        grpc.Channel
	Endpoint       prsim.Endpoint
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.runtime"}
}

func (that *runtimeAware) Init() error {
	that.RemoteNet = macro.Load(prsim.INetwork).Get(macro.MeshMPI).(prsim.Network)
	that.Cache = macro.Load(prsim.ICache).Get(macro.MeshSPI).(prsim.Cache)
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Network = macro.Load(prsim.INetwork).Get(macro.MeshSys).(prsim.Network)
	that.Dispatcher = macro.Load(prsim.IDispatcher).Get(macro.MeshSPI).(prsim.Dispatcher)
	that.Scheduler = macro.Load(prsim.IScheduler).Get(macro.MeshSys).(prsim.Scheduler)
	that.RemoteRegistry = macro.Load(prsim.IRegistry).Get(macro.MeshMPI).(prsim.Registry)
	that.LocalRegistry = macro.Load(prsim.IRegistry).Get(macro.MeshSys).(prsim.Registry)
	that.KV = macro.Load(prsim.IKV).Get(tool.Ternary(plugin.PROXY.Match(), macro.MeshMPI, macro.MeshSys)).(prsim.KV)
	that.Builtin = macro.Load(prsim.IBuiltin).Get(macro.MeshSys).(prsim.Builtin)
	that.RemoteBuiltin = macro.Load(prsim.IBuiltin).Get(macro.MeshMPI).(prsim.Builtin)
	that.Channel = macro.Load(grpc.IChannel).Get(grpc.Name).(grpc.Channel)
	that.Endpoint = macro.Load(prsim.IEndpoint).Get(macro.MeshMPI).(prsim.Endpoint)
	return nil
}
