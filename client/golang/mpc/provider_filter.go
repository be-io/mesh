/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"math"
	"strconv"
	"time"
)

func init() {
	var _ Filter = new(providerFilter)
	macro.Provide(IFilter, new(providerFilter))
}

type providerFilter struct {
}

func (that *providerFilter) Att() *macro.Att {
	return &macro.Att{Name: PROVIDER, Pattern: PROVIDER, Priority: math.MaxInt}
}

func (that *providerFilter) Invoke(ctx context.Context, invoker Invoker, invocation Invocation) (interface{}, error) {
	attachments := invocation.GetParameters().GetAttachments(ctx)
	if nil == attachments {
		attachments = map[string]string{}
	}
	cdc, ok := macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	if !ok {
		return nil, cause.Errorf("No codec named %s exist ", codec.JSON)
	}
	traceId := prsim.MeshTraceId.Get(attachments)
	spanId := prsim.MeshSpanId.Get(attachments)
	timestamp := that.resolveTimestamp(ctx, prsim.MeshTimestamp.Get(attachments))
	runMode := prsim.FromString(prsim.MeshRunMode.Get(attachments))
	urn := prsim.MeshUrn.Get(attachments)
	var pp types.Location
	if _, err := cdc.DecodeString(prsim.MeshProvider.Get(attachments), &pp); nil != err {
		log.Error(ctx, err.Error())
		pp = types.Location{}
	}
	if nil != invocation.GetURN() {
		urn = invocation.GetURN().String()
	}
	ntx := &MeshContext{
		Context:     ctx,
		TraceId:     traceId,
		SpanId:      spanId,
		Timestamp:   timestamp,
		RunMode:     runMode,
		URN:         urn,
		Consumer:    &pp,
		Attachments: map[string]string{},
	}
	for key, value := range attachments {
		ntx.GetAttachments()[key] = value
	}
	mtx := ContextWith(ctx)
	mtx.RewriteContext(ntx)

	digest := NewDigest(mtx)

	ret, err := invoker.Invoke(mtx, invocation)
	if nil != err {
		digest.Print(mtx, "P", cause.Coder(err))
	} else {
		digest.Print(mtx, "P", cause.Success.Code)
	}
	return ret, err
}

func (that *providerFilter) Name() string {
	return PROVIDER
}

func (that *providerFilter) Pattern() string {
	return PROVIDER
}

func (that *providerFilter) resolveTimestamp(ctx context.Context, v string) int64 {
	if "" == v {
		return time.Now().UnixMilli()
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if nil != err {
		log.Error(ctx, err.Error())
		return time.Now().UnixMilli()
	}
	return n
}
