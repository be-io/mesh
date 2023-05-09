/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package raft

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
	Codec   codec.Codec
	Cluster channel
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.raft.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Cluster = macro.Load(iCluster).Get(macro.MeshMPI).(channel)
	return nil
}
