/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"bytes"
	"github.com/be-io/mesh/client/golang/cause"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

func init() {
	var _ encoding.Codec = new(Marshaller)
}

var Codec = &Marshaller{Default: &protoMarshaller{}}

type Frame struct {
	payload []byte
}

type Marshaller struct {
	Default encoding.Codec
}

func (that *Marshaller) Marshal(v interface{}) ([]byte, error) {
	if nil == v {
		return nil, nil
	}
	if x, ok := v.([]byte); ok {
		return x, nil
	}
	if x, ok := v.(*[]byte); ok {
		return *x, nil
	}
	if x, ok := v.([]uint8); ok {
		return x, nil
	}
	if x, ok := v.(*[]uint8); ok {
		return *x, nil
	}
	if x, ok := v.(*Frame); ok {
		return x.payload, nil
	}
	if x, ok := v.(bytes.Buffer); ok {
		return x.Bytes(), nil
	}
	if x, ok := v.(*bytes.Buffer); ok {
		return x.Bytes(), nil
	}
	if nil != that.Default {
		return that.Default.Marshal(v)
	}
	return nil, cause.Errorf("Cant not serialize input arguments.")
}

func (that *Marshaller) Unmarshal(data []byte, v interface{}) error {
	if nil == data {
		return nil
	}
	if x, ok := v.([]byte); ok {
		copy(x, data)
		return nil
	}
	if x, ok := v.(*[]byte); ok {
		copy(*x, data)
		return nil
	}
	if x, ok := v.([]uint8); ok {
		copy(x, data)
		return nil
	}
	if x, ok := v.(*[]uint8); ok {
		copy(*x, data)
		return nil
	}
	if x, ok := v.(*Frame); ok {
		copy(x.payload, data)
		return nil
	}
	if x, ok := v.(bytes.Buffer); ok {
		x.Write(data)
		return nil
	}
	if x, ok := v.(*bytes.Buffer); ok {
		x.Write(data)
		return nil
	}
	if nil != that.Default {
		return that.Default.Unmarshal(data, v)
	}
	v = data
	return nil
}

func (that *Marshaller) Name() string {
	return "mesh"
}

// protoMarshaller is a Codec implementation with protobuf. It is the default rawCodec for gRPC.
type protoMarshaller struct{}

func (*protoMarshaller) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (*protoMarshaller) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (*protoMarshaller) Name() string {
	return "proto"
}
