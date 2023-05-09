/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v2/pkg/server/service/loadbalancer/wrr"
	"net/http"
)

func init() {
	wrr.Provide(new(SubsetLBStrategy))
}

type SubsetLBStrategy struct {
}

func (that *SubsetLBStrategy) Name() string {
	return "subset@mesh"
}

func (that *SubsetLBStrategy) Priority() int {
	return 1
}

func (that *SubsetLBStrategy) Next(w http.ResponseWriter, req *http.Request, servers []wrr.Server) []wrr.Server {
	subset := prsim.MeshSubset.GetHeader(req.Header)
	if "" == subset {
		return that.skipServerNoSubset(servers)
	}
	var ss []wrr.Server
	for _, server := range servers {
		if nil == server.URL() {
			continue
		}
		if types.FromURL(server.URL()).GetS() == subset {
			ss = append(ss, server)
		}
	}
	if len(ss) > 0 {
		return ss
	}
	return servers
}

func (that *SubsetLBStrategy) skipServerNoSubset(servers []wrr.Server) []wrr.Server {
	var ss []wrr.Server
	for _, server := range servers {
		if nil == server.URL() {
			continue
		}
		var s = types.FromURL(server.URL()).GetS()
		if len(s) == 0 {
			ss = append(ss, server)
		}
	}
	if len(ss) == 0 {
		return servers
	}
	return ss
}
