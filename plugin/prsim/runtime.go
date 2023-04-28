/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	_ "github.com/be-io/mesh/client/golang/proxy"
	"github.com/be-io/mesh/client/golang/prsim"
	_ "github.com/be-io/mesh/client/golang/system"
	"github.com/be-io/mesh/plugin/cache"
	"github.com/be-io/mesh/plugin/kms"
	"github.com/be-io/mesh/plugin/metabase"
	nsqio "github.com/be-io/mesh/plugin/nsq"
	"github.com/be-io/mesh/plugin/redis"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	KMS              prsim.KMS
	Codec            codec.Codec
	RemoteNet        prsim.Network
	LocalNet         prsim.Network
	RemoteSubscriber prsim.Subscriber
	NSQPublisher     prsim.Publisher
	Sequence         prsim.Sequence
	Listener         prsim.Listener
	Registry         prsim.Registry
	Cache            prsim.Cache
	Local            prsim.Cache
	Redis            prsim.Cache
	KV               prsim.KV
	RemoteRegistry   prsim.Registry
	Dispatcher       prsim.Dispatcher
	Scheduler        prsim.Scheduler
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.prsim.runtime"}
}

func (that *runtimeAware) Init() error {
	that.KMS = macro.Load(prsim.IKMS).Get(kms.Local).(prsim.KMS)
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.RemoteNet = macro.Load(prsim.INetwork).Get(macro.MeshMPI).(prsim.Network)
	that.LocalNet = macro.Load(prsim.INetwork).Get(metabase.Name).(prsim.Network)
	that.RemoteSubscriber = macro.Load(prsim.ISubscriber).Get(macro.MeshMPI).(prsim.Subscriber)
	that.NSQPublisher = macro.Load(prsim.IPublisher).Get(nsqio.Name).(prsim.Publisher)
	that.Sequence = macro.Load(prsim.ISequence).Get(metabase.Name).(prsim.Sequence)
	that.Listener = macro.Load(prsim.IListener).Get(PRSListener).(prsim.Listener)
	that.Registry = macro.Load(prsim.IRegistry).Get(macro.MeshSPI).(prsim.Registry)
	that.RemoteRegistry = macro.Load(prsim.IRegistry).Get(macro.MeshMPI).(prsim.Registry)
	that.Cache = macro.Load(prsim.ICache).Get(macro.MeshSys).(prsim.Cache)
	that.Local = macro.Load(prsim.ICache).Get(cache.Unshared).(prsim.Cache)
	that.Redis = macro.Load(prsim.ICache).Get(redis.Name).(prsim.Cache)
	that.KV = macro.Load(prsim.IKV).Get(metabase.Name).(prsim.KV)
	that.Dispatcher = macro.Load(prsim.IDispatcher).Get(macro.MeshSPI).(prsim.Dispatcher)
	that.Scheduler = macro.Load(prsim.IScheduler).Get(macro.MeshSys).(prsim.Scheduler)
	return nil
}
