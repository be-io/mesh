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
	"encoding/json"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"io"
	"strings"
)

func init() {
	var _ prsim.EndpointSticker[*bytes.Buffer, json.RawMessage] = new(proxyDashboard)
	macro.Provide(prsim.IEndpointSticker, new(proxyDashboard))
	var _ prsim.EndpointSticker[*bytes.Buffer, json.RawMessage] = new(proxyConfiguration)
	macro.Provide(prsim.IEndpointSticker, new(proxyConfiguration))

	var _ prsim.EndpointSticker[*types.Paging, json.RawMessage] = new(proxyEndpoints)
	macro.Provide(prsim.IEndpointSticker, new(proxyEndpoints))
	var _ prsim.EndpointSticker[*endpointDescriber, json.RawMessage] = new(proxyEndpointDescribe)
	macro.Provide(prsim.IEndpointSticker, new(proxyEndpointDescribe))

	var _ prsim.EndpointSticker[*types.Paging, json.RawMessage] = new(proxyRoutes)
	macro.Provide(prsim.IEndpointSticker, new(proxyRoutes))
	var _ prsim.EndpointSticker[*endpointDescriber, json.RawMessage] = new(proxyRouteDescribe)
	macro.Provide(prsim.IEndpointSticker, new(proxyRouteDescribe))

	var _ prsim.EndpointSticker[*types.Paging, json.RawMessage] = new(proxyServices)
	macro.Provide(prsim.IEndpointSticker, new(proxyServices))
	var _ prsim.EndpointSticker[*endpointDescriber, json.RawMessage] = new(proxyServiceDescribe)
	macro.Provide(prsim.IEndpointSticker, new(proxyServiceDescribe))

	var _ prsim.EndpointSticker[*types.Paging, json.RawMessage] = new(proxyMiddlewares)
	macro.Provide(prsim.IEndpointSticker, new(proxyMiddlewares))
	var _ prsim.EndpointSticker[*endpointDescriber, json.RawMessage] = new(proxyMiddlewareDescribe)
	macro.Provide(prsim.IEndpointSticker, new(proxyMiddlewareDescribe))
}

func endpoint(path string) (*bytes.Buffer, error) {
	resp, err := tool.Client.Get(fmt.Sprintf("http://127.0.0.1:%s/%s", strings.Split(proxy.TransportC, ":")[1], path))
	if nil != err {
		return nil, cause.Error(err)
	}
	defer func() { log.Catch(resp.Body.Close()) }()
	var buff bytes.Buffer
	if _, err = io.Copy(&buff, resp.Body); nil != err {
		return nil, cause.Error(err)
	}
	return &buff, nil
}

type endpointDescriber struct {
	Kind  string `json:"kind" yaml:"kind" xml:"kind"`
	Index int    `json:"index" yaml:"index" xml:"index"`
}

type proxyDashboard struct {
}

func (that *proxyDashboard) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.dashboard"}
}

func (that *proxyDashboard) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.dashboard"}
}

func (that *proxyDashboard) I() *bytes.Buffer {
	return new(bytes.Buffer)
}

func (that *proxyDashboard) O() json.RawMessage {
	return []byte{}
}

func (that *proxyDashboard) Stick(ctx context.Context, varg *bytes.Buffer) (json.RawMessage, error) {
	buff, err := endpoint("api/overview")
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyConfiguration struct {
}

func (that *proxyConfiguration) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.configuration"}
}

func (that *proxyConfiguration) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.configuration"}
}

func (that *proxyConfiguration) I() *bytes.Buffer {
	return new(bytes.Buffer)
}

func (that *proxyConfiguration) O() json.RawMessage {
	return []byte{}
}

func (that *proxyConfiguration) Stick(ctx context.Context, varg *bytes.Buffer) (json.RawMessage, error) {
	buff, err := endpoint("api/rawdata")
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyEndpoints struct {
}

func (that *proxyEndpoints) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.endpoints"}
}

func (that *proxyEndpoints) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.endpoints"}
}

func (that *proxyEndpoints) I() *types.Paging {
	return new(types.Paging)
}

func (that *proxyEndpoints) O() json.RawMessage {
	return []byte{}
}

func (that *proxyEndpoints) Stick(ctx context.Context, varg *types.Paging) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/entrypoints?pag=%d&per_page=%d", varg.Index, varg.Limit))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyEndpointDescribe struct {
}

func (that *proxyEndpointDescribe) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.endpoint"}
}

func (that *proxyEndpointDescribe) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.endpoint"}
}

func (that *proxyEndpointDescribe) I() *endpointDescriber {
	return new(endpointDescriber)
}

func (that *proxyEndpointDescribe) O() json.RawMessage {
	return []byte{}
}

func (that *proxyEndpointDescribe) Stick(ctx context.Context, varg *endpointDescriber) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/entrypoints/%d", varg.Index))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyRoutes struct {
}

func (that *proxyRoutes) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.routes"}
}

func (that *proxyRoutes) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.routes"}
}

func (that *proxyRoutes) I() *types.Paging {
	return new(types.Paging)
}

func (that *proxyRoutes) O() json.RawMessage {
	return []byte{}
}

func (that *proxyRoutes) Stick(ctx context.Context, varg *types.Paging) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/%s/routers?pag=%d&per_page=%d", strings.ToLower(varg.GetFactor("kind")), varg.Index, varg.Limit))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyRouteDescribe struct {
}

func (that *proxyRouteDescribe) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.route"}
}

func (that *proxyRouteDescribe) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.route"}
}

func (that *proxyRouteDescribe) I() *endpointDescriber {
	return new(endpointDescriber)
}

func (that *proxyRouteDescribe) O() json.RawMessage {
	return []byte{}
}

func (that *proxyRouteDescribe) Stick(ctx context.Context, varg *endpointDescriber) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/%s/routers/%d", strings.ToLower(varg.Kind), varg.Index))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyServices struct {
}

func (that *proxyServices) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.services"}
}

func (that *proxyServices) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.services"}
}

func (that *proxyServices) I() *types.Paging {
	return new(types.Paging)
}

func (that *proxyServices) O() json.RawMessage {
	return []byte{}
}

func (that *proxyServices) Stick(ctx context.Context, varg *types.Paging) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/%s/services?pag=%d&per_page=%d", strings.ToLower(varg.GetFactor("kind")), varg.Index, varg.Limit))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyServiceDescribe struct {
}

func (that *proxyServiceDescribe) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.service"}
}

func (that *proxyServiceDescribe) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.service"}
}

func (that *proxyServiceDescribe) I() *endpointDescriber {
	return new(endpointDescriber)
}

func (that *proxyServiceDescribe) O() json.RawMessage {
	return []byte{}
}

func (that *proxyServiceDescribe) Stick(ctx context.Context, varg *endpointDescriber) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/%s/services/%d", strings.ToLower(varg.Kind), varg.Index))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyMiddlewares struct {
}

func (that *proxyMiddlewares) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.middlewares"}
}

func (that *proxyMiddlewares) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.middlewares"}
}

func (that *proxyMiddlewares) I() *types.Paging {
	return new(types.Paging)
}

func (that *proxyMiddlewares) O() json.RawMessage {
	return []byte{}
}

func (that *proxyMiddlewares) Stick(ctx context.Context, varg *types.Paging) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/%s/middlewares?pag=%d&per_page=%d", strings.ToLower(varg.GetFactor("kind")), varg.Index, varg.Limit))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}

type proxyMiddlewareDescribe struct {
}

func (that *proxyMiddlewareDescribe) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.middleware"}
}

func (that *proxyMiddlewareDescribe) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.middleware"}
}

func (that *proxyMiddlewareDescribe) I() *endpointDescriber {
	return new(endpointDescriber)
}

func (that *proxyMiddlewareDescribe) O() json.RawMessage {
	return []byte{}
}

func (that *proxyMiddlewareDescribe) Stick(ctx context.Context, varg *endpointDescriber) (json.RawMessage, error) {
	buff, err := endpoint(fmt.Sprintf("api/%s/middlewares/%d", strings.ToLower(varg.Kind), varg.Index))
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff.Bytes(), nil
}
