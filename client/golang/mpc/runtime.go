/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	Network   prsim.Network
	Registry  prsim.Registry
	Scheduler prsim.Scheduler
	JSON      codec.Codec
	Builtin   prsim.Builtin
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.client.mpc.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Network = macro.Load(prsim.INetwork).GetAny(macro.MeshSys).(prsim.Network)
	that.Registry = macro.Load(prsim.IRegistry).Get(macro.MeshSys).(prsim.Registry)
	that.Scheduler = macro.Load(prsim.IScheduler).Get(macro.MeshSys).(prsim.Scheduler)
	that.JSON = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Builtin = macro.Load(prsim.IBuiltin).Get(macro.MeshSys).(prsim.Builtin)
	return nil
}
