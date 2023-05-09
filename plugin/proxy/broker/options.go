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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type CallOptions struct {
	// Headers is a slice of metadata pointers which should all be set when
	// response header metadata is received.
	Headers []*metadata.MD
	// Trailers is a slice of metadata pointers which should all be set when
	// response trailer metadata is received.
	Trailers []*metadata.MD
	// Peer is a slice of peer pointers which should all be set when the
	// remote peer is known.
	Peer []*peer.Peer
	// Creds are per-RPC credentials to use for a call.
	Creds credentials.PerRPCCredentials
	// MaxRecv is the maximum number of bytes to receive for a single message
	// in a call.
	MaxRecv int
	// MaxSend is the maximum number of bytes to send for a single message in
	// a call.
	MaxSend int
}

func (that *CallOptions) SetHeaders(md metadata.MD) {
	for _, hdr := range that.Headers {
		*hdr = md
	}
}

func (that *CallOptions) SetTrailers(md metadata.MD) {
	for _, tlr := range that.Trailers {
		*tlr = md
	}
}

func (that *CallOptions) SetPeer(p *peer.Peer) {
	for _, pr := range that.Peer {
		*pr = *p
	}
}

func GetCallOptions(opts []grpc.CallOption) *CallOptions {
	var copts CallOptions
	for _, o := range opts {
		switch o := o.(type) {
		case grpc.HeaderCallOption:
			copts.Headers = append(copts.Headers, o.HeaderAddr)
		case grpc.TrailerCallOption:
			copts.Trailers = append(copts.Trailers, o.TrailerAddr)
		case grpc.PeerCallOption:
			copts.Peer = append(copts.Peer, o.PeerAddr)
		case grpc.PerRPCCredsCallOption:
			copts.Creds = o.Creds
		case grpc.MaxRecvMsgSizeCallOption:
			copts.MaxRecv = o.MaxRecvMsgSize
		case grpc.MaxSendMsgSizeCallOption:
			copts.MaxSend = o.MaxSendMsgSize
		}
	}
	return &copts
}

func ApplyPerRPCCreds(ctx context.Context, copts *CallOptions, uri string, isChannelSecure bool) (context.Context, error) {
	if copts.Creds != nil {
		if copts.Creds.RequireTransportSecurity() && !isChannelSecure {
			return nil, fmt.Errorf("transport security is required")
		}
		md, err := copts.Creds.GetRequestMetadata(ctx, uri)
		if nil != err {
			return nil, err
		}
		if len(md) > 0 {
			reqHeaders, ok := metadata.FromOutgoingContext(ctx)
			if ok {
				reqHeaders = metadata.Join(reqHeaders, metadata.New(md))
			} else {
				reqHeaders = metadata.New(md)
			}
			ctx = metadata.NewOutgoingContext(ctx, reqHeaders)
		}
	}
	return ctx, nil
}
