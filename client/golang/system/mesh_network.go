/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
)

func init() {
	var _ prsim.Network = network
	macro.Provide(prsim.INetwork, network)
}

var (
	network    = new(systemNetwork)
	environKey = &prsim.Key{Name: "mesh.mpc.environ", Dft: func() interface{} { return nil }}
	Environ    = new(macro.Once[*types.Environ]).With(func() *types.Environ {
		ctx := mpc.Context()
		env, err := network.GetEnviron(ctx)
		if nil != err {
			log.Error(ctx, err.Error())
		}
		return env
	})
)

type systemNetwork struct {
	environ *types.Environ
}

func (that *systemNetwork) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *systemNetwork) GetEnviron(ctx context.Context) (*types.Environ, error) {
	if nil != that.environ {
		return that.environ, nil
	}
	mtx := mpc.CloneContext(ctx)
	attr := mtx.GetAttribute(environKey)
	if env, ok := attr.(*types.Environ); ok {
		return env, nil
	}
	environ := &types.Environ{
		NodeId: types.LocalNodeId,
		InstId: types.LocalInstId,
	}
	mtx.SetAttribute(environKey, environ)
	env, err := aware.Network.GetEnviron(mtx)
	if nil != err {
		return nil, cause.Error(err)
	}
	that.environ = env
	return that.environ, nil
}

func (that *systemNetwork) Accessible(ctx context.Context, route *types.Route) (bool, error) {
	return aware.Network.Accessible(ctx, route)
}

func (that *systemNetwork) Refresh(ctx context.Context, routes []*types.Route) error {
	return aware.Network.Refresh(ctx, routes)
}

func (that *systemNetwork) GetRoute(ctx context.Context, nodeId string) (*types.Route, error) {
	return aware.Network.GetRoute(ctx, nodeId)
}

func (that *systemNetwork) GetRoutes(ctx context.Context) ([]*types.Route, error) {
	return aware.Network.GetRoutes(ctx)
}

func (that *systemNetwork) GetDomains(ctx context.Context, kind string) ([]*types.Domain, error) {
	return aware.Network.GetDomains(ctx, kind)
}

func (that *systemNetwork) PutDomains(ctx context.Context, kind string, domains []*types.Domain) error {
	return aware.Network.PutDomains(ctx, kind, domains)
}

func (that *systemNetwork) Weave(ctx context.Context, route *types.Route) error {
	return aware.Network.Weave(ctx, route)
}

func (that *systemNetwork) Ack(ctx context.Context, route *types.Route) error {
	return aware.Network.Ack(ctx, route)
}

func (that *systemNetwork) Disable(ctx context.Context, nodeId string) error {
	return aware.Network.Disable(ctx, nodeId)
}

func (that *systemNetwork) Enable(ctx context.Context, nodeId string) error {
	return aware.Network.Enable(ctx, nodeId)
}

func (that *systemNetwork) Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Route], error) {
	return aware.Network.Index(ctx, index)
}

func (that *systemNetwork) Version(ctx context.Context, nodeId string) (*types.Versions, error) {
	return aware.Network.Version(ctx, nodeId)
}

func (that *systemNetwork) Instx(ctx context.Context, index *types.Paging) (*types.Page[*types.Institution], error) {
	return aware.Network.Instx(ctx, index)
}

func (that *systemNetwork) Instr(ctx context.Context, institutions []*types.Institution) error {
	return aware.Network.Instr(ctx, institutions)
}

func (that *systemNetwork) Ally(ctx context.Context, nodeIds []string) error {
	return aware.Network.Ally(ctx, nodeIds)
}

func (that *systemNetwork) Disband(ctx context.Context, nodeIds []string) error {
	return aware.Network.Disband(ctx, nodeIds)
}

func (that *systemNetwork) Assert(ctx context.Context, feature string, nodeIds []string) (bool, error) {
	return aware.Network.Assert(ctx, feature, nodeIds)
}
