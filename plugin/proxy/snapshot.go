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
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	mtypes "github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/paerser/types"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/tls"
	"strings"
	"time"
)

type Routes struct {
	routers  map[string]*dynamic.Router
	services map[string]*dynamic.Service
}

type SnapShot struct {
	env     *mtypes.Environ
	routers []*mtypes.Route
	groups  map[string]*mtypes.Route
	proxies mtypes.MetadataRegistrations
	servers mtypes.MetadataRegistrations
	sets    mtypes.MetadataRegistrations
	complex mtypes.MetadataRegistrations
	lic     *mtypes.License
}

func (that *Routes) IfAbsent(name string, fn func(n string) *dynamic.Service) string {
	if nil == that.services {
		that.services = map[string]*dynamic.Service{}
	}
	if nil == that.services[name] {
		that.services[name] = fn(name)
	}
	return name
}

func (that *Routes) Route(ctx context.Context, name string, router *dynamic.Router) {
	if nil != that.routers[name] {
		log.Warn(ctx, "Route %s overlap write.", name)
	}
	that.routers[name] = router
}

func (that *SnapShot) remoteService(ctx context.Context, name string, routes *Routes, schema string, nodeId string, uris ...mtypes.URC) *dynamic.Service {
	if len(that.proxies) > 0 {
		var servers []dynamic.Server
		for _, px := range that.proxies {
			address := tool.Ternary(strings.Contains(px.Address, "://"), px.Address, fmt.Sprintf("%s://%s", tool.Anyone(schema, "https"), px.Address))
			servers = append(servers, dynamic.Server{
				URL: mtypes.ParseURL(ctx, address).SetP("").SetK("").SetT("").String(),
			})
		}
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: nodeId, Servers: servers}}
	}
	urc := tool.Anyone(uris...)
	paas := tool.Ternary("" != urc.URL(ctx).GetD(), true, false)
	members, others := urc.Failover()
	var servers []dynamic.Server
	for _, addr := range members {
		address := tool.Ternary(strings.Contains(addr, "://"), addr, fmt.Sprintf("%s://%s", tool.Anyone(schema, "https"), addr))
		servers = append(servers, dynamic.Server{
			URL: mtypes.ParseURL(ctx, address).SetP("").SetK("").SetT("").String(),
		})
	}
	if len(others) < 1 {
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{
			ServersTransport: nodeId,
			Servers:          servers,
			PassHostHeader:   &paas,
		}}
	}
	var failoverServers []dynamic.Server
	for _, addr := range others {
		address := tool.Ternary(strings.Contains(addr, "://"), addr, fmt.Sprintf("%s://%s", tool.Anyone(schema, "https"), addr))
		failoverServers = append(failoverServers, dynamic.Server{
			URL: mtypes.ParseURL(ctx, address).SetP("").SetK("").SetT("").String(),
		})
	}
	master := routes.IfAbsent(fmt.Sprintf("%s#failover", name), func(n string) *dynamic.Service {
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{
			ServersTransport: nodeId,
			Servers:          servers,
			PassHostHeader:   &paas,
			HealthCheck: &dynamic.ServerHealthCheck{
				Path:     "/stats",
				Interval: types.Duration(12 * time.Second),
				Timeout:  types.Duration(10 * time.Second),
			},
		}}
	})
	slave := routes.IfAbsent(fmt.Sprintf("%s#failover", name), func(n string) *dynamic.Service {
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{
			ServersTransport: nodeId,
			Servers:          failoverServers,
			PassHostHeader:   &paas,
			HealthCheck: &dynamic.ServerHealthCheck{
				Path:     "/stats",
				Interval: types.Duration(12 * time.Second),
				Timeout:  types.Duration(10 * time.Second),
			},
		}}
	})
	return &dynamic.Service{Failover: &dynamic.Failover{Service: master, Fallback: slave, HealthCheck: &dynamic.HealthCheck{}}}
}

func (that *SnapShot) withRemoteRoute(ctx context.Context, routes *Routes) {
	for _, route := range that.routers {
		if tool.IsLocalEnv(that.env, route.NodeId) || proxy.RecursionBreak(ctx, route.URC()) {
			continue
		}
		if "" != mtypes.ParseURL(ctx, route.URC().String()).GetT() {
			ts := routes.IfAbsent(fmt.Sprintf("%s#tensor", route.NodeId), func(n string) *dynamic.Service {
				return that.remoteService(ctx, n, routes, "https", route.NodeId, that.supervise(route), route.URC())
			})
			tr := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", mtypes.URNMatcher("tensor", route.ID(ctx).SEQ, mtypes.CN, route.NodeId, route.InstId))
			routes.Route(ctx, fmt.Sprintf("%s#tensor#secure", route.NodeId), &dynamic.Router{
				EntryPoints: []string{TransportX, TransportY},
				Middlewares: []string{PluginBarrier, PluginForwarder},
				Service:     ts,
				Rule:        tr,
				Priority:    2500,
				TLS: &dynamic.RouterTLSConfig{
					Options: route.NodeId,
					Domains: proxy.Domains(),
				},
			})
			routes.Route(ctx, fmt.Sprintf("%s#tensor#insecure", route.NodeId), &dynamic.Router{
				EntryPoints: []string{TransportX, TransportY},
				Middlewares: []string{PluginBarrier, PluginForwarder},
				Service:     ts,
				Rule:        tr,
				Priority:    2500,
			})
		}
		rs := routes.IfAbsent(route.NodeId, func(n string) *dynamic.Service {
			return that.remoteService(ctx, n, routes, "https", route.NodeId, that.supervise(route), route.URC())
		})
		rr := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", mtypes.URNMatcher("([-a-zA-Z\\d]+)", route.ID(ctx).SEQ, mtypes.CN, route.NodeId, route.InstId))
		routes.Route(ctx, fmt.Sprintf("%s#secure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier, PluginForwarder},
			Service:     rs,
			Rule:        rr,
			Priority:    2000,
			TLS: &dynamic.RouterTLSConfig{
				Options: route.NodeId,
				Domains: proxy.Domains(),
			},
		})
		routes.Route(ctx, fmt.Sprintf("%s#insecure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier, PluginForwarder},
			Service:     rs,
			Rule:        rr,
			Priority:    2000,
		})
	}
}

// supervise redirect request to spec inst if in union
func (that *SnapShot) supervise(ne *mtypes.Route) mtypes.URC {
	if tool.Contains(that.lic.Group, that.env.NodeId) || tool.Contains(that.lic.Group, that.env.InstId) {
		return ""
	}
	groups := strings.Split(ne.Group, ",")
	for _, group := range groups {
		if nil != that.groups[group] {
			return that.groups[group].URC()
		}
	}
	return ""
}

func (that *SnapShot) withClusterSetRoute(ctx context.Context, routes *Routes, name string, flag string,
	instances mtypes.MetadataRegistrations, priority int, uri func(registration *mtypes.MetadataRegistration) string,
	pattern func() string) {
	cs := routes.IfAbsent(fmt.Sprintf("%s#%s#%s", mtypes.LocalNodeId, name, flag), func(n string) *dynamic.Service {
		var servers []dynamic.Server
		for _, registration := range instances {
			if proxy.RecursionBreak(ctx, mtypes.URC(registration.Address)) {
				continue
			}
			if !plugin.PROXY.Match() {
				servers = append(servers, dynamic.Server{URL: uri(registration)})
				continue
			}
			for _, addr := range tool.Address.Get().All() {
				servers = append(servers, dynamic.Server{URL: fmt.Sprintf("https://%s", addr)})
			}
		}
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: mtypes.LocalNodeId, Servers: servers}}
	})
	cr := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", pattern())
	routes.Route(ctx, fmt.Sprintf("%s#secure", cs), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     cs,
		Rule:        cr,
		Priority:    priority,
		TLS: &dynamic.RouterTLSConfig{
			Options: mtypes.LocalNodeId,
			Domains: proxy.Domains(),
		},
	})
	routes.Route(ctx, fmt.Sprintf("%s#insecure", cs), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     cs,
		Rule:        cr,
		Priority:    priority,
	})
}

func (that *SnapShot) withClusterRoute(ctx context.Context, routes *Routes) {
	nvs := map[string]mtypes.MetadataRegistrations{}
	for _, registration := range that.sets {
		name := formatLegacyName(registration.Name)
		nvs[name] = append(nvs[name], registration)
	}
	for name, instances := range nvs {
		that.withClusterSetRoute(ctx, routes, name, "h2", instances, 1000, func(registration *mtypes.MetadataRegistration) string {
			return mtypes.ParseURL(ctx, registration.Address).SetS(prsim.MeshSubset.Get(registration.Attachments)).SetSchema("h2c").URL().String()
		}, func() string {
			return mtypes.URNMatcher(name, that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId)
		})
		that.withClusterSetRoute(ctx, routes, name, "h1", instances, 1200, func(registration *mtypes.MetadataRegistration) string {
			h, p := assemblies.ParseHost(ctx, registration.Address)
			uri := mtypes.ParseURL(ctx, registration.Address).SetS(prsim.MeshSubset.Get(registration.Attachments)).SetSchema("http").SetHost(fmt.Sprintf("%s:%d", h, p)).URL().String()
			return uri
		}, func() string {
			return mtypes.URNMatcher(fmt.Sprintf("h1\\.%s", name), that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId)
		})
	}
	css := map[string]mtypes.MetadataRegistrations{}
	for _, registration := range that.sets {
		for _, service := range registration.InferService() {
			if (service.Flags & 4) == 4 {
				urn := mtypes.FromURN(ctx, service.URN)
				css[urn.Name] = append(css[urn.Name], registration)
			}
		}
	}
	for name, instances := range css {
		that.withClusterSetRoute(ctx, routes, tool.Hash(name), "c2", instances, 1500, func(registration *mtypes.MetadataRegistration) string {
			return mtypes.ParseURL(ctx, fmt.Sprintf("h2c://%s", registration.Address)).SetS(prsim.MeshSubset.Get(registration.Attachments)).URL().String()
		}, func() string {
			return mtypes.URNPatternMatcher(strings.Join(tool.Reverse(mtypes.AsArray(name)), "."), that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId)
		})
	}
}

func (that *SnapShot) bridgeService(ctx context.Context, routes *Routes) string {
	return routes.IfAbsent(fmt.Sprintf("%s#h12", mtypes.LocalNodeId), func(n string) *dynamic.Service {
		var servers []dynamic.Server
		if plugin.PROXY.Match() {
			for _, addr := range tool.Address.Get().All() {
				servers = append(servers, dynamic.Server{URL: fmt.Sprintf("https://%s", addr)})
			}
			return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: mtypes.LocalNodeId, Servers: servers}}
		}
		servers = append(servers, dynamic.Server{URL: fmt.Sprintf("%s://127.0.0.1:8865", "http")})
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: mtypes.LocalNodeId, Servers: servers}}
	})
}

func (that *SnapShot) withBridgeRoute(ctx context.Context, routes *Routes) {
	for _, route := range that.routers {
		if tool.IsLocalEnv(that.env, route.NodeId) || proxy.RecursionBreak(ctx, route.URC()) || "1" != mtypes.ParseURL(ctx, route.URC().String()).GetH() {
			continue
		}
		h21s := routes.IfAbsent(fmt.Sprintf("%s#h21", route.NodeId), func(n string) *dynamic.Service {
			return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{
				ServersTransport: mtypes.LocalNodeId,
				Servers:          []dynamic.Server{{URL: fmt.Sprintf("h2c://127.0.0.1:%d", tool.Runtime.Get().Port)}}}}
		})
		h21r := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", mtypes.URNMatcher("([-a-zA-Z\\d]+)", route.ID(ctx).SEQ, mtypes.CN, route.NodeId, route.InstId))
		routes.Route(ctx, fmt.Sprintf("%s#h21#secure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier, PluginForwarder, PluginHath, PluginReplace},
			Service:     h21s,
			Rule:        h21r,
			Priority:    1500,
			TLS: &dynamic.RouterTLSConfig{
				Options: route.NodeId,
				Domains: proxy.Domains(),
			},
		})
		routes.Route(ctx, fmt.Sprintf("%s#h21#insecure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier, PluginForwarder, PluginHath, PluginReplace},
			Service:     h21s,
			Rule:        h21r,
			Priority:    1500,
		})
		h11s := routes.IfAbsent(fmt.Sprintf("%s#h11", route.NodeId), func(n string) *dynamic.Service {
			return that.remoteService(ctx, n, routes, "https", route.NodeId, that.supervise(route), route.URC())
		})
		h11r := fmt.Sprintf("PathRegexp(`/mesh-h12/*`) && HeaderRegexp(`mesh-urn`, `%s`)", mtypes.URNMatcher("([-a-zA-Z\\d]+)", route.ID(ctx).SEQ, mtypes.CN, route.NodeId, route.InstId))
		routes.Route(ctx, fmt.Sprintf("%s#h11#secure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier, PluginForwarder},
			Service:     h11s,
			Rule:        h11r,
			Priority:    2000,
			TLS: &dynamic.RouterTLSConfig{
				Options: route.NodeId,
				Domains: proxy.Domains(),
			},
		})
		routes.Route(ctx, fmt.Sprintf("%s#h11#insecure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier, PluginForwarder},
			Service:     h11s,
			Rule:        h11r,
			Priority:    2000,
		})
	}
	h21r := fmt.Sprintf("PathRegexp(`/mesh-h12/*`) && HeaderRegexp(`mesh-urn`, `%s`)", mtypes.URNMatcher("([-a-zA-Z\\d]+)", that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId))
	h21s := that.bridgeService(ctx, routes)
	routes.Route(ctx, fmt.Sprintf("%s#h12#secure", mtypes.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     h21s,
		Rule:        h21r,
		Priority:    2000,
		TLS: &dynamic.RouterTLSConfig{
			Options: mtypes.LocalNodeId,
			Domains: proxy.Domains(),
		},
	})
	routes.Route(ctx, fmt.Sprintf("%s#h12#insecure", mtypes.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     h21s,
		Rule:        h21r,
		Priority:    2000,
	})
}

func (that *SnapShot) mdcService(ctx context.Context, routes *Routes, secure bool) string {
	return routes.IfAbsent(fmt.Sprintf("%s#mdc#%s", mtypes.LocalNodeId, tool.Ternary(secure, "secure", "insecure")), func(n string) *dynamic.Service {
		var servers []dynamic.Server
		for _, addr := range tool.MDC.Get() {
			servers = append(servers, dynamic.Server{URL: fmt.Sprintf("%s://%s", tool.Ternary(secure, "https", "h2c"), addr)})
		}
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: mtypes.LocalNodeId, Servers: servers}}
	})
}

func (that *SnapShot) withMDCRoute(ctx context.Context, routes *Routes) {
	if len(tool.MDC.Get()) < 1 {
		return
	}
	cr := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", mtypes.URNMatcher("([-a-zA-Z\\d]+)", that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId))
	routes.Route(ctx, fmt.Sprintf("%s#mdc#secure", mtypes.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     that.mdcService(ctx, routes, true),
		Rule:        cr,
		Priority:    0,
		TLS: &dynamic.RouterTLSConfig{
			Options: mtypes.LocalNodeId,
			Domains: proxy.Domains(),
		},
	})
	routes.Route(ctx, fmt.Sprintf("%s#mdc#insecure", mtypes.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     that.mdcService(ctx, routes, false),
		Rule:        cr,
		Priority:    0,
	})
}

func (that *SnapShot) proxyService(ctx context.Context, routes *Routes, secure bool) string {
	return routes.IfAbsent(fmt.Sprintf("%s#sx#%s", mtypes.LocalNodeId, tool.Ternary(secure, "secure", "insecure")), func(n string) *dynamic.Service {
		var servers []dynamic.Server
		for _, registration := range that.servers {
			v := tool.MapBy(registration.Attachments, fmt.Sprintf("%s.version", tool.Name.Get()))
			addr := fmt.Sprintf("%s://%s", tool.Ternary(secure, "https", "h2c"), registration.Address)
			servers = append(servers, dynamic.Server{URL: mtypes.ParseURL(ctx, addr).SetV(v).String()})
		}
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: mtypes.LocalNodeId, Servers: servers}}
	})
}

func (that *SnapShot) withProxyRoute(ctx context.Context, routes *Routes) {
	if len(that.servers) < 1 {
		return
	}
	rr := fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", mtypes.URNMatcher("([-a-zA-Z\\d]+)", that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId))
	routes.Route(ctx, fmt.Sprintf("%s#sx#secure", mtypes.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     that.proxyService(ctx, routes, true),
		Rule:        rr,
		Priority:    500,
		TLS: &dynamic.RouterTLSConfig{
			Options: mtypes.LocalNodeId,
			Domains: proxy.Domains(),
		},
	})
	routes.Route(ctx, fmt.Sprintf("%s#sx#insecure", mtypes.LocalNodeId), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier, PluginForwarder},
		Service:     that.proxyService(ctx, routes, false),
		Rule:        rr,
		Priority:    100,
	})
}

func (that *SnapShot) complexService(ctx context.Context, routes *Routes, group string, services []*mtypes.Service) string {
	return routes.IfAbsent(fmt.Sprintf("%s#%s", mtypes.LocalNodeId, group), func(n string) *dynamic.Service {
		var servers []dynamic.Server
		if plugin.PROXY.Match() {
			for _, addr := range tool.Address.Get().All() {
				servers = append(servers, dynamic.Server{URL: fmt.Sprintf("https://%s", addr)})
			}
			return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: mtypes.LocalNodeId, Servers: servers}}
		}
		for _, service := range services {
			if proxy.RecursionBreak(ctx, mtypes.URC(service.Address)) {
				continue
			}
			proto := tool.Ternary(service.Kind == mtypes.Restful, "http", "h2c")
			uri := mtypes.ParseURL(ctx, fmt.Sprintf("%s://%s", proto, service.Address)).SetS(prsim.MeshSubset.Get(service.Attrs)).URL().String()
			servers = append(servers, dynamic.Server{URL: uri})
		}
		return &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{ServersTransport: mtypes.LocalNodeId, Servers: servers}}
	})
}

func (that *SnapShot) withComplexRoute(ctx context.Context, routes *Routes) {
	if len(that.complex) < 1 {
		return
	}
	groups := map[string][]*mtypes.Service{}
	for _, registration := range that.complex {
		for _, service := range registration.InferService() {
			service.Attrs = tool.Merge(registration.Attachments, service.Attrs)
			group := fmt.Sprintf("%s#%s", formatLegacyName(registration.Name), strings.ToLower(service.Kind))
			groups[group] = append(groups[group], service)
		}
	}
	for group, services := range groups {
		var matchers []string
		var middlewares = []string{PluginBarrier, PluginForwarder}
		for _, service := range services {
			suffix := tool.Ternary(strings.HasSuffix(service.URN, "/*"), "", "/*")
			matchers = append(matchers, fmt.Sprintf("PathRegexp(`%s`)", strings.ReplaceAll(fmt.Sprintf("/%s%s", service.URN, suffix), "//", "/")))
			middlewares = tool.Distinct(tool.WithoutZero(append(middlewares, strings.Split(tool.MapBy(service.Attrs, "plugins"), ",")...)))
		}
		cr := strings.Join(matchers, " || ")
		cs := that.complexService(ctx, routes, group, services)
		routes.Route(ctx, fmt.Sprintf("%s#secure", cs), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY, TransportA, TransportB},
			Middlewares: middlewares,
			Service:     cs,
			Rule:        cr,
			Priority:    200,
			TLS: &dynamic.RouterTLSConfig{
				Options: mtypes.LocalNodeId,
				Domains: proxy.Domains(),
			},
		})
		routes.Route(ctx, fmt.Sprintf("%s#insecure", cs), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY, TransportA, TransportB},
			Middlewares: middlewares,
			Service:     cs,
			Rule:        cr,
			Priority:    200,
		})
	}
}

func (that *SnapShot) routesMessage(ctx context.Context) *dynamic.Message {
	routes := &Routes{services: map[string]*dynamic.Service{}, routers: map[string]*dynamic.Router{}}
	assemblies.WithAsmRoute(ctx, routes)
	tunnel.WithFOPSRoute(ctx, that.routers, routes, that)
	that.withRemoteRoute(ctx, routes)
	that.withClusterRoute(ctx, routes)
	that.withProxyRoute(ctx, routes)
	that.withMDCRoute(ctx, routes)
	that.withComplexRoute(ctx, routes)
	that.withBridgeRoute(ctx, routes)
	return &dynamic.Message{
		ProviderName: ProviderName,
		Configuration: &dynamic.Configuration{
			TLS: that.withTLS(ctx),
			HTTP: &dynamic.HTTPConfiguration{
				Routers:  routes.routers,
				Services: routes.services,
				Middlewares: map[string]*dynamic.Middleware{
					PluginRetry: {
						Retry: &dynamic.Retry{
							Attempts:        3,
							InitialInterval: types.Duration(time.Millisecond * 100),
						},
					},
					PluginRewrite: {
						StripPrefix: &dynamic.StripPrefix{
							Prefixes: rewritePaths,
						},
					},
					PluginReplace: {
						ReplacePath: &dynamic.ReplacePath{
							Path: "/mesh-h21/v1",
						},
					},
					PluginErrors:    {},
					PluginHeader:    {},
					PluginBarrier:   {},
					PluginForwarder: {},
					PluginSubfilter: {},
					PluginWedge:     {},
					PluginHath:      {},
				},
				Models:            map[string]*dynamic.Model{},
				ServersTransports: that.withTransports(ctx),
			},
		},
	}
}

func (that *SnapShot) withTransports(ctx context.Context) map[string]*dynamic.ServersTransport {
	defaultCts := &mtypes.RouteCertificate{
		HostRoot:  that.env.RootCrt,
		HostKey:   that.env.RootKey,
		HostCrt:   that.env.RootCrt,
		GuestRoot: that.env.RootCrt,
		GuestKey:  that.env.RootKey,
		GuestCrt:  that.env.RootCrt,
	}
	transports := map[string]*dynamic.ServersTransport{}
	routes := append(that.routers, &mtypes.Route{NodeId: "default"}, &mtypes.Route{NodeId: mtypes.LocalNodeId})
	for _, route := range routes {
		pxy := tool.Anyone(mtypes.ParseURL(ctx, tool.Anyone(that.supervise(route), route.URC()).String()).GetP(), route.Proxy)
		certification := route.GetCertificate(ctx).Override(defaultCts)
		var roots []tls.FileOrContent
		if "" != certification.HostRoot {
			roots = append(roots, tls.FileOrContent(certification.HostRoot))
		}
		if "" != certification.GuestRoot {
			roots = append(roots, tls.FileOrContent(certification.GuestRoot))
		}
		transports[route.NodeId] = &dynamic.ServersTransport{
			ServerName:         fmt.Sprintf("%s.%s.%s", that.env.NodeId, route.NodeId, mtypes.MeshDomain),
			InsecureSkipVerify: true,
			RootCAs:            roots,
			Certificates: tls.Certificates{
				{
					CertFile: tls.FileOrContent(certification.GuestCrt),
					KeyFile:  tls.FileOrContent(certification.GuestKey),
				},
			},
			MaxIdleConnsPerHost: 1,
			ForwardingTimeouts: &dynamic.ForwardingTimeouts{
				DialTimeout:     types.Duration(time.Second * 24),
				PingTimeout:     types.Duration(time.Second * 12),
				IdleConnTimeout: types.Duration(time.Minute * 1),
				ReadIdleTimeout: types.Duration(time.Minute * 2),
			},
			Proxy: tool.Ternary(len(that.proxies) > 0 || !tool.Proxy.Get().Empty(), "", pxy),
		}
	}
	return transports
}

func (that *SnapShot) withTLS(ctx context.Context) *dynamic.TLSConfiguration {
	defaultCts := &mtypes.RouteCertificate{
		HostRoot:  that.env.RootCrt,
		HostKey:   that.env.RootKey,
		HostCrt:   that.env.RootCrt,
		GuestRoot: that.env.RootCrt,
		GuestKey:  that.env.RootKey,
		GuestCrt:  that.env.RootCrt,
	}
	stores := map[string]tls.Store{
		"default": {
			DefaultCertificate: &tls.Certificate{
				CertFile: tls.FileOrContent(defaultCts.HostCrt),
				KeyFile:  tls.FileOrContent(defaultCts.HostKey),
			},
		},
		mtypes.LocalNodeId: {
			DefaultCertificate: &tls.Certificate{
				CertFile: tls.FileOrContent(defaultCts.HostCrt),
				KeyFile:  tls.FileOrContent(defaultCts.HostKey),
			},
		},
	}
	options := map[string]tls.Options{"default": *tls.DefaultTLSOptions.DeepCopy(), mtypes.LocalNodeId: *tls.DefaultTLSOptions.DeepCopy()}
	for _, route := range that.routers {
		certification := route.GetCertificate(ctx).Override(defaultCts)
		stores[route.NodeId] = tls.Store{
			DefaultCertificate: &tls.Certificate{
				CertFile: tls.FileOrContent(certification.HostCrt),
				KeyFile:  tls.FileOrContent(certification.HostKey),
			},
		}
		options[route.NodeId] = *tls.DefaultTLSOptions.DeepCopy()
	}
	return &dynamic.TLSConfiguration{Options: options, Stores: stores}
}
