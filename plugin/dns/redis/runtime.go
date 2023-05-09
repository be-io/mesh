/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	Cache   prsim.Cache
	Network prsim.Network
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.mdns.redis.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Cache = macro.Load(prsim.ICache).Get(macro.MeshSys).(prsim.Cache)
	that.Network = macro.Load(prsim.INetwork).Get(macro.MeshSys).(prsim.Network)
	return nil
}
