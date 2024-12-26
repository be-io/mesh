/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/crypt"
	"github.com/opendatav/mesh/client/golang/macro"
	_ "github.com/opendatav/mesh/client/golang/proxy"
	"github.com/opendatav/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(boostRuntimeAware)

type boostRuntimeAware struct {
	KMS         prsim.KMS
	YAML        codec.Codec
	JSON        codec.Codec
	RSA2        prsim.Cryptor
	Singularity prsim.Singularity
	Dispatcher  prsim.Dispatcher
}

func (that *boostRuntimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.client.boost.runtime"}
}

func (that *boostRuntimeAware) Init() error {
	that.KMS = macro.Load(prsim.IKMS).Get(macro.MeshSys).(prsim.KMS)
	that.YAML = macro.Load(codec.ICodec).Get(codec.YAML).(codec.Codec)
	that.JSON = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.RSA2 = macro.Load(prsim.ICryptor).Get(crypt.RSA2).(prsim.Cryptor)
	that.Singularity = macro.Load(prsim.ISingularity).GetAny(macro.MeshSPI, macro.MeshNop).(prsim.Singularity)
	that.Dispatcher = macro.Load(prsim.IDispatcher).Get(macro.MeshSPI).(prsim.Dispatcher)
	return nil
}
