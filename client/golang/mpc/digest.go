/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/prsim"
	"strconv"
	"time"
)

type Digest struct {
	TraceId string
	SpanId  string
	Mode    string
	CNI     string
	CII     string
	CIP     string
	CHost   string
	PNI     string
	PII     string
	PIP     string
	PHost   string
	URN     string
	Now     int64
}

func NewDigest(ctx prsim.Context) *Digest {
	return &Digest{
		TraceId: ctx.GetTraceId(),
		SpanId:  ctx.GetSpanId(),
		Mode:    strconv.Itoa(ctx.GetRunMode()),
		CNI:     ctx.GetConsumer(ctx).NodeId,
		CII:     ctx.GetConsumer(ctx).InstId,
		CIP:     ctx.GetConsumer(ctx).IP,
		CHost:   ctx.GetConsumer(ctx).Host,
		PNI:     ctx.GetProvider(ctx).NodeId,
		PII:     ctx.GetProvider(ctx).InstId,
		PIP:     ctx.GetProvider(ctx).IP,
		PHost:   ctx.GetProvider(ctx).Host,
		URN:     ctx.GetUrn(),
		Now:     time.Now().UnixMilli(),
	}
}

func (that *Digest) Print(ctx prsim.Context, pattern string, code string) {
	now := time.Now().UnixMilli()
	log.Info(ctx, "%s,%s,%d,%d,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%v,%s,%s",
		that.TraceId,
		that.SpanId,
		ctx.GetTimestamp(),
		that.Now,
		now-ctx.GetTimestamp(),
		now-that.Now,
		that.Mode,
		pattern,
		that.CNI,
		that.CII,
		that.PNI,
		that.PII,
		that.CIP,
		that.PIP,
		that.CHost,
		that.PHost,
		ctx.GetAttribute(AddressKey),
		that.URN,
		code,
	)
}
