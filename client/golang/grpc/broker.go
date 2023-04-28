/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"context"
	"google.golang.org/grpc"
	"net/http"
)

var IHandler = (*http.Handler)(nil)

var IRPCService = (*RPCService)(nil)

type RPCService interface {
	Metadata() (*grpc.ServiceDesc, interface{})
}

var IChannel = (*Channel)(nil)

type Channel interface {
	DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (Future, context.CancelFunc, error)
}
