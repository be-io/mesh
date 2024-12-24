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
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/tls"
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
	master := routes.IfAbsent(fmt.Sprintf("%s#master", name), func(n string) *dynamic.Service {
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
	slave := routes.IfAbsent(fmt.Sprintf("%s#slave", name), func(n string) *dynamic.Service {
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
		rs := routes.IfAbsent(route.NodeId, func(n string) *dynamic.Service {
			return that.remoteService(ctx, n, routes, "https", route.NodeId, that.supervise(route), route.URC())
		})
		rr := fmt.Sprintf("HeaderRegexp(`x-ptp-target-node-id`, `%s`) || HeaderRegexp(`x-ptp-target-inst-id`, `%s`) || HeaderRegexp(`mesh-urn`, `%s`)", route.NodeId, route.InstId, mtypes.URNMatcher("([-a-zA-Z\\d]+)", route.ID(ctx).SEQ, mtypes.CN, route.NodeId, route.InstId))
		routes.Route(ctx, fmt.Sprintf("%s#secure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier},
			Service:     rs,
			Rule:        rr,
			Priority:    3500,
			TLS: &dynamic.RouterTLSConfig{
				Options: route.NodeId,
				Domains: proxy.Domains(),
			},
		})
		routes.Route(ctx, fmt.Sprintf("%s#insecure", route.NodeId), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY},
			Middlewares: []string{PluginBarrier},
			Service:     rs,
			Rule:        rr,
			Priority:    3100,
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
	cr := pattern()
	routes.Route(ctx, fmt.Sprintf("%s#secure", cs), &dynamic.Router{
		EntryPoints: []string{TransportX, TransportY},
		Middlewares: []string{PluginBarrier},
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
		Middlewares: []string{PluginBarrier},
		Service:     cs,
		Rule:        cr,
		Priority:    priority,
	})
}

func (that *SnapShot) withClusterRoute(ctx context.Context, routes *Routes) {
	nvs := map[string]mtypes.MetadataRegistrations{}
	for _, registration := range that.sets {
		nvs[registration.Name] = append(nvs[registration.Name], registration)
	}
	for name, instances := range nvs {
		if tool.Name.Get() == name {
			that.withClusterSetRoute(ctx, routes, name, "fh2", instances, 4500, func(registration *mtypes.MetadataRegistration) string {
				return mtypes.ParseURL(ctx, registration.Address).SetS(prsim.MeshSubset.Get(registration.Attachments)).SetSchema("h2c").URL().String()
			}, func() string {
				return "PathRegexp(`/org.ppc.ptp.PrivateTransferTransport/.*`)"
			})
			that.withClusterSetRoute(ctx, routes, name, "fh1", instances, 4100, func(registration *mtypes.MetadataRegistration) string {
				h, p := assemblies.ParseHost(ctx, registration.Address)
				uri := mtypes.ParseURL(ctx, registration.Address).SetS(prsim.MeshSubset.Get(registration.Attachments)).SetHost(fmt.Sprintf("%s:%d", h, p)).SetSchema("http").URL().String()
				return uri
			}, func() string {
				return "PathRegexp(`/v1/interconn/chan/(pop|push|peek|release)`)"
			})
		}
		that.withClusterSetRoute(ctx, routes, name, "h2", instances, 1000, func(registration *mtypes.MetadataRegistration) string {
			return mtypes.ParseURL(ctx, registration.Address).SetS(prsim.MeshSubset.Get(registration.Attachments)).SetSchema("h2c").URL().String()
		}, func() string {
			urn := mtypes.URNMatcher(name, that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId)
			if tool.Name.Get() != name {
				return fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", urn)
			}
			return fmt.Sprintf("PathRegexp(`/org.ppc.ptp\\..*`) || HeaderRegexp(`mesh-urn`, `%s`)", urn)
		})
		that.withClusterSetRoute(ctx, routes, name, "h1", instances, 1200, func(registration *mtypes.MetadataRegistration) string {
			h, p := assemblies.ParseHost(ctx, registration.Address)
			uri := mtypes.ParseURL(ctx, registration.Address).SetS(prsim.MeshSubset.Get(registration.Attachments)).SetHost(fmt.Sprintf("%s:%d", h, p)).SetSchema("http").URL().String()
			return uri
		}, func() string {
			urn := mtypes.URNMatcher(fmt.Sprintf("h1\\.%s", name), that.env.ID(ctx).SEQ, mtypes.CN, that.env.NodeId, that.env.InstId, mtypes.LocalNodeId, mtypes.LocalInstId)
			if tool.Name.Get() != name {
				return fmt.Sprintf("HeaderRegexp(`mesh-urn`, `%s`)", urn)
			}
			return fmt.Sprintf("PathRegexp(`/v1/interconn/.*`) || HeaderRegexp(`mesh-urn`, `%s`)", urn)
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
	priorities := make(map[string]int)
	for _, registration := range that.complex {
		for _, service := range registration.InferService() {
			service.Attrs = tool.Merge(registration.Attachments, service.Attrs)
			group := fmt.Sprintf("%s#%s", registration.Name, strings.ToLower(service.Kind))
			groups[group] = append(groups[group], service)
			if len(service.Attrs) > 0 && service.Attrs["priority"] == "high" {
				priorities[group] = 100
			}
		}
	}
	for group, services := range groups {
		var matchers []string
		var middlewares = []string{PluginBarrier}
		for _, service := range services {
			urn := tool.Ternary(strings.HasSuffix(service.URN, "/*"), fmt.Sprintf("%s/.*", strings.TrimSuffix(service.URN, "/*")), service.URN)
			matchers = append(matchers, fmt.Sprintf("PathRegexp(`%s`)", strings.ReplaceAll(fmt.Sprintf("/%s", urn), "//", "/")))
			middlewares = tool.Distinct(tool.WithoutZero(append(middlewares, strings.Split(tool.MapBy(service.Attrs, "plugins"), ",")...)))
		}
		cr := strings.Join(matchers, " || ")
		cs := that.complexService(ctx, routes, group, services)
		routes.Route(ctx, fmt.Sprintf("%s#secure", cs), &dynamic.Router{
			EntryPoints: []string{TransportX, TransportY, TransportA, TransportB},
			Middlewares: middlewares,
			Service:     cs,
			Rule:        cr,
			Priority:    4200 + priorities[group],
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
			Priority:    4200 + priorities[group],
		})
	}
}

func (that *SnapShot) dyn(ctx context.Context, route *mtypes.Route, middlewares ...string) []string {
	var ms []string
	for _, m := range middlewares {
		ms = append(ms, m)
	}
	return ms
}

func (that *SnapShot) routesMessage(ctx context.Context) *dynamic.Message {
	routes := &Routes{services: map[string]*dynamic.Service{}, routers: map[string]*dynamic.Router{}}
	assemblies.WithAsmRoute(ctx, routes)
	that.withRemoteRoute(ctx, routes)
	that.withClusterRoute(ctx, routes)
	that.withComplexRoute(ctx, routes)
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
					PluginErrors:  {},
					PluginHeader:  {},
					PluginBarrier: {},
					PluginHath:    {},
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
			InsecureSkipVerify: proxy.InsecureSkip,
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
