/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

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
	Codec     codec.Codec
	Cache     prsim.Cache
	Scheduler prsim.Scheduler
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.redis.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Cache = macro.Load(prsim.ICache).Get(macro.MeshSPI).(prsim.Cache)
	that.Scheduler = macro.Load(prsim.IScheduler).Get(macro.MeshSys).(prsim.Scheduler)
	return nil
}
