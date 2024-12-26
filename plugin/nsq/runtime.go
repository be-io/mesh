/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package nsqio

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
	Codec           codec.Codec
	Subscriber      prsim.Subscriber
	RemotePublisher prsim.Publisher
	Network         prsim.Network
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.nsq.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Subscriber = macro.Load(prsim.ISubscriber).Get(macro.MeshSPI).(prsim.Subscriber)
	that.RemotePublisher = macro.Load(prsim.IPublisher).Get(macro.MeshMPI).(prsim.Publisher)
	that.Network = macro.Load(prsim.INetwork).GetAny(macro.MeshSys).(prsim.Network)
	return nil
}
