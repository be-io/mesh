/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package ptp

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/grpc"
	httpx "github.com/be-io/mesh/client/golang/http"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/ugorji/go/codec"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func init() {
	var _ http.Handler = new(refresh)
	macro.Provide(grpc.IHandler, new(refresh))

	var _ http.Handler = new(register)
	macro.Provide(grpc.IHandler, new(register))

	var _ http.Handler = new(update)
	macro.Provide(grpc.IHandler, new(update))

	var _ http.Handler = new(weave)
	macro.Provide(grpc.IHandler, new(weave))

	var _ http.Handler = new(ack)
	macro.Provide(grpc.IHandler, new(ack))

	var _ http.Handler = pushService
	macro.Provide(grpc.IHandler, pushService)

	var _ http.Handler = popService
	macro.Provide(grpc.IHandler, popService)

	var _ http.Handler = peekService
	macro.Provide(grpc.IHandler, peekService)

	var _ http.Handler = releaseService
	macro.Provide(grpc.IHandler, releaseService)

	var _ http.Handler = new(invoke)
	macro.Provide(grpc.IHandler, new(invoke))

	var _ http.Handler = new(transport)
	macro.Provide(grpc.IHandler, new(transport))
}

var (
	pushService    = new(push)
	popService     = new(pop)
	peekService    = new(peek)
	releaseService = new(release)
)

func WithHTTPBound(ctx context.Context, fn func() ([]byte, error)) *TransportOutbound {
	buff, err := fn()
	if nil == err {
		return &TransportOutbound{
			Metadata: map[string]string{},
			Payload:  buff,
			Code:     cause.Success.Code,
			Message:  cause.Success.Message,
		}
	}
	log.Warn(ctx, err.Error())
	code, msg := cause.Parse(err)
	return &TransportOutbound{
		Metadata: map[string]string{},
		Payload:  buff,
		Code:     code,
		Message:  msg,
	}
}

func ServeHTTP[T any](w http.ResponseWriter, r *http.Request, fn func(ctx prsim.Context, contentType string, input *T) ([]byte, error)) {
	ctx := mpc.HTTPTracerContext(r)
	for k, v := range prsim.GetMetadata(r.Header) {
		ctx.GetAttachments()[k] = v
	}
	prsim.ContentLength.Del(ctx.GetAttachments())
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, "%v", err)
			log.Error(ctx, string(debug.Stack()))
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		}
	}()
	contentType := httpx.FormatContentType(r.Header.Get("Content-Type"))
	binding := httpx.Default(r.Method, contentType)
	var input T
	if err := binding.Bind(r, &input); nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := fn(ctx, contentType, &input)
	if ok := cause.WriteHTTPError(err, w); ok {
		return
	}
	v := WithHTTPBound(ctx, func() ([]byte, error) { return b, err })
	buff, err := Encode(v, contentType)
	if nil != err {
		log.Warn(ctx, err.Error())
		http.Error(w, cause.DeError(err).Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType)
	if _, err = w.Write(buff); nil != err {
		log.Warn(ctx, err.Error())
		http.Error(w, cause.DeError(err).Error(), http.StatusInternalServerError)
		return
	}
}

func TransportHTTP(ctx prsim.Context, input any, uri string, contentType string) ([]byte, error) {
	if err := aware.WeavedCheck(ctx); nil != err {
		return nil, cause.Error(err)
	}
	b, err := Encode(input, contentType)
	if nil != err {
		return nil, cause.Error(err)
	}
	i := &Inbound{
		Metadata: ctx.GetAttachments(),
		Payload:  b,
	}
	bi, err := Encode(i, contentType)
	if nil != err {
		return nil, cause.Error(err)
	}
	in, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://127.0.0.1:7304/v1/interconn/chan/invoke", bytes.NewBuffer(bi))
	if nil != err {
		return nil, cause.Error(err)
	}
	prsim.SetMetadata(ctx, in.Header)
	in.Header.Set("Content-Type", contentType)
	prsim.MeshURI.SetHeader(in.Header, uri)
	out, err := tool.Client.Do(in)
	if nil != err {
		return nil, cause.Error(err)
	}
	defer func() { log.Catch(out.Body.Close()) }()
	bo, err := io.ReadAll(out.Body)
	if nil != err {
		return nil, cause.Error(err)
	}
	if out.StatusCode != http.StatusOK {
		return nil, cause.Errorh(out.StatusCode, string(bo))
	}
	oct := tool.Anyone(prsim.ContentType.Get(ctx.GetAttachments()), out.Header.Get("Content-Type"))
	o := new(Outbound)
	if err = Decode(bo, o, oct); nil != err {
		return nil, cause.Error(err)
	}
	if o.Code != cause.Success.Code {
		return nil, cause.Errorm(o.Code, o.Message)
	}
	return o.Payload, nil
}

func TransportGRPC(ctx prsim.Context, input any, uri string) (*Outbound, error) {
	if err := aware.WeavedCheck(ctx); nil != err {
		return nil, cause.Error(err)
	}
	b, err := Encode(input, httpx.MIMEPROTOBUF)
	if nil != err {
		return nil, cause.Error(err)
	}
	i := &Inbound{
		Metadata: ctx.GetAttachments(),
		Payload:  b,
	}
	prsim.MeshURI.Set(ctx.GetAttachments(), uri)
	c, cancel, err := aware.Channel.DialContext(ctx, "127.0.0.1:7304")
	defer cancel()
	if nil != err {
		return nil, cause.Error(err)
	}
	defer func() { log.Catch(c.Close()) }()
	return NewPrivateTransferProtocolClient(c).Invoke(ctx, i)
}

func Encode(v any, contentType string) ([]byte, error) {
	switch contentType {
	case httpx.MIMEJSON:
		return json.Marshal(v)
	case httpx.MIMEXML, httpx.MIMEXML2:
		return xml.Marshal(v)
	case httpx.MIMEPROTOBUF:
		if m, ok := v.(proto.Message); ok {
			return proto.Marshal(m)
		}
		return nil, cause.Compatible.Error()
	case httpx.MIMEMSGPACK, httpx.MIMEMSGPACK2:
		buff := &bytes.Buffer{}
		cdc := new(codec.MsgpackHandle)
		if err := codec.NewEncoder(buff, cdc).Encode(v); err != nil {
			return nil, err
		}
		return buff.Bytes(), nil
	case httpx.MIMEYAML:
		return yaml.Marshal(v)
	case httpx.MIMETOML:
		return toml.Marshal(v)
	case httpx.MIMEMultipartPOSTForm:
		return json.Marshal(v)
	default: // case MIMEPOSTForm:
		return json.Marshal(v)
	}
}

func Decode(b []byte, ptr any, contentType string) error {
	switch contentType {
	case httpx.MIMEJSON:
		return json.Unmarshal(b, ptr)
	case httpx.MIMEXML, httpx.MIMEXML2:
		return xml.Unmarshal(b, ptr)
	case httpx.MIMEPROTOBUF:
		if m, ok := ptr.(proto.Message); ok {
			return proto.Unmarshal(b, m)
		}
		return cause.Compatible.Error()
	case httpx.MIMEMSGPACK, httpx.MIMEMSGPACK2:
		buff := &bytes.Buffer{}
		cdc := new(codec.MsgpackHandle)
		return codec.NewDecoder(buff, cdc).Decode(ptr)
	case httpx.MIMEYAML:
		return yaml.Unmarshal(b, ptr)
	case httpx.MIMETOML:
		return toml.Unmarshal(b, ptr)
	case httpx.MIMEMultipartPOSTForm:
		return json.Unmarshal(b, ptr)
	default: // case MIMEPOSTForm:
		return json.Unmarshal(b, ptr)
	}
}

type refresh struct {
}

func (that *refresh) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.node.refresh",
		Pattern: "/v1/interconn/node/refresh",
	}
}

func (that *refresh) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, func(ctx prsim.Context, contentType string, input *types.Environ) ([]byte, error) {
		return nil, aware.KMS.Reset(ctx, input)
	})
}

type register struct {
}

func (that *register) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.registry.register",
		Pattern: "/v1/interconn/registry/register",
	}
}

func (that *register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, func(ctx prsim.Context, contentType string, input *types.Registration[any]) ([]byte, error) {
		if nil == input || "" == input.InstanceId || "" == input.Address || "" == input.Name || "" == input.Kind || nil == input.Content {
			return nil, cause.Validate.Error()
		}
		return nil, aware.Registry.Register(ctx, input)
	})
}

type update struct {
}

func (that *update) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.net.refresh",
		Pattern: "/v1/interconn/net/refresh",
	}
}

func (that *update) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, func(ctx prsim.Context, contentType string, input *types.Route) ([]byte, error) {
		if nil == input || "" == input.NodeId || "" == input.InstId || "" == input.Address {
			return nil, cause.Validate.Error()
		}
		return nil, aware.Network.Refresh(ctx, []*types.Route{input})
	})
}

type weave struct {
}

func (that *weave) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.net.weave",
		Pattern: "/v1/interconn/net/weave",
	}
}

func (that *weave) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, func(ctx prsim.Context, contentType string, input *types.Route) ([]byte, error) {
		if nil == input || "" == input.NodeId {
			return nil, cause.Validate.Error()
		}
		return nil, aware.Network.Weave(ctx, input)
	})
}

type ack struct {
}

func (that *ack) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.net.ack",
		Pattern: "/v1/interconn/net/ack",
	}
}

func (that *ack) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, func(ctx prsim.Context, contentType string, input *types.Route) ([]byte, error) {
		if nil == input || "" == input.NodeId {
			return nil, cause.Validate.Error()
		}
		return nil, aware.Network.Ack(ctx, input)
	})
}

type push struct {
}

func (that *push) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.chan.push",
		Pattern: "/v1/interconn/chan/push",
	}
}

func (that *push) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, that.ServePush)
}

func (that *push) ServePush(ctx prsim.Context, contentType string, input *PushInbound) ([]byte, error) {
	if nil == input || nil == input.Payload {
		return nil, cause.Validate.Error()
	}
	env, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	if x := prsim.MeshTargetNodeId.Get(ctx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
		topic := tool.Anyone(prsim.MeshTopic.Get(ctx.GetAttachments()), input.Topic)
		return nil, aware.Session.Push(ctx, input.Payload, input.Metadata, topic)
	}
	return TransportHTTP(ctx, input, fmt.Sprintf("https://ptp.cn%s", that.Att().Pattern), contentType)
}

type pop struct {
}

func (that *pop) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.chan.pop",
		Pattern: "/v1/interconn/chan/pop",
	}
}

func (that *pop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, that.ServePop)
}

func (that *pop) ServePop(ctx prsim.Context, contentType string, input *PopInbound) ([]byte, error) {
	if nil == input {
		return nil, cause.Validate.Error()
	}
	env, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	if x := prsim.MeshTargetNodeId.Get(ctx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
		timeout := types.Duration(time.Duration(input.Timeout) * time.Millisecond)
		if x := prsim.MeshTimeout.Get(ctx.GetAttachments()); "" != x {
			if t, err := strconv.Atoi(x); nil == err && t > 0 {
				timeout = types.Duration(time.Duration(t) * time.Millisecond)
			}
		}
		topic := tool.Anyone(prsim.MeshTopic.Get(ctx.GetAttachments()), input.Topic)
		return aware.Session.Pop(ctx, timeout, topic)
	}
	return TransportHTTP(ctx, input, fmt.Sprintf("https://ptp.cn%s", that.Att().Pattern), contentType)
}

type peek struct {
}

func (that *peek) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.chan.peek",
		Pattern: "/v1/interconn/chan/peek",
	}
}

func (that *peek) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, that.ServePeek)
}

func (that *peek) ServePeek(ctx prsim.Context, contentType string, input *PeekInbound) ([]byte, error) {
	if nil == input {
		return nil, cause.Validate.Error()
	}
	env, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	if x := prsim.MeshTargetNodeId.Get(ctx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
		topic := tool.Anyone(prsim.MeshTopic.Get(ctx.GetAttachments()), input.Topic)
		return aware.Session.Peek(ctx, topic)
	}
	return TransportHTTP(ctx, input, fmt.Sprintf("https://ptp.cn%s", that.Att().Pattern), contentType)
}

type release struct {
}

func (that *release) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.chan.release",
		Pattern: "/v1/interconn/chan/release",
	}
}

func (that *release) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, that.ServeRelease)
}

func (that *release) ServeRelease(ctx prsim.Context, contentType string, input *ReleaseInbound) ([]byte, error) {
	if nil == input {
		return nil, cause.Validate.Error()
	}
	env, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	if x := prsim.MeshTargetNodeId.Get(ctx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
		timeout := types.Duration(time.Duration(input.Timeout) * time.Millisecond)
		if x := prsim.MeshTimeout.Get(ctx.GetAttachments()); "" != x {
			if t, err := strconv.Atoi(x); nil == err && t > 0 {
				timeout = types.Duration(time.Duration(t) * time.Millisecond)
			}
		}
		topic := tool.Anyone(prsim.MeshTopic.Get(ctx.GetAttachments()), input.Topic)
		return nil, aware.Session.Release(ctx, timeout, topic)
	}
	return TransportHTTP(ctx, input, fmt.Sprintf("https://ptp.cn%s", that.Att().Pattern), contentType)
}

type invoke struct {
}

func (that *invoke) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.chan.invoke",
		Pattern: "/v1/interconn/chan/invoke",
	}
}

func (that *invoke) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, func(ctx prsim.Context, contentType string, input *Inbound) ([]byte, error) {
		if nil == input {
			return nil, cause.Validate.Error()
		}
		uri, err := types.FormatURL(prsim.MeshURI.Get(ctx.GetAttachments()))
		if nil != err {
			return nil, cause.Error(err)
		}
		switch uri.Path {
		case "/v1/interconn/chan/pop":
			pi := new(PopInbound)
			if err = Decode(input.Payload, pi, contentType); nil != err {
				return nil, cause.Error(err)
			}
			return popService.ServePop(ctx, contentType, pi)
		case "/v1/interconn/chan/push":
			pi := new(PushInbound)
			if err = Decode(input.Payload, pi, contentType); nil != err {
				return nil, cause.Error(err)
			}
			return pushService.ServePush(ctx, contentType, pi)
		case "/v1/interconn/chan/peek":
			pi := new(PeekInbound)
			if err = Decode(input.Payload, pi, contentType); nil != err {
				return nil, cause.Error(err)
			}
			return peekService.ServePeek(ctx, contentType, pi)
		case "/v1/interconn/chan/release":
			pi := new(ReleaseInbound)
			if err = Decode(input.Payload, pi, contentType); nil != err {
				return nil, cause.Error(err)
			}
			return releaseService.ServeRelease(ctx, contentType, pi)
		case "/v1/interconn/net/weave":
			var route types.Route
			if _, err = aware.Codec.Decode(bytes.NewBuffer(input.Payload), &route); nil != err {
				return nil, cause.Error(err)
			}
			return nil, aware.Network.Weave(ctx, &route)
		case "/v1/interconn/net/ack":
			var route types.Route
			if _, err = aware.Codec.Decode(bytes.NewBuffer(input.Payload), &route); nil != err {
				return nil, cause.Error(err)
			}
			return nil, aware.Network.Ack(ctx, &route)
		default:
			return nil, cause.NotFound.Error()
		}
	})
}

type transport struct {
}

func (that *transport) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.ptp.chan.transport",
		Pattern: "/v1/interconn/chan/transport",
	}
}

func (that *transport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(w, r, func(ctx prsim.Context, contentType string, input *Inbound) ([]byte, error) {
		if nil == input {
			return nil, cause.Validate.Error()
		}
		return nil, cause.UrnNotPermit.Error()
	})
}
