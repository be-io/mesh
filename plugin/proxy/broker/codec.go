/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

import (
	"github.com/be-io/mesh/client/golang/grpc"
	"google.golang.org/grpc/encoding"
)

const (
	UnaryContentTypeV1  = "application/x-protobuf"
	StreamContentTypeV1 = "application/x-mesh-proto+v1"
	ApplicationJson     = "application/json"
)

func getUnaryCodec(contentType string) encoding.Codec {
	// mediaType, _, _ := mime.ParseMediaType(contentType)
	return grpc.Codec
}

func getStreamingCodec(contentType string) encoding.Codec {
	// mediaType, _, _ := mime.ParseMediaType(contentType)
	return grpc.Codec
}
