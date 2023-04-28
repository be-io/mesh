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
)

var ICodec = (*Codec)(nil)

type Codec interface {
	macro.SPI

	// Encode nil able if value is nil
	Encode(value any) (*bytes.Buffer, error)

	Decode(value *bytes.Buffer, kind any) (any, error)

	EncodeString(value any) (string, error)

	DecodeString(value string, kind any) (any, error)
}

func EncodeString(value any, codec Codec) (string, error) {
	if x, err := codec.Encode(value); nil != err {
		return "", cause.Error(err)
	} else {
		return string(x.Bytes()), nil
	}
}

func DecodeString(value string, kind any, codec Codec) (any, error) {
	if x, err := codec.Decode(bytes.NewBufferString(value), kind); nil != err {
		return kind, cause.Error(err)
	} else {
		return x, nil
	}
}

func Encode(value any, fn func(input any) (*bytes.Buffer, error)) (*bytes.Buffer, error) {
	if nil == value {
		return nil, nil
	}
	if buf, ok := value.([]byte); ok {
		return bytes.NewBuffer(buf), nil
	}
	if buf, ok := value.(*[]byte); ok {
		return bytes.NewBuffer(*buf), nil
	}
	if buf, ok := value.([]uint8); ok {
		return bytes.NewBuffer(buf), nil
	}
	if buf, ok := value.(*[]uint8); ok {
		return bytes.NewBuffer(*buf), nil
	}
	if buf, ok := value.(*bytes.Buffer); ok {
		if nil == buf {
			return &bytes.Buffer{}, nil
		} else {
			return bytes.NewBuffer(buf.Bytes()), nil
		}
	}
	if buf, ok := value.(bytes.Buffer); ok {
		return bytes.NewBuffer(buf.Bytes()), nil
	}
	return fn(value)
}

func Decode(value *bytes.Buffer, kind any, fn func(input *bytes.Buffer, kind any) (any, error)) (any, error) {
	if nil == value || nil == value.Bytes() {
		return kind, nil
	}
	if _, ok := kind.([]byte); ok {
		return value.Bytes(), nil
	}
	if _, ok := kind.(*[]byte); ok {
		return value.Bytes(), nil
	}
	if _, ok := kind.([]uint8); ok {
		return value.Bytes(), nil
	}
	if _, ok := kind.(*[]uint8); ok {
		return value.Bytes(), nil
	}
	if buf, ok := kind.(*bytes.Buffer); ok {
		buf.Write(value.Bytes())
		return buf, nil
	}
	if buf, ok := kind.(bytes.Buffer); ok {
		buf.Write(value.Bytes())
		return &buf, nil
	}
	return fn(value, kind)
}
