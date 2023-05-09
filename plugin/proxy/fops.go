/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"bytes"
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"regexp"
	"strings"
)

func compile(padding string) *regexp.Regexp {
	exp := &bytes.Buffer{}
	for index, p := range patterns {
		exp.WriteString(fmt.Sprintf("%s%s%s", tool.Ternary(0 == index, "", "|"), p, padding))
	}
	patt, err := regexp.Compile(exp.String())
	if nil != err {
		log.Error0("Compiling regex %s: %s", exp.String(), err.Error())
	}
	return patt
}

var patterns = []string{
	"href=\"/",
	"src=\"/",
	".href = \"/",
	".src = \"/",
	"baseAssets:\"/",
	"apiEndpoint:\"/",
	"projectEndpoint:\"/",
	"proxyEndpoint:\"/",
	"wsEndpoint:\"/",
	"catalogEndpoint:\"/",
	"authEndpoint:\"/",
	"telemetryEndpoint:\"/",
	"webhookEndpoint:\"/",
	"magicEndpoint:\"/",
	"kubernetesEndpoint:\"/",
	"kubectlEndpoint:\"/",
	"kubernetesDashboard:\"/",
	"mesosEndpoint:\"/",
	"swarmDashboard:\"/",
	"legacyApiEndpoint:\"/",
	"\":\"https?:\\/\\/[a-zA-Z0-9|\\.|-|_|:]+\\/",
	"\":\"https?:\\\\/\\\\/[a-zA-Z0-9|\\.|-|_|:]+\\\\/",
	"\":\"wss?:\\/\\/[a-zA-Z0-9|\\.|-|_|:]+\\/",
	"\":\"wss?:\\\\/\\\\/[a-zA-Z0-9|\\.|-|_|:]+\\\\/",
	"currentEndpoint\":\"端点 \\(",
	"legacyEndpoint\":\"端点 \\(",
	//
	"url:\"/",
	"\\.p=\"/",
	"url\\(\"/",
}

var pattern = compile("")

var tunnel = new(fops)

type FopsRule struct {
	Path    string `json:"path"`
	Server  string `json:"server"`
	Filter  string `json:"filter"`
	Replace string `json:"replace"`
}

type fops struct {
}

func (that *fops) WithFOPSRoute(ctx context.Context, routers []*types.Route, routes *Routes, snapshot *SnapShot) {
	ent, err := aware.KV.Get(ctx, "mesh.plugin.proxy.fops")
	if nil != err {
		log.Warn(ctx, err.Error())
	}
	var rs []*FopsRule
	if nil != ent {
		if err = ent.TryReadObject(&rs); nil != err {
			log.Warn(ctx, err.Error())
		}
	}
	var rewrites []*Rewrite
	if err = log.PError(ctx, func() error {
		for _, route := range routers {
			fs := routes.IfAbsent(fmt.Sprintf("%s#fops", route.NodeId), func(n string) *dynamic.Service {
				return snapshot.remoteService(ctx, n, routes, "https", route.NodeId, snapshot.supervise(route), route.URC().URL(ctx).SetD("1").URC())
			})
			ps := fmt.Sprintf("PathRegexp(`/%s/mesh/fops/*`)", route.ID(ctx).SEQ)
			routes.Route(ctx, fmt.Sprintf("%s#mesh#fops#insecure", route.NodeId), &dynamic.Router{
				EntryPoints: []string{TransportA, TransportB},
				Middlewares: []string{PluginBarrier, PluginForwarder, PluginSubfilter},
				Service:     fs,
				Rule:        ps,
				Priority:    2000,
			})
			routes.Route(ctx, fmt.Sprintf("%s#mesh#fops#secure", route.NodeId), &dynamic.Router{
				EntryPoints: []string{TransportA, TransportB},
				Middlewares: []string{PluginBarrier, PluginForwarder, PluginSubfilter},
				Service:     fs,
				Rule:        ps,
				Priority:    2000,
				TLS: &dynamic.RouterTLSConfig{
					Options: route.NodeId,
					Domains: proxy.Domains(),
				},
			})
			rewrites = append(rewrites, &Rewrite{
				regex:   pattern,
				prefix:  strings.TrimSpace(fmt.Sprintf("/%s", route.ID(ctx).SEQ)),
				filter:  nil,
				replace: "",
			})
		}
		if plugin.PROXY.Match() {
			return nil
		}
		for _, vs := range rs {
			if nil == vs || "" == vs.Path || "" == vs.Server {
				continue
			}
			ns := fmt.Sprintf("%s#%s#fops", types.LocalNodeId, vs.Path)
			addr := tool.Ternary(strings.Contains(vs.Server, "://"), vs.Server, fmt.Sprintf("%s://%s", "http", vs.Server))
			routes.services[ns] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: addr}}}}
			//
			routes.Route(ctx, fmt.Sprintf("%s#%s#fops#insecure", vs.Path, types.LocalNodeId), &dynamic.Router{
				EntryPoints: []string{TransportY},
				Middlewares: []string{PluginBarrier, PluginForwarder, PluginSubfilter},
				Service:     ns,
				Rule:        fmt.Sprintf("PathRegexp(`/mesh/fops/%s/*`)", vs.Path),
				Priority:    200,
			})
			routes.Route(ctx, fmt.Sprintf("%s#%s#fops#secure", vs.Path, types.LocalNodeId), &dynamic.Router{
				EntryPoints: []string{TransportY},
				Middlewares: []string{PluginBarrier, PluginForwarder, PluginSubfilter},
				Service:     ns,
				Rule:        fmt.Sprintf("PathRegexp(`/mesh/fops/%s/*`)", vs.Path),
				Priority:    200,
				TLS: &dynamic.RouterTLSConfig{
					Options: types.LocalNodeId,
					Domains: proxy.Domains(),
				},
			})
			rewrite := &Rewrite{
				regex:   pattern,
				prefix:  strings.TrimSpace(fmt.Sprintf("/mesh/fops/%s", vs.Path)),
				replace: vs.Replace,
			}
			if "" != vs.Filter {
				if rewrite.filter, err = regexp.Compile(vs.Filter); nil != err {
					log.Warn(ctx, err.Error())
				}
			}
			rewrites = append(rewrites, rewrite)
		}
		filter.Refresh(rewrites)
		return nil
	}); nil != err {
		log.Error(ctx, err.Error())
	}
}
