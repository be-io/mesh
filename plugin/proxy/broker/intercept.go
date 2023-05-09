/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type WrappedClientConn interface {
	grpc.ClientConnInterface
	Unwrap() grpc.ClientConnInterface
}

func InterceptClientConn(ch grpc.ClientConnInterface, unaryInt grpc.UnaryClientInterceptor, streamInt grpc.StreamClientInterceptor) grpc.ClientConnInterface {
	if unaryInt == nil && streamInt == nil {
		return ch
	}
	return &interceptedChannel{ch: ch, unaryInt: unaryInt, streamInt: streamInt}
}

type interceptedChannel struct {
	ch        grpc.ClientConnInterface
	unaryInt  grpc.UnaryClientInterceptor
	streamInt grpc.StreamClientInterceptor
}

func (that *interceptedChannel) Unwrap() grpc.ClientConnInterface {
	return that.ch
}

func unwrap(ch grpc.ClientConnInterface) grpc.ClientConnInterface {
	// completely unwrap to find the root ClientConn
	for {
		w, ok := ch.(WrappedClientConn)
		if !ok {
			return ch
		}
		ch = w.Unwrap()
	}
}

func (that *interceptedChannel) Invoke(ctx context.Context, methodName string, req, resp interface{}, opts ...grpc.CallOption) error {
	if that.unaryInt == nil {
		return that.ch.Invoke(ctx, methodName, req, resp, opts...)
	}
	cc, _ := unwrap(that.ch).(*grpc.ClientConn)
	return that.unaryInt(ctx, methodName, req, resp, cc, that.unaryInvoker, opts...)
}

func (that *interceptedChannel) unaryInvoker(ctx context.Context, methodName string, req, resp interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
	return that.ch.Invoke(ctx, methodName, req, resp, opts...)
}

func (that *interceptedChannel) NewStream(ctx context.Context, desc *grpc.StreamDesc, methodName string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if that.streamInt == nil {
		return that.ch.NewStream(ctx, desc, methodName, opts...)
	}
	cc, _ := that.ch.(*grpc.ClientConn)
	return that.streamInt(ctx, desc, cc, methodName, that.streamer, opts...)
}

func (that *interceptedChannel) streamer(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, methodName string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return that.ch.NewStream(ctx, desc, methodName, opts...)
}

var _ grpc.ClientConnInterface = (*interceptedChannel)(nil)

func WithInterceptor(reg grpc.ServiceRegistrar, unaryInt grpc.UnaryServerInterceptor, streamInt grpc.StreamServerInterceptor) grpc.ServiceRegistrar {
	if unaryInt == nil && streamInt == nil {
		return reg
	}
	return &interceptingRegistry{reg: reg, unaryInt: unaryInt, streamInt: streamInt}
}

type interceptingRegistry struct {
	reg       grpc.ServiceRegistrar
	unaryInt  grpc.UnaryServerInterceptor
	streamInt grpc.StreamServerInterceptor
}

func (that *interceptingRegistry) RegisterService(desc *grpc.ServiceDesc, srv interface{}) {
	that.reg.RegisterService(InterceptServer(desc, that.unaryInt, that.streamInt), srv)
}

func InterceptServer(svcDesc *grpc.ServiceDesc, unaryInt grpc.UnaryServerInterceptor, streamInt grpc.StreamServerInterceptor) *grpc.ServiceDesc {
	if unaryInt == nil && streamInt == nil {
		return svcDesc
	}
	intercepted := *svcDesc

	if unaryInt != nil {
		intercepted.Methods = make([]grpc.MethodDesc, len(svcDesc.Methods))
		for i, md := range svcDesc.Methods {
			origHandler := md.Handler
			intercepted.Methods[i] = grpc.MethodDesc{
				MethodName: md.MethodName,
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					combinedInterceptor := unaryInt
					if interceptor != nil {
						// combine unaryInt with the interceptor provided to handler
						combinedInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
							h := func(ctx context.Context, req interface{}) (interface{}, error) {
								return unaryInt(ctx, req, info, handler)
							}
							// we first call provided interceptor, but supply a handler that will call unaryInt
							return interceptor(ctx, req, info, h)
						}
					}
					return origHandler(srv, ctx, dec, combinedInterceptor)
				},
			}
		}
	}

	if streamInt != nil {
		intercepted.Streams = make([]grpc.StreamDesc, len(svcDesc.Streams))
		for i, sd := range svcDesc.Streams {
			origHandler := sd.Handler
			info := &grpc.StreamServerInfo{
				FullMethod:     fmt.Sprintf("/%s/%s", svcDesc.ServiceName, sd.StreamName),
				IsClientStream: sd.ClientStreams,
				IsServerStream: sd.ServerStreams,
			}
			intercepted.Streams[i] = grpc.StreamDesc{
				StreamName:    sd.StreamName,
				ClientStreams: sd.ClientStreams,
				ServerStreams: sd.ServerStreams,
				Handler: func(srv interface{}, stream grpc.ServerStream) error {
					return streamInt(srv, stream, info, origHandler)
				},
			}
		}
	}

	return &intercepted
}
