/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"context"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	var _ grpc.StreamClientInterceptor = Interceptors.ClientStream
	var _ grpc.StreamServerInterceptor = Interceptors.ServerStream
	var _ grpc.UnaryClientInterceptor = Interceptors.ClientUnary
	var _ grpc.UnaryServerInterceptor = Interceptors.ServerUnary
}

var Interceptors = new(grpcInterceptor)

type grpcInterceptor struct {
}

func (that *grpcInterceptor) ClientContext(ctx context.Context) context.Context {
	mtx := mpc.ContextWith(ctx)
	attachments := make(map[string][]string, 19)
	prsim.SetMetadata(mtx, attachments)
	var kvs []string
	headers, ok := mtx.GetAttribute(mpc.HeaderKey).(map[string]string)
	if ok {
		for k, v := range headers {
			kvs = append(kvs, k, v)
		}
	}
	for k, v := range attachments {
		if len(v) > 0 {
			kvs = append(kvs, k, v[0])
		}
	}
	return metadata.AppendToOutgoingContext(ctx, kvs...)
}

func (that *grpcInterceptor) ClientStream(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, stream grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return stream(that.ClientContext(ctx), desc, cc, method, opts...)
}

func (that *grpcInterceptor) ClientUnary(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoke grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoke(that.ClientContext(ctx), method, req, reply, cc, opts...)
}

func (that *grpcInterceptor) ServerStream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handle grpc.StreamHandler) error {
	return handle(srv, ss)
}

func (that *grpcInterceptor) ServerUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handle grpc.UnaryHandler) (resp interface{}, err error) {
	return handle(that.ServerContext(ctx), req)
}

func (that *grpcInterceptor) ServerContext(ctx context.Context) context.Context {
	mtx := mpc.ContextWith(ctx)
	if "" != mtx.GetUrn() {
		return mtx
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	mtx.RewriteContext(&mpc.MeshContext{
		URN:         tool.Anyone(prsim.MeshUrn.GetHeader(md), tool.Anyone(md.Get(":Authority")...)),
		TraceId:     prsim.MeshTraceId.GetHeader(md),
		SpanId:      prsim.MeshSpanId.GetHeader(md),
		Attachments: prsim.GetMetadata(md),
	})
	return mtx
}
