/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proto

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"reflect"
	"unsafe"
)

func init() {
	Register(&pointers{codecs: map[reflect.Type]codec{}})
}

type pointers struct {
	codecs map[reflect.Type]codec
}

func (that *pointers) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if reflect.Ptr != kind.Kind() {
		return nil, nil
	}
	if nil != that.codecs[kind] {
		return that.codecs[kind], nil
	}
	cdc, err := infer(ctx, kind.Elem())
	if nil != err {
		return nil, cause.Error(err)
	}
	return &ptr{kind: kind, codec: cdc}, nil
}

type ptr struct {
	codec codec
	kind  reflect.Type
}

func (that *ptr) wire() wireType {
	return that.codec.wire()
}

func (that *ptr) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if !flags.has(inline) {
			pointer = *(*unsafe.Pointer)(pointer)
		}
		return that.codec.size(ctx, pointer, flags.without(inline).with(wantzero))
	}
	return 0
}

func (that *ptr) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if !flags.has(inline) {
			pointer = *(*unsafe.Pointer)(pointer)
		}
		return that.codec.encode(ctx, buff, pointer, flags.without(inline).with(wantzero))
	}
	return 0, nil
}

func (that *ptr) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v := (*unsafe.Pointer)(pointer)
	if *v == nil {
		*v = unsafe.Pointer(reflect.New(that.kind.Elem()).Pointer())
	}
	return that.codec.decode(ctx, buff, *v, flags)
}
