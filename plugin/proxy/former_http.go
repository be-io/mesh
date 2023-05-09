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
	"github.com/traefik/traefik/v2/pkg/server/middleware"
	"net/http"
)

func init() {
	var _ http.Handler = new(formerHttp)
	middleware.Provide(formerHttps)
}

var formerHttps = &formerHttpMiddleware{}

type formerHttp struct {
	name string
	next http.Handler
}

func (that *formerHttp) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	that.next.ServeHTTP(writer, request)
}

type formerHttpMiddleware struct {
}

func (that *formerHttpMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginTransformer, ProviderName)
}

func (that *formerHttpMiddleware) Priority() int {
	return 0
}

func (that *formerHttpMiddleware) Scope() int {
	return 0
}

func (that *formerHttpMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &formerHttp{next: next, name: name}, nil
}
