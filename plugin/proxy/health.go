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
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
	"net/http"
	"strings"
)

func init() {
	var _ http.Handler = new(health)
	middleware.Provide(healths)
}

var healths = &healthMiddleware{}

type health struct {
	name string
	next http.Handler
}

func (that *health) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.URL.Path, "/f5/checkHealth") {
		if _, err := writer.Write([]byte("@the@thealth@is@good@")); nil != err {
			log.Error0("Health response with %s", err.Error())
		}
		return
	}
	if request.URL.Path == "/stats" {
		if _, err := writer.Write([]byte(cause.Success.Code)); nil != err {
			log.Error0("Health response with %s", err.Error())
		}
		return
	}
	that.next.ServeHTTP(writer, request)
}

type healthMiddleware struct {
}

func (that *healthMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginHealth, ProviderName)
}

func (that *healthMiddleware) Priority() int {
	return 0
}

func (that *healthMiddleware) Scope() int {
	return 0
}

func (that *healthMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &health{next: next, name: name}, nil
}
