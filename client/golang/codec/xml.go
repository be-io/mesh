/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package codec

import (
	"bytes"
	"encoding/xml"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
)

func init() {
	macro.Provide(ICodec, new(Xml))
}

const XML = "xml"

type Xml struct {
}

func (that *Xml) Att() *macro.Att {
	return &macro.Att{Name: XML}
}

func (that *Xml) Encode(value interface{}) (*bytes.Buffer, error) {
	if nil == value {
		return nil, nil
	}
	return Encode(value, func(input any) (*bytes.Buffer, error) {
		if buf, err := xml.Marshal(value); nil != err {
			return nil, cause.Error(err)
		} else {
			return bytes.NewBuffer(buf), nil
		}
	})
}

func (that *Xml) Decode(value *bytes.Buffer, kind interface{}) (interface{}, error) {
	if nil == value {
		return kind, nil
	}
	return Decode(value, kind, func(input *bytes.Buffer, kind any) (any, error) {
		if err := xml.Unmarshal(value.Bytes(), kind); nil != err {
			return kind, cause.Error(err)
		} else {
			return kind, nil
		}
	})
}

func (that *Xml) EncodeString(value interface{}) (string, error) {
	return EncodeString(value, that)
}

func (that *Xml) DecodeString(value string, kind interface{}) (interface{}, error) {
	return DecodeString(value, kind, that)
}
