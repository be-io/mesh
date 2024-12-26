/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/macro"
	_ "github.com/opendatav/mesh/client/golang/proxy"
	"github.com/opendatav/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	Codec         codec.Codec
	Network       prsim.Network
	Sequence      prsim.Sequence
	Registry      prsim.Registry
	Licenser      prsim.Licenser
	Cache         prsim.Cache
	DataHouse     prsim.DataHouse
	Graphics      prsim.Graphics
	Tokenizer     prsim.Tokenizer
	Cryptor       prsim.Cryptor
	KV            prsim.KV
	Scheduler     prsim.Scheduler
	Commercialize prsim.Commercialize
	Devops        prsim.Devops
	Endpoint      prsim.Endpoint
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.cmd.exe"}
}

func (that *runtimeAware) Init() error {
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Network = macro.Load(prsim.INetwork).Get(macro.MeshMPI).(prsim.Network)
	that.Sequence = macro.Load(prsim.ISequence).Get(macro.MeshMPI).(prsim.Sequence)
	that.Registry = macro.Load(prsim.IRegistry).Get(macro.MeshMPI).(prsim.Registry)
	that.Licenser = macro.Load(prsim.ILicenser).Get(macro.MeshMPI).(prsim.Licenser)
	that.Cache = macro.Load(prsim.ICache).Get(macro.MeshMPI).(prsim.Cache)
	that.DataHouse = macro.Load(prsim.IDataHouse).Get(macro.MeshMPI).(prsim.DataHouse)
	that.Graphics = macro.Load(prsim.IGraphics).Get(macro.MeshMPI).(prsim.Graphics)
	that.Tokenizer = macro.Load(prsim.ITokenizer).Get(macro.MeshMPI).(prsim.Tokenizer)
	that.Cryptor = macro.Load(prsim.ICryptor).Get(macro.MeshMPI).(prsim.Cryptor)
	that.KV = macro.Load(prsim.IKV).Get(macro.MeshMPI).(prsim.KV)
	that.Scheduler = macro.Load(prsim.IScheduler).Get(macro.MeshMPI).(prsim.Scheduler)
	that.Commercialize = macro.Load(prsim.ICommercialize).Get(macro.MeshMPI).(prsim.Commercialize)
	that.Devops = macro.Load(prsim.IDevops).Get(macro.MeshSys).(prsim.Devops)
	that.Endpoint = macro.Load(prsim.IEndpoint).Get(macro.MeshMPI).(prsim.Endpoint)
	return nil
}
