/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/system"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	_ "github.com/be-io/mesh/plugin/cache"
	_ "github.com/be-io/mesh/plugin/kms"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	NetworkDomainKey     = "mesh.network.graph.domain"
	NetworkFeaturePrefix = "mesh.network.feature"
)

var _ prsim.Network = new(PRSINetwork)

// PRSINetwork
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSINetwork struct {
}

func (that *PRSINetwork) GetEnviron(ctx context.Context) (*types.Environ, error) {
	return aware.KMS.Environ(ctx)
}

func (that *PRSINetwork) Accessible(ctx context.Context, route *types.Route) (bool, error) {
	if nil == route {
		return true, nil
	}
	env, err := that.GetEnviron(ctx)
	if nil != err {
		return false, cause.Error(err)
	}
	nodeId := tool.Anyone(route.NodeId, route.InstId)
	if tool.IsLocalEnv(env, nodeId) {
		return true, nil
	}
	if "" == route.URC().String() {
		route, err = that.GetRoute(ctx, nodeId)
		if nil != err {
			return false, cause.Error(err)
		}
		if nil == route {
			return false, cause.Errorable(cause.NetNotWeave)
		}
	}
	//if !tool.Addressable(ctx, route.Address) {
	//	return false, cause.Errorable(cause.AddressError)
	//}
	if err = that.FlushProxy(ctx, route); nil != err {
		return false, cause.Error(err)
	}
	mtx := mpc.ContextWith(ctx).Resume(ctx)
	mtx.SetAttribute(mpc.TimeoutKey, time.Second*6)
	mtx.GetPrincipals().Push(&types.Principal{NodeId: route.NodeId, InstId: route.InstId})
	defer func() { mtx.GetPrincipals().Pop() }()
	testRoute := route.Copy()
	testRoute.Address = ""
	return aware.RemoteNet.Accessible(mtx, testRoute)
}

func (that *PRSINetwork) Refresh(ctx context.Context, routes []*types.Route) error {
	environ, err := that.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	myEnvFn := func(route *types.Route) (*types.Environ, error) {
		if tool.IsLocalEnv(environ, route.NodeId, route.InstId) {
			return environ, nil
		}
		//if !tool.Addressable(ctx, route.Address) {
		//	return nil, cause.Errorable(cause.AddressError)
		//}
		if err = that.FlushProxy(ctx, route); nil != err {
			return nil, cause.Error(err)
		}
		mtx := mpc.ContextWith(ctx).Resume(ctx)
		mtx.SetAttribute(mpc.TimeoutKey, time.Second*3)
		mtx.GetPrincipals().Push(&types.Principal{NodeId: route.NodeId, InstId: route.InstId})
		defer func() { mtx.GetPrincipals().Pop() }()
		return aware.RemoteNet.GetEnviron(mtx)
	}
	text, _ := aware.Codec.EncodeString(routes)
	log.Info(ctx, "Network refresh event=%s", text)
	for _, route := range routes {
		if "" != route.HostCrt {
			b, _ := pem.Decode([]byte(route.HostCrt))
			if _, err = x509.ParseCertificate(b.Bytes); nil != err {
				return cause.Error(err)
			}
		}
		if "" != route.GuestCrt {
			b, _ := pem.Decode([]byte(route.GuestCrt))
			if _, err = x509.ParseCertificate(b.Bytes); nil != err {
				return cause.Error(err)
			}
		}
		isLocal := tool.IsLocalEnv(environ, route.NodeId, route.InstId)
		if myEnv, err := myEnvFn(route); nil == err {
			route.Status = tool.Ternary(isLocal, types.Virtual.Not(route.Status|int32(types.Connected)|int32(types.Weaved)), route.Status|int32(types.Connected))
			route.NodeId = myEnv.NodeId
			route.InstId = myEnv.InstId
			route.Name = myEnv.InstName
			route.InstName = myEnv.InstName
		} else {
			log.Warn(ctx, "Determine route %s env with error, %s", route.NodeId, err.Error())
			route.Status = tool.Ternary(isLocal, types.Virtual.Not(route.Status|int32(types.Connected)|int32(types.Weaved)), route.Status&int32(^types.Connected))
		}
		if "" == route.NodeId || "" == route.InstId || "" == route.URC().String() {
			log.Warn(ctx, "%s|%s|%s|%s", route.NodeId, route.InstId, route.URC().String(), route.InstName)
			return cause.ValidateErrorf("Refresh network route, nodeId, instId, instName, address are required.")
		}
	}
	if err = aware.LocalNet.Refresh(ctx, routes); nil != err {
		return cause.Error(err)
	}
	return tabledataCaster.RouteManRefresh(ctx)
}

func (that *PRSINetwork) GetRoute(ctx context.Context, nodeId string) (*types.Route, error) {
	route, err := aware.LocalNet.GetRoute(ctx, nodeId)
	if nil != err {
		return nil, cause.Error(err)
	}
	if nil != route {
		return route, nil
	}
	if types.NodeSEQ(nodeId) == types.NodeSEQ(system.Environ.Get().NodeId) {
		instId := tool.Ternary(len(nodeId) == len(system.Environ.Get().InstId), nodeId, system.Environ.Get().InstId)
		return types.LocRoute(system.Environ.Get(), instId), nil
	}
	routes, err := that.GetRoutes(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	for _, r := range routes {
		if types.NodeSEQ(nodeId) == types.NodeSEQ(r.NodeId) {
			cr := r.Copy()
			cr.InstId = nodeId
			return cr, nil
		}
	}
	return nil, nil
}

func (that *PRSINetwork) GetRoutes(ctx context.Context) ([]*types.Route, error) {
	return aware.LocalNet.GetRoutes(ctx)
}

func (that *PRSINetwork) GetDomains(ctx context.Context, kind string) ([]*types.Domain, error) {
	var domains []*types.Domain
	if err := system.GetWithCache(ctx, aware.Cache, fmt.Sprintf("%s.%s", NetworkDomainKey, kind), &domains); nil != err {
		log.Warn(ctx, err.Error())
		return domains, cause.Error(err)
	}
	if len(domains) > 0 {
		return domains, nil
	}
	environ := os.Getenv("MESH_NET_DOMAINS")
	if "" == environ || !codec.Jsonizer.Valid([]byte(environ)) {
		return domains, nil
	}
	if err := codec.Jsonizer.Unmarshal([]byte(environ), &domains); nil != err {
		log.Warn(ctx, "Cant resolve domains environs, %s", err.Error())
	}
	return domains, nil
}

func (that *PRSINetwork) PutDomains(ctx context.Context, kind string, domains []*types.Domain) error {
	return system.PutWithCache(ctx, aware.Cache, fmt.Sprintf("%s.%s", NetworkDomainKey, kind), domains, time.Hour*24*365)
}

func (that *PRSINetwork) Weave(ctx context.Context, route *types.Route) error {
	environ, err := that.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	if "" == route.NodeId {
		return cause.Errorf("Network weave must bind nodeId.")
	}
	edge, err := func() (*types.Route, error) {
		edge, err := aware.LocalNet.GetRoute(ctx, route.NodeId)
		if nil != err {
			return nil, cause.Error(err)
		}
		if nil != edge {
			return edge, nil
		}
		return route, nil
	}()
	if nil != err {
		return cause.Error(err)
	}
	issue, err := aware.KMS.ApplyIssue(ctx, &types.KeyCsr{
		CNO:      environ.NodeId,
		PNO:      environ.NodeId,
		Domain:   fmt.Sprintf("%s.%s", environ.NodeId, types.MeshDomain),
		Subject:  environ.InstName,
		Length:   2048,
		ExpireAt: types.Time(time.Now().AddDate(10, 0, 0)),
		Mail:     fmt.Sprintf("mail@%s.%s", environ.NodeId, types.MeshDomain),
		IsCA:     false,
		CaCert:   environ.RootCrt,
		CaKey:    environ.RootKey,
	})
	if nil != err {
		return cause.Error(err)
	}
	sets := types.KeysSet(issue)
	if mpc.ContextWith(ctx).GetAttachments()["mesh.net.weave"] != "guest" {
		certification := &types.RouteCertificate{
			HostCrt: sets.Get(types.IssueCrtKey),
			HostKey: sets.Get(types.IssuePrivateKey),
		}
		if err = func() error {
			mtx := mpc.ContextWith(ctx)
			mtx.GetAttachments()["mesh.net.weave"] = "guest"
			mtx.GetPrincipals().Push(&types.Principal{NodeId: edge.NodeId})
			defer func() { mtx.GetPrincipals().Pop() }()
			return aware.RemoteNet.Weave(mtx, that.ExchangeRoute(ctx, edge.ExpireAt, tool.NextID(), certification, edge.Extra))
		}(); nil != err {
			return cause.Error(err)
		}
		edge.HostRoot = environ.RootCrt
		edge.HostCrt = certification.HostCrt
		edge.HostKey = certification.HostKey
		edge.Status = edge.Status | int32(types.Approving)
		if err = aware.LocalNet.Weave(ctx, edge); nil != err {
			return cause.Error(err)
		}
		return cause.Error(tabledataCaster.RouteManRefresh(ctx))
	}
	if "" == route.NodeId || "" == route.InstId || "" == route.Name || "" == route.URC().String() || "" == route.HostKey || "" == route.HostCrt || "" == route.AuthCode {
		log.Warn(ctx, "%s|%s|%s|%s|%s|%s|%s", route.NodeId, route.InstId, route.Name, route.URC().String(), route.HostKey, route.HostCrt, route.AuthCode)
		return cause.Errorf("Network weave must have enough parameters.")
	}
	edge.AuthCode = tool.Anyone(edge.AuthCode, route.AuthCode)
	edge.HostRoot = environ.RootCrt
	edge.HostCrt = sets.Get(types.IssueCrtKey)
	edge.HostKey = sets.Get(types.IssuePrivateKey)
	if "" == edge.GuestCrt || "" == edge.GuestKey {
		edge.GuestCrt = route.HostCrt
		edge.GuestKey = route.HostKey
	}
	if access, err := that.Accessible(ctx, edge); access {
		if nil != err {
			log.Warn(ctx, err.Error())
		}
		edge.Status = edge.Status | int32(types.Connected)
	}
	if err = aware.LocalNet.Weave(ctx, edge); nil != err {
		return cause.Error(err)
	}
	return cause.Error(tabledataCaster.RouteManRefresh(ctx))
}

func (that *PRSINetwork) Ack(ctx context.Context, route *types.Route) error {
	if "" == route.NodeId {
		return cause.Errorf("Network weave ack must bind nodeId.")
	}
	edge, err := aware.LocalNet.GetRoute(ctx, route.NodeId)
	if nil != err {
		return cause.Error(err)
	}
	if nil == edge {
		return cause.Errorf("Network weave ack must after send weaving.")
	}
	if mpc.ContextWith(ctx).GetAttachments()["mesh.net.ack"] != "guest" {
		certification := new(types.RouteCertificate)
		certification.HostCrt = edge.HostCrt
		certification.HostKey = edge.HostKey
		if err = func() error {
			mtx := mpc.ContextWith(ctx)
			mtx.GetAttachments()["mesh.net.ack"] = "guest"
			mtx.GetPrincipals().Push(&types.Principal{NodeId: edge.NodeId})
			defer func() { mtx.GetPrincipals().Pop() }()
			return aware.RemoteNet.Ack(mtx, that.ExchangeRoute(ctx, edge.ExpireAt, edge.AuthCode, certification, edge.Extra))
		}(); nil != err {
			return cause.Error(err)
		}
		if err = aware.LocalNet.Ack(ctx, edge); nil != err {
			return cause.Error(err)
		}
		return cause.Error(tabledataCaster.RouteManRefresh(ctx))
	}
	if "" == route.NodeId || "" == route.HostKey || "" == route.HostCrt || "" == route.AuthCode {
		log.Info(ctx, "%s|%s|%s|%s", route.NodeId, route.HostKey, route.HostCrt, route.AuthCode)
		return cause.Errorf("Network weave must have enough parameters.")
	}
	if "" == edge.GuestCrt || "" == edge.GuestKey {
		edge.GuestCrt = route.HostCrt
		edge.GuestKey = route.HostKey
	}
	edge.AuthCode = route.AuthCode
	if err = aware.LocalNet.Ack(ctx, edge); nil != err {
		return cause.Error(err)
	}
	return cause.Error(tabledataCaster.RouteManRefresh(ctx))
}

func (that *PRSINetwork) Disable(ctx context.Context, nodeId string) error {
	env, err := that.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	if tool.IsLocalEnv(env, nodeId) {
		return cause.Errorable(cause.Validate)
	}
	if err := aware.LocalNet.Disable(ctx, nodeId); nil != err {
		return cause.Error(err)
	}
	return cause.Error(tabledataCaster.RouteManRefresh(ctx))
}

func (that *PRSINetwork) Enable(ctx context.Context, nodeId string) error {
	if err := aware.LocalNet.Enable(ctx, nodeId); nil != err {
		return cause.Error(err)
	}
	return cause.Error(tabledataCaster.RouteManRefresh(ctx))
}

func (that *PRSINetwork) Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Route], error) {
	return aware.LocalNet.Index(ctx, index)
}

func (that *PRSINetwork) Version(ctx context.Context, nodeId string) (*types.Versions, error) {
	environ, err := that.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	if "" == nodeId || tool.IsLocalEnv(environ, nodeId) {
		registrations, err := aware.Registry.Export(ctx, types.METADATA)
		if nil != err {
			return nil, cause.Error(err)
		}
		infos := map[string]string{}
		for _, registration := range registrations {
			for key, value := range registration.Attachments {
				infos[key] = value
			}
		}
		infos[fmt.Sprintf("mesh.%s", types.CommitID)] = prsim.CommitID
		infos[fmt.Sprintf("mesh.%s", types.OS)] = prsim.GOOS
		infos[fmt.Sprintf("mesh.%s", types.ARCH)] = prsim.GOARCH
		infos[fmt.Sprintf("mesh.%s", types.Version)] = prsim.Version
		return &types.Versions{
			Version: prsim.Version,
			Infos:   infos,
		}, nil
	}
	return aware.RemoteNet.Version(mpc.ContextWith(ctx), nodeId)
}

func (that *PRSINetwork) Instx(ctx context.Context, index *types.Paging) (*types.Page[*types.Institution], error) {
	return aware.LocalNet.Instx(ctx, index)
}

func (that *PRSINetwork) Instr(ctx context.Context, institutions []*types.Institution) error {
	return aware.LocalNet.Instr(ctx, institutions)
}

func (that *PRSINetwork) Ally(ctx context.Context, nodeIds []string) error {
	return that.operateAlly(ctx, nodeIds, int32(types.Connected)|int32(types.Weaved))
}

func (that *PRSINetwork) Disband(ctx context.Context, nodeIds []string) error {
	return that.operateAlly(ctx, nodeIds, int32(types.Removed))
}

func (that *PRSINetwork) Assert(ctx context.Context, feature string, nodeIds []string) (bool, error) {
	sort.Strings(nodeIds)
	key := fmt.Sprintf("%s.%s", NetworkFeaturePrefix, feature)
	name := strings.Join(nodeIds, "|")
	ent, err := aware.Cache.HGet(ctx, key, name)
	if nil != err {
		return false, cause.Error(err)
	}
	if nil != err {
		var ok bool
		if err = ent.Entity.TryReadObject(&ok); nil != err {
			return false, cause.Error(err)
		}
		return ok, nil
	}
	ret, err := that.featureAssert(ctx, feature, nodeIds)
	if nil != err {
		return false, cause.Error(err)
	}
	entity, err := new(types.CacheEntity).Wrap(name, ret, time.Hour*6)
	if nil != err {
		return false, cause.Error(err)
	}
	if err = aware.Cache.HSet(ctx, key, entity); nil != err {
		return false, cause.Error(err)
	}
	return ret, nil
}

func (that *PRSINetwork) featureAssert(ctx context.Context, feature string, nodeIds []string) (bool, error) {
	getter := func(nodeId string) ([]*types.Route, error) {
		mtx := mpc.ContextWith(ctx).Resume(ctx)
		mtx.GetPrincipals().Push(&types.Principal{NodeId: nodeId, InstId: nodeId})
		defer func() {
			mtx.GetPrincipals().Pop()
		}()
		return aware.RemoteNet.GetRoutes(mtx)
	}
	intersections := map[string]int{}
	for _, nodeId := range nodeIds {
		routes, err := getter(nodeId)
		if nil != err {
			return false, cause.Error(err)
		}
		for _, route := range routes {
			intersections[route.NodeId] = intersections[route.NodeId] + 1
		}
	}
	return false, nil
}

func (that *PRSINetwork) GetRemoteAddr(ctx context.Context, tar *types.Route) (types.URC, error) {
	mtx := mpc.ContextWith(ctx).Resume(ctx)
	mtx.GetPrincipals().Push(&types.Principal{NodeId: tar.NodeId, InstId: tar.InstId})
	defer func() {
		mtx.GetPrincipals().Pop()
	}()
	route, err := aware.RemoteNet.GetRoute(mtx, tar.NodeId)
	if nil != err {
		return "", cause.Error(err)
	}
	if nil == route {
		return "", nil
	}
	return route.URC(), nil
}

func (that *PRSINetwork) FlushProxy(ctx context.Context, route *types.Route) error {
	router, err := that.GetRoute(ctx, tool.Anyone(route.NodeId, route.InstId))
	if nil != err {
		return cause.Error(err)
	}
	if nil == router {
		if "" == route.NodeId || "" == route.InstId || "" == route.URC().String() {
			log.Warn(ctx, "%s|%s|%s|%s", route.NodeId, route.InstId, route.URC().String(), route.InstName)
			return cause.ValidateErrorf("Network accessible check, nodeId, instId, address are required.")
		}
		if err = tabledataCaster.RouteAutoRefresh(ctx, route); nil != err {
			return cause.Error(err)
		}
		if err = proxyFlush(ctx); nil != err {
			return cause.Error(err)
		}
	}
	return nil
}

func (that *PRSINetwork) ExchangeRoute(ctx context.Context, expireAt int64, code string, certificate *types.RouteCertificate, extra string) *types.Route {
	environ, err := that.GetEnviron(ctx)
	if nil != err {
		log.Error(ctx, err.Error())
		environ = &types.Environ{}
	}
	me, err := that.GetRoute(ctx, environ.NodeId)
	if nil != err {
		log.Error(ctx, err.Error())
	}
	if nil == me {
		me = &types.Route{Address: tool.Address.Get().Any()}
	}
	return &types.Route{
		NodeId:    environ.NodeId,
		InstId:    environ.InstId,
		Name:      environ.InstName,
		InstName:  environ.InstName,
		Address:   me.Address,
		Describe:  environ.InstName,
		HostRoot:  certificate.HostRoot,
		HostKey:   certificate.HostKey,
		HostCrt:   certificate.HostCrt,
		GuestRoot: certificate.GuestRoot,
		GuestKey:  certificate.GuestKey,
		GuestCrt:  certificate.GuestCrt,
		Status:    0,
		Version:   0,
		AuthCode:  code,
		ExpireAt:  expireAt,
		Extra:     extra,
		CreateAt:  "",
		UpdateAt:  "",
		CreateBy:  environ.InstName,
		UpdateBy:  environ.InstName,
	}
}

func (that *PRSINetwork) operateAlly(ctx context.Context, nodeIds []string, status int32) error {
	routes, err := that.GetRoutes(ctx)
	if nil != err {
		return cause.Error(err)
	}
	env, err := that.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	var members []*types.Route
	var address = map[string]string{}
	for _, route := range routes {
		if tool.IsLocalEnv(env, route.NodeId, route.InstId) || types.Removed.Is(route.Status) || types.Disabled.Is(route.Status) {
			log.Warn(ctx, "Node %s has been removed, disabled or islocal, dont broadcast routes change event.", route.NodeId)
			continue
		}
		if tool.Contains(nodeIds, route.NodeId) || tool.Contains(nodeIds, route.InstId) {
			members = append(members, route)
			if addr, thx := that.GetRemoteAddr(ctx, route); nil != thx {
				return cause.Error(thx)
			} else {
				address[route.NodeId] = addr.String()
			}
		}
	}
	if len(members) < 2 {
		return cause.Errorable(cause.Validate)
	}
	ally := func(tar *types.Principal, rs []*types.Route) error {
		mtx := mpc.ContextWith(ctx).Resume(ctx)
		mtx.GetAttachments()["omega.tenant.id"] = tar.InstId
		mtx.GetAttachments()["omega.inst.id"] = tar.InstId
		mtx.GetPrincipals().Push(tar)
		defer func() {
			mtx.GetPrincipals().Pop()
		}()
		return aware.RemoteNet.Refresh(mtx, rs)
	}
	filter := func(nodeId string, rs []*types.Route) []*types.Route {
		var ss []*types.Route
		for _, route := range rs {
			payload := route.Copy()
			payload.Address = address[route.NodeId]
			payload.Group = env.InstId
			payload.AuthCode = tool.Anyone(payload.AuthCode, tool.NextID())
			if nodeId != route.NodeId && nodeId != route.InstId {
				payload.Status = types.Virtual.Or(status)
			}
			ss = append(ss, payload)
		}
		return ss
	}
	for _, member := range members {
		tar := &types.Principal{NodeId: member.NodeId, InstId: member.InstId}
		if err = ally(tar, filter(member.NodeId, members)); nil != err {
			return cause.Error(err)
		}
	}
	return nil
}
