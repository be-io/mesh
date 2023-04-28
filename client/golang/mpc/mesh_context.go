/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	var _ prsim.Context = new(MeshContext)
}

var (
	TimeoutKey     = &prsim.Key{Name: "mesh.mpc.timeout", Dft: func() interface{} { return time.Second * 10 }}
	AddressKey     = &prsim.Key{Name: "mesh.mpc.address", Dft: func() interface{} { return "" }}
	InsecureKey    = &prsim.Key{Name: "mesh.mpc.insecure", Dft: func() interface{} { return true }}
	CertificateKey = &prsim.Key{Name: "mesh.mpc.certificate", Dft: func() interface{} { return &tls.Config{} }}
	RemoteUname    = &prsim.Key{Name: "mesh.mpc.uname", Dft: func() interface{} { return "" }}
	RemoteName     = &prsim.Key{Name: "mesh.mpc.name", Dft: func() interface{} { return "" }}
	HeaderKey      = &prsim.Key{Name: "mesh.mpc.headers", Dft: func() interface{} { return map[string]string{} }}
)

type MeshContext struct {
	Context     context.Context
	TraceId     string
	SpanId      string
	Timestamp   int64
	RunMode     int
	URN         string
	Consumer    *types.Location
	calls       int
	Attachments map[string]string
	Attributes  map[string]interface{}
	Principals  *prsim.Deque[*types.Principal]
}

func (that *MeshContext) GetTraceId() string {
	return tool.Anyone(that.TraceId, prsim.MeshTraceId.Get(that.Attachments))
}

func (that *MeshContext) GetSpanId() string {
	return tool.Anyone(that.SpanId, prsim.MeshSpanId.Get(that.Attachments))
}

func (that *MeshContext) GetTimestamp() int64 {
	if 0 != that.Timestamp {
		return that.Timestamp
	}
	return time.Now().UnixMilli()
}

func (that *MeshContext) GetRunMode() int {
	if 0 != that.RunMode {
		return that.RunMode
	}
	return int(prsim.Routine)
}

func (that *MeshContext) GetUrn() string {
	return tool.Anyone(that.URN, prsim.MeshUrn.Get(that.Attachments))
}

func (that *MeshContext) GetConsumer(ctx context.Context) *types.Location {
	if nil != that.Consumer {
		return that.Consumer
	}
	return that.GetProvider(ctx)
}

func (that *MeshContext) GetProvider(ctx context.Context) *types.Location {
	return Locale(ctx)
}

func (that *MeshContext) GetAttachments() map[string]string {
	return that.Attachments
}

func (that *MeshContext) GetAttachment(name string) string {
	if "" != that.Attachments[name] {
		return that.Attachments[name]
	}
	for k, v := range that.Attachments {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	return ""
}

func (that *MeshContext) GetPrincipals() *prsim.Deque[*types.Principal] {
	if nil == that.Principals {
		that.Principals = &prsim.Deque[*types.Principal]{}
	}
	return that.Principals
}

func (that *MeshContext) GetAttributes() map[string]interface{} {
	return that.Attributes
}

func (that *MeshContext) GetAttribute(key *prsim.Key) interface{} {
	if nil == that.Attributes[key.Name] {
		return key.Dft()
	}
	return that.Attributes[key.Name]
}

func (that *MeshContext) SetAttribute(key *prsim.Key, value interface{}) {
	that.Attributes[key.Name] = value
}

func (that *MeshContext) RewriteURN(urn string) {
	that.URN = urn
	prsim.MeshUrn.Set(that.Attachments, urn)
}

func (that *MeshContext) RewriteContext(context prsim.Context) {
	if "" != context.GetTraceId() {
		that.TraceId = context.GetTraceId()
	}
	if "" != context.GetSpanId() {
		that.SpanId = context.GetSpanId()
	}
	if 0 != context.GetTimestamp() {
		that.Timestamp = context.GetTimestamp()
	}
	if 0 != context.GetRunMode() {
		that.RunMode = context.GetRunMode()
	}
	if "" != context.GetUrn() {
		that.URN = context.GetUrn()
	}
	if nil != context.GetConsumer(context) {
		that.Consumer = context.GetConsumer(context)
	}
	if nil != context.GetAttachments() {
		for key, value := range context.GetAttachments() {
			if "" != value {
				that.Attachments[key] = value
			}
		}
	}
	if nil != context.GetAttributes() {
		for key, value := range context.GetAttributes() {
			if nil != value && "" != value {
				that.Attributes[key] = value
			}
		}
	}
	if nil != that.GetPrincipals() {
		that.GetPrincipals().Add(context.GetPrincipals())
	}
}

func (that *MeshContext) Resume(ctx context.Context) prsim.Context {
	mtx := &MeshContext{
		Context:     ctx,
		Principals:  &prsim.Deque[*types.Principal]{},
		Attachments: map[string]string{},
		Attributes:  map[string]interface{}{},
		Consumer:    &types.Location{}}
	that.calls = that.calls + 1
	mtx.RewriteContext(that)
	mtx.SpanId = tool.NewSpanId(that.SpanId, that.calls)
	return mtx
}

func (that *MeshContext) WithTimeout(timeout time.Duration) context.CancelFunc {
	ttx, cancel := context.WithTimeout(that.Context, timeout)
	that.Context = ttx
	return cancel
}

func (that *MeshContext) Deadline() (deadline time.Time, ok bool) {
	return that.Context.Deadline()
}

func (that *MeshContext) Done() <-chan struct{} {
	return that.Context.Done()
}

func (that *MeshContext) Err() error {
	return that.Context.Err()
}

func (that *MeshContext) Value(key interface{}) interface{} {
	sk := fmt.Sprintf("%v", key)
	if nil != that.Attributes && nil != that.Attributes[sk] && "" != that.Attributes[sk] {
		return that.Attributes[sk]
	}
	if nil != that.Attachments && "" != that.Attachments[sk] {
		return that.Attachments[sk]
	}
	return that.Context.Value(key)
}

func CopyContext(ctx context.Context, source map[string]string) prsim.Context {
	mtx := ContextWith(ctx)
	for k, v := range source {
		mtx.GetAttachments()[k] = v
	}
	return mtx
}

func Context() prsim.Context {
	return contextWith(macro.Context())
}

func HTTPTracerContext(r *http.Request) prsim.Context {
	return TracerContext(prsim.MeshTraceId.GetHeader(r.Header), prsim.MeshSpanId.GetHeader(r.Header))
}

func TracerContext(traceId string, spanId string) prsim.Context {
	return contextWithTracer(macro.Context(), traceId, spanId)
}

func CloneContext(ctx context.Context) prsim.Context {
	mtx := Context()
	if c, ok := ctx.(prsim.Context); ok {
		mtx.RewriteContext(c)
	}
	return mtx
}

func WithContext(ctx context.Context, mtx prsim.Context) context.Context {
	return context.WithValue(ctx, "mesh.mpc.context", mtx)
}

func ContextWith(ctx context.Context) prsim.Context {
	if mtx, ok := ctx.(prsim.Context); ok {
		return mtx
	}
	if mtx, ok := ctx.Value("mesh.mpc.context").(prsim.Context); ok {
		return mtx
	}
	return contextWith(ctx)
}

func Dump(ctx context.Context) map[string]string {
	return ContextWith(ctx).GetAttachments()
}

func Locale(ctx context.Context) *types.Location {
	environ, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		log.Error(ctx, err.Error())
		environ = &types.Environ{}
	}
	return &types.Location{
		Principal: types.Principal{
			NodeId: environ.NodeId,
			InstId: environ.InstId,
		},
		IP:   tool.IP.Get(),
		Port: strconv.Itoa(tool.Runtime.Get().Port),
		Host: tool.Host.Get(),
		Name: tool.Name.Get(),
	}
}

func contextWithTracer(ctx context.Context, traceId string, spanId string) prsim.Context {
	attachments := map[string]string{}
	prsim.MeshTraceId.Set(attachments, traceId)
	prsim.MeshSpanId.Set(attachments, spanId)
	prsim.MeshTimestamp.Set(attachments, strconv.FormatInt(time.Now().UnixMilli(), 10))
	prsim.MeshRunMode.Set(attachments, strconv.Itoa(int(prsim.Routine)))
	prsim.MeshConsumer.Set(attachments, "{}")
	prsim.MeshProvider.Set(attachments, "{}")
	prsim.MeshUrn.Set(attachments, "")
	return &MeshContext{
		Context:     ctx,
		Principals:  &prsim.Deque[*types.Principal]{},
		Attachments: attachments,
		Attributes:  map[string]interface{}{},
		Consumer:    &types.Location{},
	}
}

func contextWith(ctx context.Context) prsim.Context {
	return contextWithTracer(ctx, tool.NewTraceId(), tool.NewSpanId("", 0))
}
