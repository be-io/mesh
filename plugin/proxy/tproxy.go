/*
 * Copyright (c) 2000, 2099, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/server/service"
	"net/http"
	"net/url"
)

func init() {
	var _ prsim.Listener = tProxies
	macro.Provide(prsim.IListener, tProxies)
	service.ProvideProxy(tProxies)
}

var tProxies = new(transportProxies)

type tproxy struct {
	endpoint string
}

func (that *tproxy) Proxy(req *http.Request) (*url.URL, error) {
	uri, err := url.Parse(that.endpoint)
	if nil != err {
		return nil, err
	}
	return http.ProxyURL(uri)(req)
}

type transportProxies struct {
}

func (that *transportProxies) Name() string {
	return "mesh"
}

func (that *transportProxies) New(endpoint string) service.Proxy {
	return &tproxy{endpoint: endpoint}
}

func (that *transportProxies) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.tproxy"}
}

func (that *transportProxies) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.NetworkRouteRefresh}
}

func (that *transportProxies) Listen(ctx context.Context, event *types.Event) error {
	var routes []*types.Route
	if err := event.TryGetObject(&routes); nil != err {
		return cause.Error(err)
	}
	return nil
}
