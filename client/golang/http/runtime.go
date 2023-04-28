/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package http

import (
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	_ "github.com/be-io/mesh/client/golang/proxy"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	Codec codec.Codec
	Eden  mpc.Eden
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.client.http.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Eden = macro.Load(mpc.IEden).Get(macro.MeshSPI).(mpc.Eden)
	return nil
}
