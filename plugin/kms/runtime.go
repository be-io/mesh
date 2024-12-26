/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package kms

import (
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/crypt"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	Cryptor prsim.Cryptor
	YAML    codec.Codec
	JSON    codec.Codec
	KV      prsim.KV
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.kms.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Cryptor = macro.Load(prsim.ICryptor).Get(crypt.RSA2).(prsim.Cryptor)
	that.YAML = macro.Load(codec.ICodec).Get(codec.YAML).(codec.Codec)
	that.JSON = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.KV = macro.Load(prsim.IKV).Get(macro.MeshSys).(prsim.KV)
	return nil
}
