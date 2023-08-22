/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package ptp

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/grpc"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func init() {
	var _ PrivateTransferProtocolServer = privateTransferProtocol
	var _ grpc.RPCService = privateTransferProtocol
	macro.Provide(grpc.IRPCService, privateTransferProtocol)

	var _ PrivateTransferTransportServer = privateTransferTransport
	var _ grpc.RPCService = privateTransferTransport
	macro.Provide(grpc.IRPCService, privateTransferTransport)
}

func WithBound(ctx context.Context, fn func() ([]byte, error)) *Outbound {
	buff, err := fn()
	if nil == err {
		return &Outbound{
			Metadata: map[string]string{},
			Payload:  buff,
			Code:     cause.Success.Code,
			Message:  cause.Success.Message,
		}
	}
	log.Warn(ctx, err.Error())
	code, msg := cause.Parse(err)
	return &Outbound{
		Metadata: map[string]string{},
		Payload:  buff,
		Code:     code,
		Message:  msg,
	}
}

func ServeGRPC[V any](ctx context.Context, fn func() (V, error)) (vx V, ex error) {
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, "%v", err)
			log.Error(ctx, string(debug.Stack()))
			ex = cause.Errorf("%v", err)
		}
	}()
	return fn()
}

var privateTransferProtocol = new(privateTransferProtocolServer)

type privateTransferProtocolServer struct {
}

func (that *privateTransferProtocolServer) Att() *macro.Att {
	return &macro.Att{Name: "ptp-in"}
}

func (that *privateTransferProtocolServer) Metadata() (*grpc2.ServiceDesc, interface{}) {
	return &PrivateTransferProtocol_ServiceDesc, privateTransferProtocol
}

func (that *privateTransferProtocolServer) Transport(server PrivateTransferProtocol_TransportServer) error {
	errs := make(chan error, 1)
	if err := tool.SharedRoutines.Get().Submit(func() {
		ctx := mpc.ContextWith(server.Context())
		for i := 0; ; i++ {
			inbound, err := server.Recv()
			if nil != err {
				errs <- cause.Error(err)
				break
			}
			out, err := privateTransferProtocol.Invoke(ctx, inbound)
			if nil != err {
				errs <- cause.Error(err)
				break
			}
			if err = server.Send(out); nil != err {
				errs <- cause.Error(err)
				break
			}
		}
	}); nil != err {
		return cause.Error(err)
	}
	select {
	case err := <-errs:
		if cause.DeError(err) == io.EOF {
			return nil
		} else {
			return status.Errorf(codes.Internal, "Stream close unexpected, %v.", err)
		}
	}
}

func (that *privateTransferProtocolServer) Invoke(ctx context.Context, inbound *Inbound) (*Outbound, error) {
	return ServeGRPC(ctx, func() (*Outbound, error) {
		if nil == inbound {
			return nil, cause.Validate.Error()
		}
		mtx := that.Context(ctx)
		uri, err := types.FormatURL(prsim.MeshURI.Get(mtx.GetAttachments()))
		if nil != err {
			return nil, cause.Error(err)
		}
		return WithBound(mtx, func() ([]byte, error) {
			switch uri.Path {
			case "/org.ppc.ptp.PrivateTransferTransport/peek":
				pi := new(PeekInbound)
				if err = proto.Unmarshal(inbound.Payload, pi); nil != err {
					return nil, cause.Error(err)
				}
				return peekService.ServePeek(mtx, pi)
			case "/org.ppc.ptp.PrivateTransferTransport/pop":
				pi := new(PopInbound)
				if err = proto.Unmarshal(inbound.Payload, pi); nil != err {
					return nil, cause.Error(err)
				}
				return popService.ServePop(mtx, pi)
			case "/org.ppc.ptp.PrivateTransferTransport/push":
				pi := new(PushInbound)
				if err = proto.Unmarshal(inbound.Payload, pi); nil != err {
					return nil, cause.Error(err)
				}
				return pushService.ServePush(mtx, pi)
			case "/org.ppc.ptp.PrivateTransferTransport/release":
				pi := new(ReleaseInbound)
				if err = proto.Unmarshal(inbound.Payload, pi); nil != err {
					return nil, cause.Error(err)
				}
				return releaseService.ServeRelease(mtx, pi)
			default:
				return nil, cause.NotFound.Error()
			}
		}), nil
	})
}

func (that *privateTransferProtocolServer) Context(ctx context.Context) prsim.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mpc.ContextWith(grpc.Interceptors.ServerContext(ctx))
	}
	return mpc.ContextWith(grpc.Interceptors.ServerContext(metadata.NewOutgoingContext(ctx, md)))
}

var privateTransferTransport = new(privateTransferTransportServer)

type privateTransferTransportServer struct {
}

func (that *privateTransferTransportServer) Att() *macro.Att {
	return &macro.Att{Name: "ptp-out"}
}

func (that *privateTransferTransportServer) Metadata() (*grpc2.ServiceDesc, interface{}) {
	return &PrivateTransferTransport_ServiceDesc, privateTransferTransport
}

func (that *privateTransferTransportServer) Peek(ctx context.Context, inbound *PeekInbound) (*TransportOutbound, error) {
	return ServeGRPC(ctx, func() (*TransportOutbound, error) {
		env, err := aware.Network.GetEnviron(ctx)
		if nil != err {
			return nil, cause.Error(err)
		}
		mtx := mpc.ContextWith(ctx)
		if x := prsim.MeshTargetNodeId.Get(mtx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
			topic := tool.Anyone(prsim.MeshTopic.Get(mtx.GetAttachments()), inbound.Topic)
			b, err := aware.Session.Peek(ctx, topic)
			if nil != err {
				return nil, cause.Error(err)
			}
			return &TransportOutbound{
				Payload: b,
				Code:    cause.Success.Code,
				Message: cause.Success.Message,
			}, nil
		}
		o, err := TransportGRPC(mtx, inbound, fmt.Sprintf("grpcs://ptp.cn%s", "/org.ppc.ptp.PrivateTransferTransport/peek"))
		if nil != err {
			return nil, cause.Error(err)
		}
		return &TransportOutbound{
			Payload: o.Payload,
			Code:    o.Code,
			Message: o.Message,
		}, nil
	})
}

func (that *privateTransferTransportServer) Pop(ctx context.Context, inbound *PopInbound) (*TransportOutbound, error) {
	return ServeGRPC(ctx, func() (*TransportOutbound, error) {
		env, err := aware.Network.GetEnviron(ctx)
		if nil != err {
			return nil, cause.Error(err)
		}
		mtx := mpc.ContextWith(ctx)
		if x := prsim.MeshTargetNodeId.Get(mtx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
			timeout := types.Duration(time.Duration(inbound.Timeout) * time.Millisecond)
			if x := prsim.MeshTimeout.Get(mtx.GetAttachments()); "" != x {
				if t, err := strconv.Atoi(x); nil == err && t > 0 {
					timeout = types.Duration(time.Duration(t) * time.Millisecond)
				}
			}
			topic := tool.Anyone(prsim.MeshTopic.Get(mtx.GetAttachments()), inbound.Topic)
			b, err := aware.Session.Pop(ctx, timeout, topic)
			if nil != err {
				return nil, cause.Error(err)
			}
			return &TransportOutbound{
				Payload: b,
				Code:    cause.Success.Code,
				Message: cause.Success.Message,
			}, nil
		}
		o, err := TransportGRPC(mtx, inbound, fmt.Sprintf("grpcs://ptp.cn%s", "/org.ppc.ptp.PrivateTransferTransport/pop"))
		if nil != err {
			return nil, cause.Error(err)
		}
		return &TransportOutbound{
			Payload: o.Payload,
			Code:    o.Code,
			Message: o.Message,
		}, nil
	})
}

func (that *privateTransferTransportServer) Push(ctx context.Context, inbound *PushInbound) (*TransportOutbound, error) {
	return ServeGRPC(ctx, func() (*TransportOutbound, error) {
		env, err := aware.Network.GetEnviron(ctx)
		if nil != err {
			return nil, cause.Error(err)
		}
		mtx := mpc.ContextWith(ctx)
		if x := prsim.MeshTargetNodeId.Get(mtx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
			topic := tool.Anyone(prsim.MeshTopic.Get(mtx.GetAttachments()), inbound.Topic)
			err = aware.Session.Push(ctx, inbound.Payload, inbound.Metadata, topic)
			if nil != err {
				return nil, cause.Error(err)
			}
			return &TransportOutbound{
				Code:    cause.Success.Code,
				Message: cause.Success.Message,
			}, nil
		}
		o, err := TransportGRPC(mtx, inbound, fmt.Sprintf("grpcs://ptp.cn%s", "/org.ppc.ptp.PrivateTransferTransport/push"))
		if nil != err {
			return nil, cause.Error(err)
		}
		return &TransportOutbound{
			Payload: o.Payload,
			Code:    o.Code,
			Message: o.Message,
		}, nil
	})
}

func (that *privateTransferTransportServer) Release(ctx context.Context, inbound *ReleaseInbound) (*TransportOutbound, error) {
	return ServeGRPC(ctx, func() (*TransportOutbound, error) {
		env, err := aware.Network.GetEnviron(ctx)
		if nil != err {
			return nil, cause.Error(err)
		}
		mtx := mpc.ContextWith(ctx)
		if x := prsim.MeshTargetNodeId.Get(mtx.GetAttachments()); "" == x || strings.EqualFold(env.NodeId, x) {
			timeout := types.Duration(time.Duration(inbound.Timeout) * time.Millisecond)
			if x := prsim.MeshTimeout.Get(mtx.GetAttachments()); "" != x {
				if t, err := strconv.Atoi(x); nil == err && t > 0 {
					timeout = types.Duration(time.Duration(t) * time.Millisecond)
				}
			}
			topic := tool.Anyone(prsim.MeshTopic.Get(mtx.GetAttachments()), inbound.Topic)
			err = aware.Session.Release(ctx, timeout, topic)
			if nil != err {
				return nil, cause.Error(err)
			}
			return &TransportOutbound{
				Code:    cause.Success.Code,
				Message: cause.Success.Message,
			}, nil
		}
		o, err := TransportGRPC(mtx, inbound, fmt.Sprintf("grpcs://ptp.cn%s", "/org.ppc.ptp.PrivateTransferTransport/release"))
		if nil != err {
			return nil, cause.Error(err)
		}
		return &TransportOutbound{
			Payload: o.Payload,
			Code:    o.Code,
			Message: o.Message,
		}, nil
	})
}
