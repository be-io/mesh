/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package ptp

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/grpc"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
	macro.Provide(prsim.IListener, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	Codec    codec.Codec
	Cache    prsim.Cache
	Session  prsim.Session
	KV       prsim.KV
	Registry prsim.Registry
	KMS      prsim.KMS
	Network  prsim.Network
	Channel  grpc.Channel
	Routes   map[string]*types.Route
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.ptp.runtime"}
}

func (that *runtimeAware) Init() error {
	that.Codec = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Cache = macro.Load(prsim.ICache).Get(macro.MeshSPI).(prsim.Cache)
	that.Session = macro.Load(prsim.ISession).Get(macro.MeshSPI).(prsim.Session)
	that.KV = macro.Load(prsim.IKV).Get(macro.MeshSPI).(prsim.KV)
	that.Registry = macro.Load(prsim.IRegistry).Get(macro.MeshSPI).(prsim.Registry)
	that.KMS = macro.Load(prsim.IKMS).Get(macro.MeshSPI).(prsim.KMS)
	that.Network = macro.Load(prsim.INetwork).Get(macro.MeshSPI).(prsim.Network)
	that.Channel = macro.Load(grpc.IChannel).Get(grpc.Name).(grpc.Channel)
	return nil
}

func (that *runtimeAware) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.NetworkRouteRefresh}
}

func (that *runtimeAware) Listen(ctx context.Context, event *types.Event) error {
	var routes []*types.Route
	if err := event.TryGetObject(&routes); nil != err {
		return cause.Error(err)
	}
	rs := map[string]*types.Route{}
	for _, route := range routes {
		rs[route.NodeId] = route
		rs[route.InstId] = route
	}
	that.Routes = rs
	return nil
}

func (that *runtimeAware) WeavedCheck(ctx prsim.Context) error {
	if nil == that.Routes {
		return cause.NetNotWeave.Error()
	}
	nodeId := prsim.MeshTargetNodeId.Get(ctx.GetAttachments())
	if nil != that.Routes[nodeId] {
		return nil
	}
	instId := prsim.MeshTargetInstId.Get(ctx.GetAttachments())
	if nil != that.Routes[instId] {
		return nil
	}
	return cause.NetNotWeave.Error()
}
