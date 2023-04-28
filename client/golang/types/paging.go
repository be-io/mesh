/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"reflect"
)

type Paging struct {
	SID    string                 `index:"0" json:"sid" xml:"sid" yaml:"sid" comment:"Session ID"`
	Index  int64                  `index:"5" json:"index" xml:"index" yaml:"index" comment:"Start index"`
	Limit  int64                  `index:"10" json:"limit" xml:"limit" yaml:"limit" comment:"Start limit"`
	Factor map[string]interface{} `index:"15" json:"factor" xml:"factor" yaml:"factor" comment:"Paging factor"`
	Order  string                 `index:"20" json:"order" xml:"order" yaml:"order" comment:"Sort order"`
}

func (that *Paging) Reset(sid string, index int64, limit int64) *Paging {
	that.SID = sid
	that.Index = index
	that.Limit = limit
	return that
}

func (that *Paging) GetFactor(key string) string {
	if nil == that.Factor || nil == that.Factor[key] || "" == that.Factor[key] {
		return ""
	}
	return fmt.Sprintf("%v", that.Factor[key])
}

func (that *Paging) TryReadFactor(ptr interface{}) error {
	buff, err := codec.Jsonizer.Marshal(that.Factor)
	if nil != err {
		return cause.Error(err)
	}
	err = codec.Jsonizer.Unmarshal(buff, ptr)
	return cause.Error(err)
}

type Page[T any] struct {
	SID   string `index:"0" json:"sid" xml:"sid" yaml:"sid" comment:"Session ID"`
	Index int64  `index:"5" json:"index" xml:"index" yaml:"index" comment:"Start index"`
	Limit int64  `index:"10" json:"limit" xml:"limit" yaml:"limit" comment:"Start limit"`
	Total int64  `index:"15" json:"total" xml:"total" yaml:"total" comment:"Start total"`
	Next  bool   `index:"20" json:"next" xml:"next" yaml:"next" comment:"Has next"`
	Data  []T    `index:"25" json:"data" xml:"data" yaml:"data" comment:"Index data"`
}

func (that *Page[T]) TryRead(ptr any) error {
	if !reflect.ValueOf(ptr).Elem().IsZero() {
		return nil
	}
	buff, err := codec.Jsonizer.Marshal(that.Data)
	if nil != err {
		return cause.Error(err)
	}
	err = codec.Jsonizer.Unmarshal(buff, ptr)
	return cause.Error(err)
}

func Reset[T any](index *Paging, total int64, data []T) *Page[T] {
	return &Page[T]{
		SID:   index.SID,
		Index: index.Index,
		Limit: index.Limit,
		Total: total,
		Next:  total > (index.Index+1)*index.Limit,
		Data:  data,
	}
}
