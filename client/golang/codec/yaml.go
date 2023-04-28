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
	"gopkg.in/yaml.v3"
)

func init() {
	macro.Provide(ICodec, new(Yaml))
}

const YAML = "yaml"

type Yaml struct {
}

func (that *Yaml) Att() *macro.Att {
	return &macro.Att{
		Name: YAML,
	}
}

func (that *Yaml) Encode(value interface{}) (*bytes.Buffer, error) {
	if nil == value {
		return nil, nil
	}
	return Encode(value, func(input any) (*bytes.Buffer, error) {
		if buf, err := yaml.Marshal(value); nil != err {
			return nil, cause.Error(err)
		} else {
			return bytes.NewBuffer(buf), nil
		}
	})
}

func (that *Yaml) Decode(value *bytes.Buffer, kind interface{}) (interface{}, error) {
	if nil == value {
		return kind, nil
	}
	return Decode(value, kind, func(input *bytes.Buffer, kind any) (any, error) {
		if err := yaml.Unmarshal(value.Bytes(), kind); nil != err {
			return kind, cause.Error(err)
		} else {
			return kind, nil
		}
	})
}
func (that *Yaml) EncodeString(value interface{}) (string, error) {
	return EncodeString(value, that)
}

func (that *Yaml) DecodeString(value string, kind interface{}) (interface{}, error) {
	return DecodeString(value, kind, that)
}
