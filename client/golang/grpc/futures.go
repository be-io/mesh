/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
)

func init() {
	var _ Future = new(connFuture)
	var _ grpc.ClientStream = new(streamFuture)
}

var IFutures = (*Futures)(nil)

type Future interface {
	io.Closer
	grpc.ClientConnInterface
}

type Futures interface {
	Get(ctx context.Context, once bool, address string, options ...grpc.DialOption) (Future, error)
}

type connFuture struct {
	future *grpc.ClientConn
}

func (that *connFuture) Close() error {
	return expects(that.future.Close())
}

func (that *connFuture) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return expects(that.future.Invoke(ctx, method, args, reply, opts...))
}

func (that *connFuture) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	stream, err := that.future.NewStream(ctx, desc, method, opts...)
	return &streamFuture{stream: stream}, expects(err)
}

type streamFuture struct {
	stream grpc.ClientStream
}

func (that *streamFuture) Header() (metadata.MD, error) {
	md, err := that.stream.Header()
	return md, expects(err)
}

func (that *streamFuture) Trailer() metadata.MD {
	return that.stream.Trailer()
}

func (that *streamFuture) CloseSend() error {
	return expects(that.stream.CloseSend())
}

func (that *streamFuture) Context() context.Context {
	return that.stream.Context()
}

func (that *streamFuture) SendMsg(m interface{}) error {
	return expects(that.stream.SendMsg(m))
}

func (that *streamFuture) RecvMsg(m interface{}) error {
	return expects(that.stream.RecvMsg(m))
}

func expects(exception error) error {
	if nil == exception {
		return nil
	}
	err, ok := status.FromError(exception)
	if !ok {
		return exception
	}
	switch err.Code() {
	case codes.OK:
		break
	case codes.Aborted:
	case codes.Canceled:
	case codes.Unknown:
	case codes.Internal:
	case codes.AlreadyExists:
		return cause.Error(exception)
	case codes.InvalidArgument:
	case codes.DataLoss:
		return cause.ValidateError(exception)
	case codes.DeadlineExceeded:
		return cause.TimeoutError(exception)
	case codes.Unavailable:
		return cause.Errorc(cause.NetUnavailable, exception)
	case codes.NotFound:
	case codes.Unimplemented:
		return cause.NotFoundError(exception)
	case codes.OutOfRange:
	case codes.ResourceExhausted:
	case codes.Unauthenticated:
	case codes.PermissionDenied:
	case codes.FailedPrecondition:
		return cause.UnauthorizedError(exception)
	}
	return exception
}
