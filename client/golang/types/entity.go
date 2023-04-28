/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"bytes"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/schema"
)

type Entity struct {
	Codec  string `index:"0" json:"codec" yaml:"codec" xml:"codec" comment:""`
	Schema string `index:"5" json:"schema" yaml:"schema" xml:"schema" comment:""`
	Buffer []byte `index:"10" json:"buffer" yaml:"buffer" xml:"buffer" comment:""`
}

func (that *Entity) AsEmpty() *Entity {
	return &Entity{Codec: codec.JSON, Schema: ""}
}

func (that *Entity) Wrap(value interface{}) (*Entity, error) {
	cdc, ok := macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	if !ok {
		return nil, cause.NoImplement(codec.JSON)
	}
	buff, err := cdc.Encode(value)
	if nil != err {
		return nil, cause.Error(err)
	}
	that.Codec = codec.JSON
	that.Schema = ""
	that.Buffer = buff.Bytes()
	return that, nil
}

// Present is the entity is present.
func (that *Entity) Present() bool {
	return len(that.Buffer) > 0
}

// ReadObject Try to get object with the type.
func (that *Entity) ReadObject() (interface{}, error) {
	ptr, err := schema.Runtime.Refine(that.Schema)
	if nil != err {
		return nil, cause.Error(err)
	}
	err = that.TryReadObject(ptr)
	return ptr, cause.Error(err)
}

// TryReadObject Try to get object with the type.
func (that *Entity) TryReadObject(ptr interface{}) error {
	if nil == that.Buffer {
		return nil
	}
	if decoder, ok := macro.Load(codec.ICodec).Get(that.Codec).(codec.Codec); ok {
		_, err := decoder.Decode(bytes.NewBuffer(that.Buffer), ptr)
		return cause.Error(err)
	}
	return cause.NoImplement(that.Codec)
}

// TryReadObject Try to get object with the type.
func (that *Entity) String() string {
	return string(that.Buffer)
}
