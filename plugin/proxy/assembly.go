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
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"strconv"
	"strings"
)

func init() {
	var _ prsim.Listener = assemblies
	macro.Provide(prsim.IListener, assemblies)
}

var rewritePaths = []string{
	"/mesh/invoke",
}

var assemblies = new(assembly)

type assembly struct {
	Servers map[string][]string
}

func (that *assembly) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.assembly"}
}

func (that *assembly) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.RegistryEventRefresh}
}

func (that *assembly) Listen(ctx context.Context, event *types.Event) error {
	var metadata types.MetadataRegistrations
	if err := event.TryGetObject(&metadata); nil != err {
		return cause.Error(err)
	}
	servers := map[string][]string{}
	for _, meta := range metadata.Of(types.METADATA) {
		host, port := that.ParseHost(ctx, meta.Address)
		subset := prsim.MeshSubset.Get(meta.Attachments)
		if "" != subset {
			servers[meta.Name] = append(servers[meta.Name], fmt.Sprintf("%s:%d?s=%s", host, that.SelectPort(ctx, port, meta.Attachments), subset))
		} else {
			servers[meta.Name] = append(servers[meta.Name], fmt.Sprintf("%s:%d", host, that.SelectPort(ctx, port, meta.Attachments)))
		}
	}
	that.Servers = servers
	return nil
}

// SelectPort For ICBC
func (that *assembly) SelectPort(ctx context.Context, dft int, attachments map[string]string) int {
	if nil == attachments || "" == attachments["_PAAS_PORT_7700"] {
		return dft
	}
	port, err := strconv.Atoi(attachments["_PAAS_PORT_7700"])
	if nil != err {
		log.Warn(ctx, err.Error())
		return dft
	}
	return port
}

func (that *assembly) ParseHost(ctx context.Context, address string) (string, int) {
	pair := strings.Split(address, ":")
	if len(pair) < 2 {
		return pair[0], 80
	}
	if port, err := strconv.Atoi(pair[1]); nil != err {
		log.Error(ctx, err.Error())
		return pair[0], 80
	} else {
		if port == 8864 {
			return pair[0], 8865
		}
		//jewel-connector
		if port == 9513 {
			return pair[0], 9512
		}
		if port%100 >= 10 || port == 9904 || port == 9902 || port == 9906 {
			return pair[0], port
		}
		return pair[0], port / 100 * 100
	}
}

func (that *assembly) WithAsmRoute(ctx context.Context, routes *Routes) {
	for name, servers := range that.Servers {
		var streams []dynamic.Server
		for _, server := range servers {
			streams = append(streams, dynamic.Server{URL: fmt.Sprintf("%s://%s", "http", server)})
		}
		ns := fmt.Sprintf("%s#%s#asm", types.LocalNodeId, name)
		routes.services[ns] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: streams}}
		//
		tr := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s\\.route\\.studio\\.([-a-z0-9]+)\\.(%s|%s|%s|%s)\\.%s\\.(net|cn|com)`)", name, types.LocalNodeId, types.LocalInstId, strings.ToLower(types.LocalNodeId), strings.ToLower(types.LocalInstId), types.CN)
		routes.Route(ctx, fmt.Sprintf("%s#%s#asms#insecure", types.LocalNodeId, name), &dynamic.Router{
			EntryPoints: []string{TransportA, TransportB},
			Middlewares: []string{PluginBarrier},
			Service:     ns,
			Rule:        tr,
			Priority:    500,
		})
		routes.Route(ctx, fmt.Sprintf("%s#%s#asms#secure", types.LocalNodeId, name), &dynamic.Router{
			EntryPoints: []string{TransportA, TransportB},
			Middlewares: []string{PluginBarrier},
			Service:     ns,
			Rule:        tr,
			Priority:    500,
			TLS: &dynamic.RouterTLSConfig{
				Options: types.LocalNodeId,
				Domains: proxy.Domains(),
			},
		})
		if tool.Name.Get() == name {
			that.AsmMesh(ctx, routes, ns, name)
		}
	}
}

func (that *assembly) AsmMesh(ctx context.Context, routes *Routes, ns string, name string) {
	routes.Route(ctx, fmt.Sprintf("%s#mesh#asm#insecure", types.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportA, TransportB, TransportX, TransportY},
		Middlewares: []string{PluginBarrier},
		Service:     ns,
		Rule:        "PathRegexp(`/mesh.*`) || PathRegexp(`/v1/interconn/.*`)",
		Priority:    200,
	})
	routes.Route(ctx, fmt.Sprintf("%s#mesh#asm#secure", types.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportA, TransportB, TransportX, TransportY},
		Middlewares: []string{PluginBarrier},
		Service:     ns,
		Rule:        "PathRegexp(`/mesh.*`) || PathRegexp(`/v1/interconn/.*`)",
		Priority:    200,
		TLS: &dynamic.RouterTLSConfig{
			Options: types.LocalNodeId,
			Domains: proxy.Domains(),
		},
	})
	env, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		log.Warn(ctx, err.Error())
		return
	}
	cr := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", types.URNFlagMatcher("([-a-zA-Z\\d]+)", "([-a-zA-Z\\d]{2}00[-a-zA-Z\\d]+)", env.ID(ctx).SEQ, types.CN, env.NodeId, env.InstId, types.LocalNodeId, types.LocalInstId))
	routes.Route(ctx, fmt.Sprintf("%s#mesh#asmh#insecure", types.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportA, TransportB},
		Middlewares: []string{PluginBarrier, PluginRewrite},
		Service:     ns,
		Rule:        cr,
		Priority:    300,
	})
	routes.Route(ctx, fmt.Sprintf("%s#mesh#asmh#secure", types.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportA, TransportB},
		Middlewares: []string{PluginBarrier, PluginRewrite},
		Service:     ns,
		Rule:        cr,
		Priority:    300,
		TLS: &dynamic.RouterTLSConfig{
			Options: types.LocalNodeId,
			Domains: proxy.Domains(),
		},
	})
}
