/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
//go:generate go run ../../../client/golang/proto/generate.go -m github.com/be-io/mesh/plugin/proxy/broker
