/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func init() {
	var consumer = &grpcConsumer{
		options: []grpc.CallOption{
			grpc.ForceCodec(Codec),
			grpc.MaxCallRecvMsgSize(MaxSize),
			grpc.MaxCallSendMsgSize(MaxSize),
			grpc.UseCompressor(gzip.Name),
		},
		conns: map[string][]Future{},
	}
	var _ mpc.Consumer = consumer
	var _ Futures = consumer
	var _ Channel = consumer
	macro.Provide(mpc.IConsumer, consumer)
	macro.Provide(IChannel, consumer)
}

var StreamDesc = &grpc.StreamDesc{ServerStreams: true, ClientStreams: true}
var ServiceDesc = &grpc.ServiceDesc{
	ServiceName: "mesh-rpc",
	HandlerType: (*StreamService)(nil),
	Metadata:    "",
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "v1",
			Handler:       Service.OnNext,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}

type grpcConsumer struct {
	options []grpc.CallOption
	conns   map[string][]Future
	sync.RWMutex
}

func (that *grpcConsumer) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *grpcConsumer) Start() error {
	return nil
}

// Consume /mesh-rpc/v1
func (that *grpcConsumer) Consume(ctx context.Context, urn string, execution mpc.Execution, inbound *bytes.Buffer) (*bytes.Buffer, error) {
	mtx := mpc.ContextWith(ctx)
	if "" == mtx.GetUrn() {
		mtx.RewriteURN(urn)
	}
	return that.Unary(mtx, urn, inbound)
}

func (that *grpcConsumer) Unary(mtx prsim.Context, urn string, inbound *bytes.Buffer) (*bytes.Buffer, error) {
	address, options, err := that.Options(mtx, urn)
	if nil != err {
		return nil, cause.Error(err)
	}
	timeout, _ := mtx.GetAttribute(mpc.TimeoutKey).(time.Duration)
	ctx, cancel := context.WithTimeout(mpc.WithContext(macro.Context(), mtx), timeout)
	defer cancel()
	conn, err := that.Get(ctx, true, address, options...)
	if nil != err {
		return nil, cause.Error(err)
	}
	defer func() { log.Catch(conn.Close()) }()
	var outbound bytes.Buffer
	if err = conn.Invoke(ctx, "/mesh-rpc/v1", inbound, &outbound, that.options...); nil != err {
		return nil, cause.Error(err)
	}
	return &outbound, nil
}

func (that *grpcConsumer) Stream(mtx prsim.Context, urn string, inbound *bytes.Buffer) (*bytes.Buffer, error) {
	address, options, err := that.Options(mtx, urn)
	if nil != err {
		return nil, cause.Error(err)
	}
	timeout, _ := mtx.GetAttribute(mpc.TimeoutKey).(time.Duration)
	ctx, cancel := context.WithTimeout(mpc.WithContext(macro.Context(), mtx), timeout)
	defer cancel()
	conn, err := that.Get(ctx, false, address, options...)
	if nil != err {
		return nil, cause.Error(err)
	}
	cs, err := conn.NewStream(mtx, StreamDesc, "/mesh-rpc/v1", that.options...)
	if nil != err {
		return nil, cause.Error(err)
	}
	if err = cs.SendMsg(inbound); nil != err {
		return nil, cause.Error(err)
	}
	var outbound bytes.Buffer
	if err = cs.RecvMsg(&outbound); nil != err {
		return nil, cause.Error(err)
	}
	return &outbound, nil
}

func (that *grpcConsumer) Close() error {
	return nil
}

func (that *grpcConsumer) Options(ctx prsim.Context, urn string) (string, []grpc.DialOption, error) {
	address, ok := ctx.GetAttribute(mpc.AddressKey).(string)
	if !ok {
		return "", nil, cause.ValidateErrorf("Invoke address cant be empty.")
	}
	if "" == address {
		address = tool.Address.Get().Any()
	}
	unsafe, ok := ctx.GetAttribute(mpc.InsecureKey).(bool)
	if !ok {
		return "", nil, cause.ValidateErrorf("Invoke insecure mode cant be empty.")
	}
	if address != tool.Address.Get().Any() {
		log.Info(ctx, "%s->%s", urn, address)
	}
	options := []grpc.DialOption{
		grpc.WithAuthority(urn),
		grpc.WithDefaultCallOptions(that.options...),
		grpc.WithStreamInterceptor(Interceptors.ClientStream),
		grpc.WithUnaryInterceptor(Interceptors.ClientUnary),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             15 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	if unsafe {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tc, ok := ctx.GetAttribute(mpc.CertificateKey).(*tls.Config)
		if !ok {
			log.Error(ctx, "Cant parse certification in context. ")
		}
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(tc)))
	}
	if !strings.Contains(address, ":") {
		address = fmt.Sprintf("%s:%d", address, tool.Port)
	}
	return address, options, nil
}

func (that *grpcConsumer) Get(ctx context.Context, once bool, address string, options ...grpc.DialOption) (Future, error) {
	if !once {
		return that.Select(ctx, address, options...)
	}
	conn, err := grpc.DialContext(ctx, address, options...)
	if nil != err {
		return nil, cause.Error(err)
	}
	return &connFuture{future: conn}, nil
}

func (that *grpcConsumer) Select(ctx context.Context, address string, options ...grpc.DialOption) (Future, error) {
	if future := func() Future {
		that.RLock()
		defer that.RUnlock()
		futures := that.conns[address]
		if len(futures) > 0 {
			return futures[rand.Intn(len(futures))]
		}
		return nil
	}(); nil != future {
		return future, nil
	}
	that.Lock()
	defer that.Unlock()
	if futures := that.conns[address]; len(futures) > 0 {
		return futures[rand.Intn(len(futures))], nil
	}
	var futures []Future
	for index := 0; index < 3; index++ {
		conn, err := grpc.DialContext(ctx, address, options...)
		if nil != err {
			return nil, cause.Error(err)
		}
		futures = append(futures, &connFuture{future: conn})
		that.conns[address] = futures
	}
	return futures[rand.Intn(len(futures))], nil
}

func (that *grpcConsumer) DialContext(ctx context.Context, target string, options ...grpc.DialOption) (Future, context.CancelFunc, error) {
	mtx := mpc.ContextWith(ctx)
	_, options, err := that.Options(mtx, mtx.GetUrn())
	if nil != err {
		return nil, nil, cause.Error(err)
	}
	timeout, _ := mtx.GetAttribute(mpc.TimeoutKey).(time.Duration)
	ttx, cancel := context.WithTimeout(mpc.WithContext(macro.Context(), mtx), timeout)
	conn, err := that.Get(ttx, true, target, options...)
	return conn, cancel, err
}
