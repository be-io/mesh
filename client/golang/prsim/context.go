/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://www.rfc-editor.org/rfc/rfc7540#section-8.1.2

const (
	MeshSpanId        Metadata = "mesh-span-id"
	MeshTimestamp     Metadata = "mesh-timestamp"
	MeshRunMode       Metadata = "mesh-run-mode"
	MeshConsumer      Metadata = "mesh-consumer"
	MeshProvider      Metadata = "mesh-provider"
	MeshUrn           Metadata = "mesh-urn"
	MeshIncomingHost  Metadata = "mesh-incoming-host"
	MeshOutgoingHost  Metadata = "mesh-outgoing-host"
	MeshIncomingProxy Metadata = "mesh-incoming-proxy"
	MeshOutgoingProxy Metadata = "mesh-outgoing-proxy"
	MeshSubset        Metadata = "mesh-subset"
	MeshPath          Metadata = "mesh-path"
	MeshGray          Metadata = "mesh-gray"
	MeshProto         Metadata = "mesh-proto"

	// INC

	MeshVersion          Metadata = "x-ptp-version"
	MeshTechProviderCode Metadata = "x-ptp-tech-provider-code"
	MeshTraceId          Metadata = "x-ptp-trace-id"
	MeshToken            Metadata = "x-ptp-token"
	MeshURI              Metadata = "x-ptp-uri"
	MeshFromInstId       Metadata = "x-ptp-source-inst-id"
	MeshFromNodeId       Metadata = "x-ptp-source-node-id"
	MeshTargetNodeId     Metadata = "x-ptp-target-node-id"
	MeshTargetInstId     Metadata = "x-ptp-target-inst-id"
	MeshSessionId        Metadata = "x-ptp-session-id"
	MeshTopic            Metadata = "x-ptp-topic"
	MeshTimeout          Metadata = "x-ptp-timeout"

	// RFC

	ContentLength   Metadata = "content-length"
	ContentEncoding Metadata = "content-encoding"
	ContentType     Metadata = "content-type"
)

var names = []Metadata{
	MeshSpanId,
	MeshTimestamp,
	MeshRunMode,
	MeshConsumer,
	MeshProvider,
	MeshUrn,
	MeshIncomingHost,
	MeshOutgoingHost,
	MeshIncomingProxy,
	MeshOutgoingProxy,
	MeshSubset,
	MeshPath,
	MeshVersion,
	MeshTechProviderCode,
	MeshTraceId,
	MeshToken,
	MeshFromInstId,
	MeshFromNodeId,
	MeshTargetNodeId,
	MeshTargetInstId,
	MeshSessionId,
}

type Metadata string

func (that Metadata) Keys(prefix string) []string {
	return []string{
		fmt.Sprintf("%s%s", prefix, http.CanonicalHeaderKey(string(that))),
		fmt.Sprintf("%s%s", prefix, http.CanonicalHeaderKey(strings.ReplaceAll(string(that), "-", "_"))),
		fmt.Sprintf("%s%s", prefix, string(that)),
		fmt.Sprintf("%s%s", prefix, tool.FistUpper(string(that))),
		fmt.Sprintf("%s%s", prefix, tool.FistUpper(strings.ToLower(strings.ReplaceAll(string(that), "-", "_")))),
	}
}

func (that Metadata) Get(attachments map[string]string) string {
	if nil != attachments {
		if v := attachments[string(that)]; "" != v {
			return v
		}
		if v := attachments[http.CanonicalHeaderKey(string(that))]; "" != v {
			return v
		}
		if v := attachments[strings.ReplaceAll(string(that), "-", "_")]; "" != v {
			return v
		}
		if v := attachments[strings.ReplaceAll(string(that), "_", "-")]; "" != v {
			return v
		}
		for k, v := range attachments {
			if string(that) == strings.ToLower(k) || string(that) == strings.ToLower(strings.ReplaceAll(k, "_", "-")) {
				return v
			}
		}
	}
	return ""
}

func (that Metadata) Set(attachments map[string]string, v string) {
	if nil != attachments && "" != v {
		hk := http.CanonicalHeaderKey(string(that))
		if _, ok := attachments[hk]; ok {
			attachments[hk] = v
		} else {
			attachments[string(that)] = v
		}
	}
}

func (that Metadata) Del(attachments map[string]string) {
	var keys []string
	for k, _ := range attachments {
		if strings.EqualFold(k, string(that)) {
			keys = append(keys, k)
		}
	}
	for _, key := range keys {
		delete(attachments, key)
	}
}

func (that Metadata) GetHeader(headers map[string][]string) string {
	if nil != headers {
		if v := headers[string(that)]; len(v) > 0 {
			return tool.Anyone(v...)
		}
		if v := headers[http.CanonicalHeaderKey(string(that))]; len(v) > 0 {
			return tool.Anyone(v...)
		}
		if v := headers[strings.ReplaceAll(string(that), "_", "-")]; len(v) > 0 {
			return tool.Anyone(v...)
		}
		if v := headers[strings.ReplaceAll(string(that), "-", "_")]; len(v) > 0 {
			return tool.Anyone(v...)
		}
		for k, v := range headers {
			if len(v) < 1 {
				continue
			}
			if string(that) == strings.ToLower(k) || string(that) == strings.ToLower(strings.ReplaceAll(k, "_", "-")) {
				return tool.Anyone(v...)
			}
		}
	}
	return ""
}

func (that Metadata) SetHeader(headers map[string][]string, v string) {
	if nil != headers && "" != v {
		that.DelHeader(headers)
		http.Header(headers).Set(string(that), v)
	}
}

func (that Metadata) DelHeader(headers map[string][]string) {
	if nil != headers {
		delete(headers, http.CanonicalHeaderKey(string(that)))
		delete(headers, string(that))
	}
}

func CopyHeadersOverride(dst http.Header, src http.Header) {
	for k, vv := range src {
		if len(vv) > 0 {
			dst[k] = vv
		}
	}
}

func IsMeshMetadata(name string) bool {
	for _, x := range names {
		if strings.EqualFold(string(x), name) {
			return true
		}
	}
	return false
}

func DstNodeId(ctx context.Context, headers http.Header) string {
	if x := MeshTargetNodeId.GetHeader(headers); "" != x {
		return x
	}
	if x := MeshTargetInstId.GetHeader(headers); "" != x {
		return x
	}
	if x := MeshUrn.GetHeader(headers); "" != x {
		return types.FromURN(ctx, x).NodeId
	}
	return ""
}

func SrcNodeId(ctx context.Context, headers http.Header) string {
	if x := MeshFromNodeId.GetHeader(headers); "" != x {
		return x
	}
	if x := MeshFromInstId.GetHeader(headers); "" != x {
		return x
	}
	return ""
}

func SetMetadata(ctx Context, dict map[string][]string) {
	MeshTraceId.SetHeader(dict, ctx.GetTraceId())
	MeshSpanId.SetHeader(dict, ctx.GetSpanId())
	MeshFromInstId.SetHeader(dict, ctx.GetConsumer(ctx).InstId)
	MeshFromNodeId.SetHeader(dict, ctx.GetConsumer(ctx).NodeId)
	MeshIncomingHost.SetHeader(dict, fmt.Sprintf("%s@%s", tool.Name.Get(), tool.Runtime.Get()))
	MeshOutgoingHost.SetHeader(dict, MeshIncomingHost.Get(ctx.GetAttachments()))
	MeshIncomingProxy.SetHeader(dict, MeshIncomingProxy.Get(ctx.GetAttachments()))
	MeshOutgoingProxy.SetHeader(dict, MeshOutgoingProxy.Get(ctx.GetAttachments()))
	MeshUrn.SetHeader(dict, ctx.GetUrn())
	MeshSubset.SetHeader(dict, MeshSubset.Get(ctx.GetAttachments()))
	MeshPath.SetHeader(dict, MeshPath.Get(ctx.GetAttachments()))
	MeshVersion.SetHeader(dict, MeshVersion.Get(ctx.GetAttachments()))
	MeshTimestamp.SetHeader(dict, fmt.Sprintf("%d", ctx.GetTimestamp()))
	MeshRunMode.SetHeader(dict, fmt.Sprintf("%d", ctx.GetRunMode()))
	// INC
	MeshTechProviderCode.SetHeader(dict, MeshTechProviderCode.Get(ctx.GetAttachments()))
	MeshToken.SetHeader(dict, MeshToken.Get(ctx.GetAttachments()))
	MeshTargetNodeId.SetHeader(dict, tool.Anyone(types.FromURN(ctx, ctx.GetUrn()).NodeId, MeshTargetNodeId.Get(ctx.GetAttachments())))
	MeshTargetInstId.SetHeader(dict, MeshTargetInstId.Get(ctx.GetAttachments()))
	MeshSessionId.SetHeader(dict, MeshSessionId.Get(ctx.GetAttachments()))
	MeshTopic.SetHeader(dict, MeshTopic.Get(ctx.GetAttachments()))
	MeshTimeout.SetHeader(dict, MeshTimeout.Get(ctx.GetAttachments()))
	MeshURI.SetHeader(dict, MeshURI.Get(ctx.GetAttachments()))
}

func GetMetadata(dict map[string][]string) map[string]string {
	attachments := map[string]string{}
	MeshTraceId.Set(attachments, MeshTraceId.GetHeader(dict))
	MeshSpanId.Set(attachments, MeshSpanId.GetHeader(dict))
	MeshFromInstId.Set(attachments, MeshFromInstId.GetHeader(dict))
	MeshFromNodeId.Set(attachments, MeshFromNodeId.GetHeader(dict))
	MeshIncomingHost.Set(attachments, MeshIncomingHost.GetHeader(dict))
	MeshOutgoingHost.Set(attachments, MeshOutgoingHost.GetHeader(dict))
	MeshIncomingProxy.Set(attachments, MeshIncomingProxy.GetHeader(dict))
	MeshOutgoingProxy.Set(attachments, MeshOutgoingProxy.GetHeader(dict))
	MeshUrn.Set(attachments, MeshUrn.GetHeader(dict))
	MeshSubset.Set(attachments, MeshSubset.GetHeader(dict))
	MeshPath.Set(attachments, MeshPath.GetHeader(dict))
	MeshVersion.Set(attachments, MeshVersion.GetHeader(dict))
	MeshTimestamp.Set(attachments, MeshTimestamp.GetHeader(dict))
	MeshRunMode.Set(attachments, MeshRunMode.GetHeader(dict))
	// INC
	MeshTechProviderCode.Set(attachments, MeshTechProviderCode.GetHeader(dict))
	MeshToken.Set(attachments, MeshToken.GetHeader(dict))
	MeshTargetNodeId.Set(attachments, MeshTargetNodeId.GetHeader(dict))
	MeshTargetInstId.Set(attachments, MeshTargetInstId.GetHeader(dict))
	MeshSessionId.Set(attachments, MeshSessionId.GetHeader(dict))
	MeshTopic.Set(attachments, MeshTopic.GetHeader(dict))
	MeshTimeout.Set(attachments, MeshTimeout.GetHeader(dict))
	MeshURI.Set(attachments, MeshURI.GetHeader(dict))
	return attachments
}

type Context interface {
	context.Context

	// GetTraceId the request trace id.
	GetTraceId() string

	// GetSpanId the request span id.
	GetSpanId() string

	// GetTimestamp the request create time.
	GetTimestamp() int64

	// GetRunMode the request run mode. {@link RunMode}
	GetRunMode() int

	// GetUrn mesh resource uniform name. Like: create.tenant.omega.json.http2.lx000001.mpi.trustbe.net
	GetUrn() string

	// GetConsumer the consumer network principal.
	GetConsumer(ctc context.Context) *types.Location

	// GetProvider the provider network principal.
	GetProvider(ctx context.Context) *types.Location

	// GetAttachments Dispatch attachments.
	GetAttachments() map[string]string

	// GetAttachment get with name sensitive.
	GetAttachment(name string) string

	// GetPrincipals Get the mpc broadcast network principals.
	GetPrincipals() *Deque[*types.Principal]

	// GetAttributes like getAttachments, but attribute wont be transfer in invoke chain.
	GetAttributes() map[string]interface{}

	// GetAttribute like getAttachments, but attribute wont be transfer in invoke chain.
	GetAttribute(key *Key) interface{}

	// SetAttribute Like putAttachments, but attribute won't be transfer in invoke chain.
	SetAttribute(key *Key, value interface{})

	// RewriteURN rewrite the urn.
	RewriteURN(urn string)

	// RewriteContext rewrite the context by another context.
	RewriteContext(context Context)

	// Resume will open a new context.
	Resume(ctx context.Context) Context

	// WithTimeout set the deadline time
	WithTimeout(timeout time.Duration) context.CancelFunc
}

type RunMode int

const (
	// Routine 正常模式
	Routine RunMode = 1
	// Perform 评测模式
	Perform RunMode = 2
	// Defense 高防模式
	Defense RunMode = 4
	// Debug 调试模式,
	Debug RunMode = 8
	// LoadTest 压测模式
	LoadTest RunMode = 16
	// Mock Mock模式
	Mock RunMode = 32
)

func (that RunMode) Matches(code int) bool {
	return (int(that) & code) == int(that)
}

func FromString(code string) int {
	v, err := strconv.Atoi(code)
	if nil != err {
		log.Error0(err.Error())
		return int(Routine)
	}
	return v
}

type Deque[T any] struct {
	data []T
}

func (that *Deque[T]) List() []T {
	return that.data
}

func (that *Deque[T]) Pop() T {
	if len(that.data) < 1 {
		var t T
		return t
	}
	element := that.data[len(that.data)-1]
	that.data = that.data[0 : len(that.data)-1]
	return element
}

func (that *Deque[T]) Push(element T) {
	that.data = append(that.data, element)
}

func (that *Deque[T]) Peek() T {
	if len(that.data) < 1 {
		var t T
		return t
	}
	return that.data[len(that.data)-1]
}

func (that *Deque[T]) Add(queue *Deque[T]) {
	that.data = append(that.data, queue.data...)
}

type Key struct {
	Name string
	Dft  func() interface{}
}
