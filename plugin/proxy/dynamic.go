/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	mtypes "github.com/opendatav/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/provider"
	"github.com/traefik/traefik/v3/pkg/safe"
	"github.com/traefik/traefik/v3/pkg/tls"
	"github.com/traefik/traefik/v3/pkg/types"
)

func init() {
	var _ provider.Provider = meshGraph
	var _ prsim.Listener = new(meshNetGraph)
	macro.Provide(prsim.IListener, meshGraph)
}

var (
	meshGraph = &meshNetGraph{groups: map[string]*mtypes.Route{}, license: &mtypes.License{}}
)

type meshNetGraph struct {
	configurationChan chan<- dynamic.Message
	routes            []*mtypes.Route
	groups            map[string]*mtypes.Route
	registrations     mtypes.MetadataRegistrations
	license           *mtypes.License
}

func (that *meshNetGraph) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.dynamic"}
}

func (that *meshNetGraph) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.NetworkRouteRefresh, prsim.RegistryEventRefresh, prsim.LicenseImports, prsim.RoutePeriodRefresh}
}

func (that *meshNetGraph) Listen(ctx context.Context, event *mtypes.Event) error {
	if event.Binding.Match(prsim.NetworkRouteRefresh) {
		var routes []*mtypes.Route
		if err := event.TryGetObject(&routes); nil != err {
			return cause.Error(err)
		}
		groups := map[string]*mtypes.Route{}
		for _, route := range routes {
			groups[route.NodeId] = route
			groups[route.InstId] = route
		}
		that.routes = routes
		that.groups = groups
	}
	if event.Binding.Match(prsim.LicenseImports) {
		var license *mtypes.License
		if err := event.TryGetObject(&license); nil != err {
			return cause.Error(err)
		}
		that.license = license
	}
	if event.Binding.Match(prsim.RegistryEventRefresh) {
		var registrations mtypes.MetadataRegistrations
		if err := event.TryGetObject(&registrations); nil != err {
			return cause.Error(err)
		}
		that.registrations = registrations
	}
	if event.Binding.Match(prsim.RoutePeriodRefresh) {
		return that.refresh(ctx)
	}
	return nil
}

func (that *meshNetGraph) refresh(ctx context.Context) error {
	log.Debug(ctx, "Refresh mesh routes doing. ")
	if nil != that.configurationChan {
		environ, err := aware.Network.GetEnviron(ctx)
		if nil != err {
			return cause.Error(err)
		}
		snap := &SnapShot{
			env:     environ,
			routers: that.routes,
			groups:  tool.Anyone(that.groups, map[string]*mtypes.Route{}),
			proxies: that.registrations.Of(mtypes.PROXY),
			servers: that.registrations.Of(mtypes.SERVER),
			sets:    that.registrations.Of(mtypes.METADATA),
			complex: that.registrations.Of(mtypes.COMPLEX),
			lic:     tool.Anyone(that.license, new(mtypes.License)),
		}
		message := snap.routesMessage(ctx)
		that.stringify(ctx, message)
		that.configurationChan <- *message
	} else {
		log.Warn(ctx, "Mesh proxy dynamic routers dont active now. ")
	}
	log.Debug(ctx, "Refresh mesh routes done. ")
	return nil
}

func (that *meshNetGraph) Init() error {
	return nil
}

func (that *meshNetGraph) Provide(configurationChan chan<- dynamic.Message, pool *safe.Pool) error {
	that.configurationChan = configurationChan
	return nil
}

func (that *meshNetGraph) stringify(ctx context.Context, message *dynamic.Message) {
	copyConf := message.Configuration.DeepCopy()
	if copyConf.TLS != nil {
		copyConf.TLS.Certificates = nil
		if copyConf.TLS.Options != nil {
			cleanedOptions := make(map[string]tls.Options, len(copyConf.TLS.Options))
			for name, option := range copyConf.TLS.Options {
				option.ClientAuth.CAFiles = []types.FileOrContent{}
				cleanedOptions[name] = option
			}
			copyConf.TLS.Options = cleanedOptions
		}
		for k := range copyConf.TLS.Stores {
			st := copyConf.TLS.Stores[k]
			st.DefaultCertificate = nil
			copyConf.TLS.Stores[k] = st
		}
	}
	if copyConf.HTTP != nil {
		for _, transport := range copyConf.HTTP.ServersTransports {
			transport.Certificates = tls.Certificates{}
			transport.RootCAs = []types.FileOrContent{}
		}
	}
	if jsonConf, err := aware.Codec.EncodeString(copyConf); nil != err {
		log.Error(ctx, "Could not marshal dynamic configuration: %v", err)
		log.Debug(ctx, "Configuration push: [struct] %#v", copyConf)
	} else {
		log.Debug(ctx, "Configuration push: %s", jsonConf)
	}
}
