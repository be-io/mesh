/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/system"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/server/service/loadbalancer/wrr"
	"net/http"
	"strings"
)

func init() {
	wrr.Provide(new(StatefulLBStrategy))
}

type StatefulLBStrategy struct {
}

func (that *StatefulLBStrategy) Name() string {
	return "stateful@mesh"
}

func (that *StatefulLBStrategy) Priority() int {
	return 0
}

func (that *StatefulLBStrategy) Next(w http.ResponseWriter, req *http.Request, servers []wrr.Server) []wrr.Server {
	uname := tool.Anyone(prsim.MeshUrn.GetHeader(req.Header), req.Host)
	if !types.MatchURNDomain(uname) {
		return servers
	}
	urn := types.FromURN(macro.Context(), uname)
	isMyNode := tool.IsLocalEnv(system.Environ.Get(), urn.NodeId)
	if isMyNode {
		prsim.MeshIncomingProxy.SetHeader(req.Header, tool.Runtime.Get().String())
	} else {
		prsim.MeshOutgoingProxy.SetHeader(req.Header, tool.Runtime.Get().String())
		prsim.MeshIncomingProxy.SetHeader(req.Header, tool.Runtime.Get().String())
	}
	if !proxy.StatefulEnable {
		return servers
	}
	if "" == urn.Name {
		return servers
	}
	if !isMyNode {
		return servers
	}
	pair := strings.Split(prsim.MeshOutgoingHost.GetHeader(req.Header), "@")
	if len(pair) < 2 || strings.Index(urn.Name, pair[0]) != 0 {
		return servers
	}
	for _, server := range servers {
		if nil == server.URL() {
			continue
		}
		if pair[1] == server.URL().Host {
			return []wrr.Server{server}
		}
	}
	return servers
}
