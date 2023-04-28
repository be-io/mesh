/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package codec

import (
	"bytes"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	var _ Codec = new(Json)
	macro.Provide(ICodec, new(Json))
}

const JSON = "json"

var Jsonizer = jsoniter.ConfigCompatibleWithStandardLibrary

type Json struct {
}

func (that *Json) Att() *macro.Att {
	return &macro.Att{
		Name: JSON,
	}
}

func (that *Json) Encode(value interface{}) (*bytes.Buffer, error) {
	if nil == value {
		return nil, nil
	}
	return Encode(value, func(input interface{}) (*bytes.Buffer, error) {
		if buf, err := Jsonizer.Marshal(value); nil != err {
			return nil, cause.Error(err)
		} else {
			return bytes.NewBuffer(buf), nil
		}
	})
}

func (that *Json) Decode(value *bytes.Buffer, kind interface{}) (interface{}, error) {
	if nil == value {
		return kind, nil
	}
	return Decode(value, kind, func(input *bytes.Buffer, kind any) (any, error) {
		if err := Jsonizer.Unmarshal(value.Bytes(), kind); nil != err {
			return kind, cause.Error(err)
		} else {
			return kind, nil
		}
	})
}

func (that *Json) EncodeString(value interface{}) (string, error) {
	return EncodeString(value, that)
}

func (that *Json) DecodeString(value string, kind interface{}) (interface{}, error) {
	return DecodeString(value, kind, that)
}
