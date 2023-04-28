/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package codec

import (
	"bytes"
	"github.com/be-io/mesh/client/golang/codec/proto"
	"github.com/be-io/mesh/client/golang/macro"
)

func init() {
	macro.Provide(ICodec, new(Protobuf))
}

const PROTOBUF = "protobuf"

type Protobuf struct {
}

func (that *Protobuf) Att() *macro.Att {
	return &macro.Att{Name: PROTOBUF}
}

func (that *Protobuf) Encode(value interface{}) (*bytes.Buffer, error) {
	if nil == value {
		return nil, nil
	}
	return Encode(value, func(input any) (*bytes.Buffer, error) {
		buff, err := proto.Marshal(macro.Context(), value)
		return bytes.NewBuffer(buff), err
	})
}

func (that *Protobuf) Decode(value *bytes.Buffer, kind interface{}) (interface{}, error) {
	return Decode(value, kind, func(input *bytes.Buffer, kind any) (any, error) {
		err := proto.Unmarshal(macro.Context(), value.Bytes(), value)
		return kind, err
	})
}

func (that *Protobuf) EncodeString(value interface{}) (string, error) {
	return EncodeString(value, that)
}

func (that *Protobuf) DecodeString(value string, kind interface{}) (interface{}, error) {
	return DecodeString(value, kind, that)
}
