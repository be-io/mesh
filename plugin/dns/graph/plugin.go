/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package graph

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/be-io/mesh/plugin/dns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"google.golang.org/grpc/resolver"
	"math"
	"net"
	"strings"
)

func init() {
	var _ plugin.Plugin = Resolver
	var _ prsim.Listener = Resolver
	macro.Provide(prsim.IListener, Resolver)
	plugin.Provide(Resolver)
}

var Resolver = new(graphPlugin)
var Backoff = []resolver.Address{{Addr: tool.Address.Get().Any()}}

type graphPlugin struct {
	routes   map[string]*types.Route
	services map[string][]*types.Service
}

func (that *graphPlugin) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.dns.resolver"}
}

func (that *graphPlugin) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.NetworkRouteRefresh, prsim.RegistryEventRefresh}
}

func (that *graphPlugin) Listen(ctx context.Context, event *types.Event) error {
	if event.Binding.Match(prsim.NetworkRouteRefresh) {
		var edges []*types.Route
		if err := event.TryGetObject(&edges); nil != err {
			return cause.Error(err)
		}
		emap := map[string]*types.Route{}
		for _, edge := range edges {
			emap[edge.NodeId] = edge
			emap[edge.InstId] = edge
			emap[strings.ToLower(edge.NodeId)] = edge
			emap[strings.ToLower(edge.InstId)] = edge
			emap[strings.ToUpper(edge.NodeId)] = edge
			emap[strings.ToUpper(edge.InstId)] = edge
		}
		that.routes = emap
	}
	if event.Binding.Match(prsim.RegistryEventRefresh) {
		var registrations types.MetadataRegistrations
		if err := event.TryGetObject(&registrations); nil != err {
			return cause.Error(err)
		}
		emap := map[string][]*types.Service{}
		for _, service := range registrations.Of(types.METADATA).InferService() {
			urn := types.FromURN(ctx, service.URN)
			if nil == emap[urn.Name] {
				emap[urn.Name] = []*types.Service{}
			}
			emap[urn.Name] = append(emap[urn.Name], service)
		}
		that.services = emap
	}
	return nil
}

func (that *graphPlugin) Name() []string {
	return []string{"graph"}
}

func (that *graphPlugin) Priority() int {
	return math.MaxInt
}

func (that *graphPlugin) ServeDNS(ctx context.Context, pip plugin.Pip, r *request.Request) ([]dns.RR, error) {
	rr, err := that.resolve(ctx, r.QName(), r.RemoteAddr())
	if nil != err {
		return nil, cause.Error(err)
	}
	if len(rr) > 0 {
		return rr, nil
	}
	return pip.Pip(ctx, r)
}

func (that *graphPlugin) ResolveAddress(authority string, check bool) ([]resolver.Address, error) {
	ctx := mpc.Context()
	urn := types.FromURN(ctx, authority)
	if !types.MatchURNDomain(authority) || "" == urn.NodeId {
		ips, err := net.LookupIP(authority)
		if nil != err {
			log.Error(ctx, err.Error())
			return Backoff, nil
		}
		var addresses []resolver.Address
		for _, ip := range ips {
			log.Info(ctx, "%s -> %s", urn.String(), ip.String())
			addresses = append(addresses, resolver.Address{Addr: ip.String()})
			return addresses, nil
		}
	}
	adds, err := that.ResolveByServices(ctx, urn, check)
	if nil != err {
		return nil, cause.Error(err)
	}
	if len(adds) > 0 {
		return adds, nil
	}
	adds, err = that.ResolveByRoutes(ctx, urn)
	if nil != err {
		return nil, cause.Error(err)
	}
	if len(adds) > 0 {
		return adds, nil
	}
	return Backoff, nil
}

func (that *graphPlugin) ResolveByRoutes(ctx context.Context, urn *types.URN) ([]resolver.Address, error) {
	var addresses []resolver.Address
	for _, route := range that.routes {
		if urn.MatchNode(ctx, route.NodeId) {
			log.Info(ctx, "%s -> %s", urn.String(), route.URC().String())
			for _, addr := range route.URC().Addrs() {
				addresses = append(addresses, resolver.Address{Addr: addr})
			}
		}
	}
	if len(addresses) < 1 {
		log.Error(ctx, "No service named %s", urn.String())
	}
	return addresses, nil
}

func (that *graphPlugin) ResolveByServices(ctx context.Context, urn *types.URN, check bool) ([]resolver.Address, error) {
	var addresses []resolver.Address
	services := that.services[urn.Name]
	if nil == services {
		log.Error(ctx, "No service named %s", urn.String())
		return addresses, nil
	}
	for _, service := range services {
		uname := types.FromURN(ctx, service.URN)
		if uname.Match(ctx, uname) {
			log.Info(ctx, "%s -> %s", urn.String(), service.Address)
			addresses = append(addresses, resolver.Address{Addr: service.Address})
		}
	}
	if len(addresses) < 1 {
		log.Error(ctx, "No service named %s", urn.String())
	}
	return addresses, nil
}
