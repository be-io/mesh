/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import (
	"context"
	"github.com/dgraph-io/badger/v4"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/system"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/timshannon/badgerhold/v4"
	"math"
	"time"
)

func init() {
	var _ prsim.Network = lan
	macro.Provide(prsim.INetwork, lan)
}

var lan = new(network)

// 原来的组网管理改成节点管理。原来的操作按钮组网改成手动组网。
// 里面新增一个状态位：是否已组网，和原有的是否联通区分开
// 新增一个标记列：节点来源（手动组网、自动发现、入网申请）
// 操作里面对组网状态进行判断，未组网的，新增一个组网申请操作，入网申请状态，新增一个申请审批
//
// 变更点：
// 1、组网管理升级成节点管理，组网只是节点管理中的一个操作。
// 2、新增组网申请与审批两个操作。
// 3、支持手动组网和自动发现组网，节点列表里面三种来源：手动组网、自动发现、入网申请
type network struct {
}

func (that *network) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *network) GetEnviron(ctx context.Context) (*types.Environ, error) {
	return system.Environ.Get(), nil
}

func (that *network) Accessible(ctx context.Context, route *types.Route) (bool, error) {
	return true, nil
}

func (that *network) Refresh(ctx context.Context, routes []*types.Route) error {
	return TX(ctx, func(ctx context.Context, tx *badger.Txn, ss *badgerhold.Store) error {
		for _, route := range routes {
			if "" == route.NodeId {
				continue
			}
			if (route.Status & int32(types.Removed)) == int32(types.Removed) {
				if err := ss.TxDelete(tx, route.NodeId, new(EdgeEnt)); nil != err {
					return cause.Error(err)
				}
				continue
			}
			extra, err := new(types.RouteAttribute).Encode(ctx, route, aware.Codec)
			if nil != err {
				log.Warn(ctx, "Encode:%s", err.Error())
				return cause.Error(err)
			}
			prev := new(EdgeEnt)
			if err = ss.TxGet(tx, route.NodeId, prev); nil != err && !IsNotFound(err) {
				return cause.Error(err)
			}
			certification := new(types.RouteCertificate)
			certification.Decode(ctx, prev.Certificate, aware.Codec)
			err = ss.TxUpsert(tx, route.NodeId, &EdgeEnt{
				NodeID:      route.NodeId,
				InstID:      tool.Anyone(route.InstId, prev.InstID),
				InstName:    tool.Anyone(route.Name, route.InstName, prev.InstName),
				Address:     tool.Anyone(route.URC().String(), prev.Address),
				Describe:    tool.Anyone(route.Describe, prev.Describe),
				Certificate: certification.Override0(route).Encode(ctx, aware.Codec),
				Status:      types.Weaved.OrElse(prev.Status, types.Connected.OrElse(prev.Status, route.Status)),
				Version:     prev.Version + 1,
				AuthCode:    tool.Anyone(route.AuthCode, prev.AuthCode),
				Extra:       tool.Anyone(extra, prev.Extra),
				ExpireAt:    time.UnixMilli(tool.Ternary(route.ExpireAt < 1, time.Now().AddDate(10, 0, 0).UnixMilli(), route.ExpireAt)),
				CreateAt:    tool.Anyone(prev.CreateAt, time.Now()),
				UpdateAt:    time.Now(),
				CreateBy:    tool.Anyone(prev.CreateBy, route.CreateBy),
				UpdateBy:    tool.Anyone(route.UpdateBy, prev.UpdateBy),
				Group:       route.Group,
			})
			if nil != err {
				return cause.Error(err)
			}
		}
		return nil
	})
}

func (that *network) GetRoute(ctx context.Context, nodeId string) (*types.Route, error) {
	return TR(ctx, func(ctx context.Context, ss *badgerhold.Store) (*types.Route, error) {
		edge := new(EdgeEnt)
		err := ss.Get(nodeId, edge)
		if nil == err {
			return that.ToRoute(ctx, edge), nil
		}
		if !IsNotFound(err) {
			return nil, cause.Error(err)
		}
		err = ss.FindOne(edge, badgerhold.Where("InstID").Eq(nodeId))
		if nil == err {
			return that.ToRoute(ctx, edge), nil
		}
		if !IsNotFound(err) {
			return nil, cause.Error(err)
		}
		return nil, nil
	})
}

func (that *network) GetRoutes(ctx context.Context) ([]*types.Route, error) {
	return TR(ctx, func(ctx context.Context, ss *badgerhold.Store) ([]*types.Route, error) {
		var edges []*EdgeEnt
		err := ss.Find(&edges, new(badgerhold.Query).Limit(math.MaxInt))
		if nil == err {
			return that.ToRoutes(ctx, edges), nil
		}
		if !IsNotFound(err) {
			return nil, cause.Error(err)
		}
		return nil, nil
	})
}

func (that *network) GetDomains(ctx context.Context, kind string) ([]*types.Domain, error) {
	return nil, nil
}

func (that *network) PutDomains(ctx context.Context, kind string, domains []*types.Domain) error {
	return nil
}

func (that *network) Weave(ctx context.Context, route *types.Route) error {
	route.Status = route.Status | int32(types.FromWeaving)
	return cause.Error(that.Refresh(ctx, []*types.Route{route}))
}

func (that *network) Ack(ctx context.Context, route *types.Route) error {
	route.Status = route.Status | int32(types.Weaved)
	return cause.Error(that.Refresh(ctx, []*types.Route{route}))
}

func (that *network) Disable(ctx context.Context, nodeId string) error {
	edge, err := that.GetRoute(ctx, nodeId)
	if nil != err {
		return cause.Error(err)
	}
	if nil == edge {
		return cause.Error(cause.Errorable(cause.NetNotWeave))
	}
	edge.Status = edge.Status | int32(types.Disabled)
	return cause.Error(that.Refresh(ctx, []*types.Route{edge}))
}

func (that *network) Enable(ctx context.Context, nodeId string) error {
	edge, err := that.GetRoute(ctx, nodeId)
	if nil != err {
		return cause.Error(err)
	}
	if nil == edge {
		return cause.Error(cause.Errorable(cause.NetNotWeave))
	}
	edge.Status = edge.Status & int32(^types.Disabled)
	return cause.Error(that.Refresh(ctx, []*types.Route{edge}))
}

func (that *network) Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Route], error) {
	return TR(ctx, func(ctx context.Context, ss *badgerhold.Store) (*types.Page[*types.Route], error) {
		con := badgerhold.Where("Status").MatchFunc(func(ra *badgerhold.RecordAccess) (bool, error) {
			// (`status` & 320) = 0
			return ra.Record() == nil, nil
		}).SortBy("NodeId").Skip(int(index.Index * index.Limit)).Limit(int(index.Limit))
		nodeId := index.Factor["node_id"]
		instName := index.Factor["inst_name"]
		if nil != nodeId && "" != nodeId {
			con.And("NodeId").Contains(nodeId)
		}
		if nil != instName && "" != instName {
			con.And("InstName").Contains(instName)
		}
		total, err := ss.Count(new(EdgeEnt), con)
		if nil != err {
			return nil, cause.Error(err)
		}
		sid := index.SID
		if "" == sid {
			sid = tool.NextID()
		}
		var edges []*EdgeEnt
		if err = ss.Find(&edges, con); nil != err {
			return nil, cause.Error(err)
		}
		return types.Reset(index, int64(total), that.ToRoutes(ctx, edges)), nil
	})
}

func (that *network) Version(ctx context.Context, nodeId string) (*types.Versions, error) {
	return &types.Versions{Version: prsim.Version}, nil
}

func (that *network) Instx(ctx context.Context, index *types.Paging) (*types.Page[*types.Institution], error) {
	return nil, nil
}

func (that *network) Instr(ctx context.Context, institutions []*types.Institution) error {
	return nil
}

func (that *network) Ally(ctx context.Context, nodeIds []string) error {
	return nil
}

func (that *network) Disband(ctx context.Context, nodeIds []string) error {
	return nil
}

func (that *network) Assert(ctx context.Context, feature string, nodeIds []string) (bool, error) {
	return false, nil
}

func (that *network) ToRoutes(ctx context.Context, edges []*EdgeEnt) []*types.Route {
	var es []*types.Route
	for _, edge := range edges {
		es = append(es, that.ToRoute(ctx, edge))
	}
	return es
}

func (that *network) ToRoute(ctx context.Context, route *EdgeEnt) *types.Route {
	certification := new(types.RouteCertificate)
	certification.Decode(ctx, route.Certificate, aware.Codec)
	attribute := new(types.RouteAttribute).Decode(ctx, route.Extra, aware.Codec)
	return &types.Route{
		NodeId:      route.NodeID,
		InstId:      route.InstID,
		Name:        route.InstName,
		InstName:    route.InstName,
		Address:     route.Address,
		Describe:    route.Describe,
		HostRoot:    certification.HostRoot,
		HostKey:     certification.HostKey,
		HostCrt:     certification.HostCrt,
		GuestRoot:   certification.GuestRoot,
		GuestKey:    certification.GuestKey,
		GuestCrt:    certification.GuestCrt,
		Status:      route.Status,
		Version:     route.Version,
		AuthCode:    route.AuthCode,
		ExpireAt:    route.ExpireAt.UnixMilli(),
		Extra:       attribute.Compat(ctx, aware.Codec),
		CreateAt:    route.CreateAt.Format(log.DateFormat),
		UpdateAt:    route.UpdateAt.Format(log.DateFormat),
		CreateBy:    route.CreateBy,
		UpdateBy:    route.UpdateBy,
		Group:       route.Group,
		Upstream:    attribute.Upstream,
		Downstream:  attribute.Downstream,
		StaticIP:    attribute.StaticIP,
		Proxy:       attribute.Proxy,
		Concurrency: attribute.Concurrency,
	}
}
