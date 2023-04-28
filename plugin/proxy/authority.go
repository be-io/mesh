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
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/system"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/server/middleware"
	"math"
	"net/http"
)

func init() {
	var _ http.Handler = new(authority)
	middleware.Provide(authorities)
}

var authorities = &authorityMiddleware{}

type authority struct {
	name string
	next http.Handler
}

func (that *authority) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodTrace {
		http.Error(writer, "Forbidden", http.StatusMethodNotAllowed)
		return
	}
	if !proxy.InsecureEnable && request.ProtoMajor < 2 && nil == request.TLS {
		http.Error(writer, "Forbidden", http.StatusBadRequest)
		return
	}
	//
	if "" == prsim.MeshTraceId.GetHeader(request.Header) {
		prsim.MeshTraceId.SetHeader(request.Header, tool.NewTraceId())
	}
	if "" == prsim.MeshFromNodeId.GetHeader(request.Header) {
		prsim.MeshFromNodeId.SetHeader(request.Header, system.Environ.Get().NodeId)
	}
	if "" == prsim.MeshFromInstId.GetHeader(request.Header) {
		prsim.MeshFromInstId.SetHeader(request.Header, system.Environ.Get().InstId)
	}
	if "" == prsim.MeshVersion.GetHeader(request.Header) {
		prsim.MeshVersion.SetHeader(request.Header, prsim.Version)
	}
	if "" == writer.Header().Get("X-Frame-Options") {
		writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
	}
	if "" == writer.Header().Get("X-XSS-Protection") {
		writer.Header().Set("X-XSS-Protection", "1")
	}
	if urn := tool.Anyone(prsim.MeshUrn.GetHeader(request.Header)); types.MatchURNDomain(urn) {
		request.Host = urn
		request.Header.Set("Host", urn)
		that.next.ServeHTTP(writer, request)
		return
	}
	if types.MatchURNDomain(request.Host) {
		prsim.MeshUrn.SetHeader(request.Header, request.Host)
		that.next.ServeHTTP(writer, request)
		return
	}
	if host := tool.Anyone(request.Header.Get("Host"), request.Header.Get("host")); types.MatchURNDomain(host) {
		prsim.MeshUrn.SetHeader(request.Header, host)
		that.next.ServeHTTP(writer, request)
		return
	}
	that.next.ServeHTTP(writer, request)
}

type authorityMiddleware struct {
}

func (that *authorityMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginAuthority, ProviderName)
}

func (that *authorityMiddleware) Priority() int {
	return math.MaxInt
}

func (that *authorityMiddleware) Scope() int {
	return 0
}

func (that *authorityMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &authority{next: next, name: name}, nil
}
