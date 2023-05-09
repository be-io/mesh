/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"bytes"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/grpc"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/plugin/proxy/broker"
	"golang.org/x/net/context"
	"golang.org/x/net/http/httpguts"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

func init() {
	h21broker := new(h21)
	var _ grpc.StreamService = h21broker
	var _ grpc.RPCService = h21broker
	macro.Provide(grpc.IRPCService, h21broker)

	server := broker.NewServer(
		broker.WithServerUnaryInterceptor(grpc.Interceptors.ServerUnary),
		broker.WithServerStreamInterceptor(grpc.Interceptors.ServerStream),
	)
	h12broker := &h12{broker: server}
	var _ grpc.StreamService = h12broker
	var _ grpc.RPCService = h12broker
	var _ http.Handler = h12broker
	macro.Provide(grpc.IRPCService, h12broker)
	macro.Provide(grpc.IHandler, h12broker)

	server.RegisterService(h12broker.Metadata())
}

const hhh = "x-mo-"

type h21 struct {
}

func (that *h21) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.bridge.h21"}
}

func (that *h21) Metadata() (*rpc.ServiceDesc, interface{}) {
	return &rpc.ServiceDesc{
		ServiceName: "mesh-h21",
		HandlerType: (*grpc.StreamService)(nil),
		Metadata:    "",
		Methods:     []rpc.MethodDesc{},
		Streams: []rpc.StreamDesc{
			{
				StreamName:    "v1",
				Handler:       that.OnNext,
				ServerStreams: true,
				ClientStreams: true,
			},
		},
	}, that
}

func (that *h21) Context(ctx context.Context) prsim.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mpc.ContextWith(grpc.Interceptors.ServerContext(ctx))
	}
	hd := metadata.MD{}
	for k, v := range md {
		if strings.EqualFold("content-length", k) {
			continue
		}
		if httpguts.ValidHeaderFieldName(k) && !prsim.IsMeshMetadata(k) {
			hd.Set(fmt.Sprintf("%s%s", hhh, k), v...)
		}
	}
	return mpc.ContextWith(grpc.Interceptors.ServerContext(metadata.NewOutgoingContext(ctx, hd)))
}

func (that *h21) OnNext(srv interface{}, ss rpc.ServerStream) error {
	ctx := that.Context(ss.Context())
	err := that.Roundtrip(ctx, srv, ss)
	if nil != err {
		log.Debug(ctx, "%s, %s", ctx.GetUrn(), err.Error())
	}
	return cause.DeError(err)
}

func (that *h21) Roundtrip(ctx prsim.Context, srv interface{}, ss rpc.ServerStream) (ex error) {
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, "%v", err)
			log.Error(ctx, string(debug.Stack()))
			ex = cause.Errorf("%v", err)
		}
	}()
	conn, err := broker.DialContext(ctx, fmt.Sprintf("%s://127.0.0.1:570", "http"))
	if nil != err {
		return cause.Error(err)
	}

	cs, err := conn.NewStream(ctx, grpc.StreamDesc, "/mesh-h12/v1")
	if nil != err {
		return cause.Error(err)
	}
	p := &pip{timestamp: time.Now(), ss: ss, cs: cs, sce: make(chan error, 1), cse: make(chan error, 1)}
	// Stream must iterate read send in go routine
	return p.Transport(ctx)
}

type h12 struct {
	broker http.Handler
}

func (that *h12) Att() *macro.Att {
	return &macro.Att{
		Name:    "mesh.plugin.proxy.bridge.h12",
		Pattern: "/mesh-h12/v1",
	}
}

func (that *h12) Metadata() (*rpc.ServiceDesc, interface{}) {
	return &rpc.ServiceDesc{
		ServiceName: "mesh-h12",
		HandlerType: (*grpc.StreamService)(nil),
		Metadata:    "",
		Methods:     []rpc.MethodDesc{},
		Streams: []rpc.StreamDesc{
			{
				StreamName:    "v1",
				Handler:       that.OnNext,
				ServerStreams: true,
				ClientStreams: true,
			},
		},
	}, that
}

func (that *h12) Context(ctx context.Context) prsim.Context {
	mtx := mpc.ContextWith(grpc.Interceptors.ServerContext(ctx))
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mtx
	}
	hs := map[string]string{}
	for k, v := range md {
		if !httpguts.ValidHeaderFieldName(k) || prsim.IsMeshMetadata(k) || strings.EqualFold("content-length", k) {
			continue
		}
		if len(v) > 0 {
			hs[strings.TrimPrefix(k, hhh)] = v[0]
		}
	}
	mtx.SetAttribute(mpc.HeaderKey, hs)
	return mtx
}

func (that *h12) OnNext(srv interface{}, ss rpc.ServerStream) error {
	ctx := that.Context(ss.Context())
	err := that.Roundtrip(ctx, srv, ss)
	if nil != err {
		log.Debug(ctx, "%s, %s", ctx.GetUrn(), err.Error())
	}
	return cause.DeError(err)
}

func (that *h12) Roundtrip(ctx prsim.Context, srv interface{}, ss rpc.ServerStream) (ex error) {
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, "%v", err)
			log.Error(ctx, string(debug.Stack()))
			ex = cause.Errorf("%v", err)
		}
	}()
	conn, cancel, err := aware.Channel.DialContext(ctx, "127.0.0.1:570")
	if nil != err {
		return cause.Error(err)
	}
	defer func() {
		cancel()
		log.Catch(conn.Close())
	}()
	cs, err := conn.NewStream(ctx, grpc.StreamDesc, tool.Anyone(prsim.MeshPath.Get(ctx.GetAttachments()), "/mesh-rpc/v1"))
	if nil != err {
		return cause.Error(err)
	}
	p := &pip{timestamp: time.Now(), ss: ss, cs: cs, sce: make(chan error, 1), cse: make(chan error, 1)}
	// Stream must iterate read send in go routine
	return p.Transport(ctx)
}

func (that *h12) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	that.broker.ServeHTTP(writer, request)
}

type pip struct {
	timestamp time.Time
	ss        rpc.ServerStream
	cs        rpc.ClientStream
	sce       chan error
	cse       chan error
}

func (that *pip) Transport(ctx prsim.Context) error {
	// https://groups.google.com/forum/#!msg/golang-nuts/pZwdYRGxCIk/qpbHxRRPJdUJ
	if err := tool.SharedRoutines.Get().Submit(func() {
		input := &bytes.Buffer{}
		for {
			input.Reset()
			if err := that.ss.RecvMsg(input); nil != err {
				that.sce <- cause.Error(err)
				break
			}
			if err := that.cs.SendMsg(input); nil != err {
				that.sce <- cause.Error(err)
				break
			}
		}
	}); nil != err {
		return cause.Error(err)
	}
	if err := tool.SharedRoutines.Get().Submit(func() {
		// client to server headers are only readable after first client msg is received but must be written to server stream
		// before the first msg is flushed. This is the only place to do it nicely.
		input := &bytes.Buffer{}
		for i := 0; ; i++ {
			input.Reset()
			if err := that.cs.RecvMsg(input); nil != err {
				that.cse <- cause.Error(err)
				break
			}
			if i == 0 {
				md, err := that.cs.Header()
				if nil != err {
					that.cse <- cause.Error(err)
					break
				}
				if err = that.ss.SendHeader(md); nil != err {
					that.cse <- cause.Error(err)
					break
				}
			}
			if err := that.ss.SendMsg(input); nil != err {
				that.cse <- cause.Error(err)
				break
			}
		}
	}); nil != err {
		return cause.Error(err)
	}
	for i := 0; i < 2; i++ {
		select {
		case err := <-that.sce:
			if cause.DeError(err) == io.EOF {
				log.Catch(that.cs.CloseSend())
			} else {
				return status.Errorf(codes.Internal, "Stream close unexpected, %v.", err)
			}
		case err := <-that.cse:
			that.ss.SetTrailer(that.cs.Trailer())
			if cause.DeError(err) != io.EOF {
				return err
			}
			return nil
		}
	}
	return status.Errorf(codes.Internal, "Stream close unexpected.")
}
