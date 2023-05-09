/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package eval

import (
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	_ "github.com/be-io/mesh/client/golang/proxy"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	KV       prsim.KV
	JSON     codec.Codec
	Sequence prsim.Sequence
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.eval.runtime"}
}

func (that *runtimeAware) Init() error {
	that.KV = macro.Load(prsim.IKV).Get(macro.MeshSys).(prsim.KV)
	that.Sequence = macro.Load(prsim.ISequence).Get(macro.MeshSys).(prsim.Sequence)
	that.JSON = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	return nil
}
