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
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/system"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
	"math"
	"net/http"
	"strings"
)

func init() {
	var _ http.Handler = new(incStack)
	middleware.Provide(incStacks)
}

var (
	incStacks = &incStackMiddleware{}
	incURL    = map[string]string{
		"v1/interconn/chan/invoke":    "mesh.chan.push",
		"v1/interconn/chan/transport": "mesh.chan.push",
		"v1/interconn/chan/push":      "mesh.chan.push",
		"v1/interconn/chan/pop":       "mesh.chan.pop",
		"v1/interconn/chan/peek":      "mesh.chan.peek",
		"v1/interconn/chan/release":   "mesh.chan.release",
	}
)

type incStack struct {
	name string
	next http.Handler
}

// ServeHTTP
func (that *incStack) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	tid := prsim.MeshTargetNodeId.GetHeader(request.Header)
	urn := prsim.MeshUrn.GetHeader(request.Header)
	if "" == tid {
		tid = types.FromURN(macro.Context(), urn).NodeId
	}
	for k, v := range incURL {
		if strings.Contains(request.URL.Path, k) {
			path := prsim.MeshPath.GetHeader(request.Header)
			prsim.MeshUrn.SetHeader(request.Header, types.AnyURN(macro.Context(), v, tid))
			if tool.IsLocalEnv(system.Environ.Get(), tid) {
				request.URL.Path = path
				if request.URL.RawPath != "" {
					request.URL.RawPath = path
				}
			}
			break
		}
		if strings.Contains(urn, v) {
			prsim.MeshTargetNodeId.SetHeader(request.Header, tid)
			prsim.MeshTargetInstId.SetHeader(request.Header, tid)
			if !tool.IsLocalEnv(system.Environ.Get(), tid) {
				request.URL.Path = k
				if request.URL.RawPath != "" {
					request.URL.RawPath = k
				}
			}
			break
		}
	}
	that.next.ServeHTTP(writer, request)
}

type incStackMiddleware struct {
}

func (that *incStackMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginINC, ProviderName)
}

func (that *incStackMiddleware) Priority() int {
	return math.MaxInt - 1
}

func (that *incStackMiddleware) Scope() int {
	return 0
}

func (that *incStackMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &incStack{next: next, name: name}, nil
}
