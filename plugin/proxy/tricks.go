/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"os"
	"strings"
	"time"
)

func init() {
	var _ prsim.Listener = tricks
	macro.Provide(prsim.IListener, tricks)

	var _ prsim.EndpointSticker[map[string]string, []*types.Trick] = new(trickEndpoint)
	macro.Provide(prsim.IEndpointSticker, new(trickEndpoint))
}

const (
	trickKey    = "mesh.proxy.trick"
	trickFormat = "Mesh trick format must be name:kind+service=proto://address/uname,uname"
)

var tricks = new(trick)

type trick struct {
}

func (that *trick) Att() *macro.Att {
	return &macro.Att{Name: "mesh.registry.event.trick"}
}

func (that *trick) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.ProxyRegisterEvent}
}

func (that *trick) Listen(ctx context.Context, event *types.Event) error {
	that.legacyAsDomain(ctx)
	environ, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	ent, err := aware.KV.Get(ctx, trickKey)
	if nil != err {
		return cause.Error(err)
	}
	tks := Tricks{}
	if err = tks.Decode(ent); nil != err {
		return cause.Error(err)
	}
	if err = tks.Put(os.Getenv("MESH_TRICK")); nil != err {
		return cause.Error(err)
	}
	if len(tks) < 1 {
		return nil
	}
	rs := map[string]*types.MetadataRegistration{}
	for _, v := range tks {
		k := fmt.Sprintf("%s:%s", v.Name, v.Kind)
		if nil == rs[k] {
			rs[k] = &types.MetadataRegistration{
				InstanceId:  v.Address,
				Name:        v.Name,
				Kind:        v.Kind,
				Address:     v.Address,
				Timestamp:   time.Now().Add(time.Minute * 10).UnixMilli(),
				Content:     &types.Metadata{},
				Attachments: map[string]string{},
			}
		}
		rs[k].Content.Services = append(rs[k].Content.Services, that.trickServices(environ, v)...)
	}
	fn := func(v *types.MetadataRegistration) *types.Registration[any] { return v.Any() }
	return cause.Error(aware.LocalRegistry.Registers(ctx, tool.MapValues(fn, rs)))
}

func (that *trick) legacyAsDomain(ctx context.Context) {
	routes, err := aware.Network.GetRoutes(ctx)
	if nil != err {
		log.Warn(ctx, err.Error())
		return
	}
	var domains []*types.Domain
	for _, edge := range routes {
		for _, name := range []string{"tensor.route.grpc"} {
			for _, nodeId := range []string{edge.NodeId, edge.InstId, strings.ToLower(edge.NodeId), strings.ToLower(edge.InstId)} {
				for _, domain := range []string{types.MeshDomain, fmt.Sprintf("%s.net", types.CN)} {
					urn := &types.URN{
						Domain: domain,
						NodeId: nodeId,
						Flag: &types.URNFlag{
							V:       "00",
							Proto:   mpc.MeshFlag.GRPC.Code(),
							Codec:   mpc.MeshFlag.JSON.Code(),
							Version: "000000",
							Zone:    "00",
							Cluster: "00",
							Cell:    "00",
							Group:   "00",
							Address: "000000000000",
							Port:    "00000",
						},
						Name: name,
					}
					domains = append(domains, &types.Domain{URN: urn.String(), Address: tool.Address.Get().Any()})
				}
			}
		}
	}
	if err = aware.Network.PutDomains(ctx, prsim.AutoDomain, domains); nil != err {
		log.Error(ctx, err.Error())
	}
}

// trickServices name:dns=uname,
func (that *trick) trickServices(environ *types.Environ, vec *types.Trick) []*types.Service {
	var services []*types.Service
	for _, name := range vec.Patterns {
		urn := &types.URN{
			Domain: types.MeshDomain,
			NodeId: environ.NodeId,
			Flag: &types.URNFlag{
				V:       "00",
				Proto:   mpc.MeshFlag.OfName(vec.Proto).Code(),
				Codec:   mpc.MeshFlag.JSON.Code(),
				Version: "000000",
				Zone:    "00",
				Cluster: "00",
				Cell:    "00",
				Group:   "00",
				Address: "000000000000",
				Port:    "00000",
			},
			Name: name,
		}
		services = append(services, &types.Service{
			URN:       tool.Ternary(types.METADATA == vec.Kind, urn.String(), name),
			Namespace: "",
			Name:      name,
			Version:   "1.0.0",
			Proto:     mpc.MeshFlag.OfName(vec.Proto).Code(),
			Codec:     mpc.MeshFlag.JSON.Name(),
			Flags:     0,
			Timeout:   10000,
			Retries:   3,
			Node:      environ.NodeId,
			Inst:      environ.InstId,
			Zone:      "",
			Cluster:   "",
			Cell:      "",
			Group:     "",
			Address:   vec.Address,
			Kind:      vec.Service,
			Lang:      "Any",
			Attrs:     map[string]string{"plugins": vec.Plugins, "priority": "high"},
		})
	}
	return services
}

type trickEndpoint struct {
}

func (that *trickEndpoint) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.trick"}
}

func (that *trickEndpoint) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.trick"}
}

func (that *trickEndpoint) I() map[string]string {
	return map[string]string{}
}

func (that *trickEndpoint) O() []*types.Trick {
	return []*types.Trick{}
}

func (that *trickEndpoint) Stick(ctx context.Context, varg map[string]string) ([]*types.Trick, error) {
	remove := tool.MapBy(varg, "r")
	write := tool.MapBy(varg, "w")
	ent, err := aware.KV.Get(ctx, trickKey)
	if nil != err {
		return nil, cause.Error(err)
	}
	tks := Tricks{}
	if err = tks.Decode(ent); nil != err {
		return nil, cause.Error(err)
	}
	if "" == remove && "" == write {
		return tool.Values(tks), nil
	}
	// name:kind+service
	if "" != remove {
		delete(tks, remove)
	}
	// name:kind+service=proto://address/uname,uname
	if err = tks.Put(write); nil != err {
		return nil, cause.Error(err)
	}
	n, err := new(types.Entity).Wrap(tks)
	if nil != err {
		return nil, cause.Error(err)
	}
	if err = aware.KV.Put(ctx, trickKey, n); nil != err {
		return nil, cause.Error(err)
	}
	return tool.Values(tks), cause.Error(aware.Scheduler.Emit(ctx, new(types.Topic).With(prsim.ProxyRegisterEvent)))
}

type Tricks map[string]*types.Trick

func (that Tricks) Decode(ent *types.Entity) error {
	if nil == ent || !ent.Present() {
		return nil
	}
	return cause.Error(ent.TryReadObject(&that))
}

func (that Tricks) Put(tricks string) error {
	for _, k := range strings.Split(strings.TrimSpace(tricks), ";") {
		if "" == k {
			return nil
		}
		pair := strings.Split(strings.TrimSpace(k), "=")
		if len(pair) < 2 {
			return cause.ValidateErrorf(trickFormat)
		}
		nks := strings.Split(pair[0], ":")
		if len(nks) < 2 {
			nks = append(nks, types.METADATA)
		}
		ks := strings.Split(nks[1], "+")
		if len(ks) < 2 {
			ks = append(ks, types.MPS)
		}
		uri, err := types.FormatURL(strings.Join(pair[1:], "="))
		if nil != err {
			return cause.Error(err)
		}
		tk := &types.Trick{
			Name:     nks[0],
			Kind:     ks[0],
			Service:  ks[1],
			Proto:    uri.Scheme,
			Patterns: strings.Split(strings.TrimPrefix(uri.Path, "/"), ","),
			Address:  uri.Host,
			Plugins:  strings.TrimSpace(uri.Query().Get("p")),
		}
		that[fmt.Sprintf("%s:%s+%s", tk.Name, tk.Kind, tk.Service)] = tk
	}
	return nil
}
