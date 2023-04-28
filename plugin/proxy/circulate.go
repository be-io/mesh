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
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	mtypes "github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v3/pkg/server/middleware"
	"math"
	"net/http"
)

func init() {
	var _ http.Handler = new(circulate)
	middleware.Provide(circulates)
}

const MDC = "mesh-mdc"

var circulates = &circulateMiddleware{MDC: tool.NewTraceId()}

type circulate struct {
	name string
	next http.Handler
}

func (that *circulate) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if macro.PermitCirculate.Enable() {
		that.next.ServeHTTP(writer, request)
		return
	}
	circulateCount := 0
	flags := request.Header.Values(MDC)
	for _, flag := range flags {
		if flag == circulates.MDC {
			circulateCount++
		}
	}
	// Max circulate times 5
	if circulateCount < 6 {
		if len(flags) < 31 {
			request.Header.Add(MDC, circulates.MDC)
		}
		that.next.ServeHTTP(writer, request)
		return
	}
	if request.ProtoMajor < 2 || nil == proxy.TCPRouters || nil == proxy.TCPRouters[TransportY] {
		writer.WriteHeader(http.StatusServiceUnavailable)
		if _, err := writer.Write([]byte("CirculateBreak")); nil != err {
			log.Error0(err.Error())
		}
		return
	}
	newRequest := request.Clone(request.Context())
	urn := prsim.MeshUrn.GetHeader(newRequest.Header)
	uname := mtypes.FromURN(mpc.Context(), urn)
	uname.NodeId = mtypes.LocalNodeId
	uname.Name = fmt.Sprintf("%s.%d", "mesh.builtin.fallback", 503)
	newRequest.Host = uname.String()
	prsim.MeshUrn.SetHeader(newRequest.Header, newRequest.Host)
	proxy.TCPRouters[TransportY].GetHTTPHandler().ServeHTTP(writer, newRequest)
}

type circulateMiddleware struct {
	MDC string
}

func (that *circulateMiddleware) Name() string {
	return fmt.Sprintf("circulate@%s", ProviderName)
}

func (that *circulateMiddleware) Priority() int {
	return math.MinInt
}

func (that *circulateMiddleware) Scope() int {
	return 0
}

func (that *circulateMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &circulate{next: next, name: name}, nil
}
