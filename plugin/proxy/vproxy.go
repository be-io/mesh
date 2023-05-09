/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v2/pkg/server/service/loadbalancer/wrr"
	"net/http"
	"strings"
)

func init() {
	var _ prsim.Listener = vp
	wrr.Provide(vp)
	macro.Provide(prsim.IListener, vp)

	var _ prsim.EndpointSticker[map[string]string, []*types.VIP] = new(vipEndpoint)
	macro.Provide(prsim.IEndpointSticker, new(vipEndpoint))
}

const (
	vipKey    = "mesh.proxy.vip"
	vipFormat = "Mesh vip format must be feature:matcher=label://host,host"
)

var vp = new(VProxyStrategy)

type VProxyStrategy struct {
	vips map[string]*types.VIP
}

func (that *VProxyStrategy) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.vip"}
}

func (that *VProxyStrategy) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.ProxyRegisterEvent}
}

func (that *VProxyStrategy) Listen(ctx context.Context, event *types.Event) error {
	if event.Binding.Match(prsim.ProxyRegisterEvent) {
		ent, err := aware.KV.Get(ctx, vipKey)
		if nil != err {
			return cause.Error(err)
		}
		if nil == ent || !ent.Present() {
			return nil
		}
		var vips map[string]*types.VIP
		if err = ent.TryReadObject(&vips); nil != err {
			return cause.Error(err)
		}
		that.vips = vips
	}
	return nil
}

func (that *VProxyStrategy) Name() string {
	return "vproxy@mesh"
}

func (that *VProxyStrategy) Priority() int {
	return 1
}

func (that *VProxyStrategy) Next(w http.ResponseWriter, req *http.Request, servers []wrr.Server) []wrr.Server {
	vips := that.vips
	if len(vips) < 1 {
		return servers
	}
	vk := that.vKey(req)
	var ss []wrr.Server
	for _, server := range servers {
		if nil == server.URL() {
			continue
		}
		for _, v := range vips {
			if v.Matches(vk, server) && tool.Contains(v.Hosts, server.URL().Hostname()) {
				ss = append(ss, server)
			}
		}
	}
	if len(ss) > 0 {
		return ss
	}
	return servers
}

func (that *VProxyStrategy) vKey(req *http.Request) *types.VIPKey {
	return &types.VIPKey{
		Version: prsim.MeshVersion.GetHeader(req.Header),
		IP:      tool.NewAddr(req.RemoteAddr, 80).Host,
		SrcID:   types.NodeSEQ(tool.Anyone(prsim.MeshFromNodeId.GetHeader(req.Header), prsim.MeshFromInstId.GetHeader(req.Header))),
		DstID:   types.NodeSEQ(types.FromURN(macro.Context(), tool.Anyone(prsim.MeshUrn.GetHeader(req.Header), req.Host)).NodeId),
		Gray:    prsim.MeshGray.GetHeader(req.Header),
	}
}

type vipEndpoint struct {
}

func (that *vipEndpoint) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.vip"}
}

func (that *vipEndpoint) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.vip"}
}

func (that *vipEndpoint) I() map[string]string {
	return map[string]string{}
}

func (that *vipEndpoint) O() []*types.VIP {
	return []*types.VIP{}
}

func (that *vipEndpoint) Stick(ctx context.Context, varg map[string]string) ([]*types.VIP, error) {
	remove := tool.MapBy(varg, "r")
	write := tool.MapBy(varg, "w")
	ent, err := aware.KV.Get(ctx, vipKey)
	if nil != err {
		return nil, cause.Error(err)
	}
	vips := map[string]*types.VIP{}
	if nil != ent && ent.Present() {
		if err = ent.TryReadObject(&vips); nil != err {
			return nil, cause.Error(err)
		}
	}
	if "" == remove && "" == write {
		return tool.Values(vips), nil
	}
	if "" != remove {
		delete(vips, remove)
	}
	if "" != write {
		kv := strings.Split(write, "=")
		if len(kv) < 2 {
			return nil, cause.ValidateErrorf(vipFormat)
		}
		kvs := strings.Split(kv[0], ":")
		if len(kvs) < 2 {
			return nil, cause.ValidateErrorf(vipFormat)
		}
		lvs := strings.Split(kv[1], "://")
		if len(lvs) < 2 {
			return nil, cause.Errorf(vipFormat)
		}
		vips[kv[0]] = &types.VIP{
			Name:    kvs[0],
			Matcher: kvs[1],
			Label:   lvs[0],
			Hosts:   strings.Split(lvs[1], ","),
		}
	}
	n, err := new(types.Entity).Wrap(vips)
	if nil != err {
		return nil, cause.Error(err)
	}
	if err = aware.KV.Put(ctx, vipKey, n); nil != err {
		return nil, cause.Error(err)
	}
	return tool.Values(vips), cause.Error(aware.Scheduler.Emit(ctx, new(types.Topic).With(prsim.ProxyRegisterEvent)))
}
