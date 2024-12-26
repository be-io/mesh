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
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/system"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/server/middleware"
	"net/http"
	"strings"
	"time"
)

func init() {
	var _ http.Handler = new(barrier)
	var _ prsim.Listener = barriers
	middleware.Provide(barriers)
	macro.Provide(prsim.IListener, barriers)
}

var barriers = &barrierMiddleware{routes: map[string]*types.Route{}}

type barrier struct {
	name string
	next http.Handler
}

func (that *barrier) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	uname := tool.Anyone(prsim.MeshUrn.GetHeader(request.Header), request.Host)
	if !types.MatchURNDomain(uname) {
		that.next.ServeHTTP(writer, request)
		return
	}
	ctx := mpc.Context()
	nodeId := tool.Anyone(prsim.MeshFromNodeId.GetHeader(request.Header), prsim.MeshFromInstId.GetHeader(request.Header))
	urn := types.FromURN(ctx, uname)
	if "" == nodeId || "" == urn.NodeId {
		that.next.ServeHTTP(writer, request)
		return
	}
	if !barriers.Disabled(ctx, nodeId, urn) {
		that.next.ServeHTTP(writer, request)
		return
	}
	writer.WriteHeader(http.StatusForbidden)
	_, err := writer.Write([]byte(cause.NetDisable.Message))
	if nil != err {
		log.Error0("%s, %s", uname, err.Error())
	}
}

type barrierMiddleware struct {
	// CopyWrite
	routes map[string]*types.Route
}

func (that *barrierMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginBarrier, ProviderName)
}

func (that *barrierMiddleware) Priority() int {
	return 0
}

func (that *barrierMiddleware) Scope() int {
	return 1
}

func (that *barrierMiddleware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.barrier"}
}

func (that *barrierMiddleware) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.NetworkRouteRefresh}
}

func (that *barrierMiddleware) Listen(ctx context.Context, event *types.Event) error {
	var routes []*types.Route
	if err := event.TryGetObject(&routes); nil != err {
		return cause.Error(err)
	}
	emap := map[string]*types.Route{}
	for _, route := range routes {
		emap[route.NodeId] = route
		emap[route.InstId] = route
		emap[strings.ToLower(route.NodeId)] = route
		emap[strings.ToLower(route.InstId)] = route
		emap[strings.ToUpper(route.NodeId)] = route
		emap[strings.ToUpper(route.InstId)] = route
	}
	that.routes = emap
	return nil
}

func (that *barrierMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &barrier{next: next, name: name}, nil
}

func (that *barrierMiddleware) Disabled(ctx context.Context, nodeId string, urn *types.URN) bool {
	if strings.Index(urn.Name, tool.Name.Get()) == 0 || tool.IsLocalEnv(system.Environ.Get(), nodeId) {
		return false
	}
	routes := that.routes
	if nil == routes || (nil == routes[nodeId] && nil == routes[urn.NodeId]) {
		return true
	}
	return that.routeDisable(routes[nodeId], true) || that.routeDisable(routes[urn.NodeId], false)
}

func (that *barrierMiddleware) routeDisable(route *types.Route, disableIfAbsent bool) bool {
	if nil == route {
		return disableIfAbsent
	}
	if (route.Status & int32(types.Disabled)) == int32(types.Disabled) {
		return true
	}
	if (route.Status & int32(types.Removed)) == int32(types.Removed) {
		return true
	}
	if route.ExpireAt > 0 && route.ExpireAt < time.Now().UnixMilli() {
		return true
	}
	return false
}
