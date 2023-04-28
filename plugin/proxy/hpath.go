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
	"github.com/traefik/traefik/v3/pkg/server/middleware"
	"net/http"
)

func init() {
	var _ http.Handler = new(hath)
	middleware.Provide(paths)
}

var paths = &hathMiddleware{}

type hath struct {
	name string
	next http.Handler
}

func (that *hath) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	prsim.MeshPath.SetHeader(request.Header, request.RequestURI)
	that.next.ServeHTTP(writer, request)
}

type hathMiddleware struct {
}

func (that *hathMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginHath, ProviderName)
}

func (that *hathMiddleware) Priority() int {
	return 0
}

func (that *hathMiddleware) Scope() int {
	return 1
}

func (that *hathMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &hath{next: next, name: name}, nil
}
