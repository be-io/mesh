/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"math"
	"strconv"
)

func init() {
	var _ Filter = new(consumerFilter)
	macro.Provide(IFilter, new(consumerFilter))
}

var ConsumerKey = &prsim.Key{Name: "mesh.mpc.Consumer", Dft: func() interface{} { return nil }}

type consumerFilter struct {
}

func (that *consumerFilter) Att() *macro.Att {
	return &macro.Att{Name: CONSUMER, Pattern: CONSUMER, Priority: math.MinInt}
}

func (that *consumerFilter) Invoke(ctx context.Context, invoker Invoker, invocation Invocation) (interface{}, error) {
	mtx := ContextWith(ctx)
	for key, value := range mtx.GetAttachments() {
		invocation.GetAttachments(mtx)[key] = value
	}
	cc, err := aware.JSON.EncodeString(mtx.GetConsumer(mtx))
	if nil != err {
		return nil, cause.Error(err)
	}
	pp, err := aware.JSON.EncodeString(func() *types.Location {
		attr := mtx.GetAttribute(ConsumerKey)
		if nil == attr {
			return mtx.GetProvider(mtx)
		}
		if locate, ok := attr.(*types.Location); ok {
			return locate
		}
		return mtx.GetProvider(mtx)
	}())
	if nil != err {
		return nil, cause.Error(err)
	}
	prsim.MeshTraceId.Set(invocation.GetAttachments(mtx), mtx.GetTraceId())
	prsim.MeshSpanId.Set(invocation.GetAttachments(mtx), mtx.GetSpanId())
	prsim.MeshTimestamp.Set(invocation.GetAttachments(mtx), strconv.FormatInt(mtx.GetTimestamp(), 10))
	prsim.MeshRunMode.Set(invocation.GetAttachments(mtx), strconv.Itoa(mtx.GetRunMode()))
	prsim.MeshUrn.Set(invocation.GetAttachments(mtx), mtx.GetUrn())
	prsim.MeshConsumer.Set(invocation.GetAttachments(mtx), cc)
	prsim.MeshProvider.Set(invocation.GetAttachments(mtx), pp)

	digest := NewDigest(mtx)

	ret, err := invoker.Invoke(mtx, invocation)
	if nil != err {
		digest.Print(mtx, "C", cause.Coder(err))
	} else {
		digest.Print(mtx, "C", cause.Success.Code)
	}
	return ret, err
}

func (that *consumerFilter) Name() string {
	return CONSUMER
}

func (that *consumerFilter) Pattern() string {
	return CONSUMER
}
