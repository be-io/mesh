/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/traefik/traefik/v3/pkg/server/middleware"
	"net/http"
	"strings"
)

func init() {
	var _ http.Handler = new(authority)
	middleware.Provide(headers)
}

var headers = &headerMiddleware{}

type header struct {
	name string
	next http.Handler
}

func (that *header) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if macro.PHeader.Enable() {
		buff := &bytes.Buffer{}
		for k, vs := range request.Header {
			buff.WriteString(k)
			buff.WriteRune('=')
			buff.WriteString(strings.Join(vs, ","))
			buff.WriteRune(';')
		}
		log.Info0("%s:%s", request.Host, buff.String())
	}
	that.next.ServeHTTP(writer, request)
}

type headerMiddleware struct {
}

func (that *headerMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginHeader, ProviderName)
}

func (that *headerMiddleware) Priority() int {
	return 0
}

func (that *headerMiddleware) Scope() int {
	return 0
}

func (that *headerMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &header{next: next, name: name}, nil
}
