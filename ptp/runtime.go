/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package ptp

import (
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	Codec    codec.Codec
	Cache    prsim.Cache
	Session  prsim.Session
	KV       prsim.KV
	Registry prsim.Registry
	KMS      prsim.KMS
	Network  prsim.Network
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.ptp.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Cache = macro.Load(prsim.ICache).Get(macro.MeshSPI).(prsim.Cache)
	that.Session = macro.Load(prsim.ISession).Get(macro.MeshSPI).(prsim.Session)
	that.KV = macro.Load(prsim.IKV).Get(macro.MeshSPI).(prsim.KV)
	that.Registry = macro.Load(prsim.IRegistry).Get(macro.MeshSPI).(prsim.Registry)
	that.KMS = macro.Load(prsim.IKMS).Get(macro.MeshSPI).(prsim.KMS)
	that.Network = macro.Load(prsim.INetwork).Get(macro.MeshSPI).(prsim.Network)
	return nil
}
