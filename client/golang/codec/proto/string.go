/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proto

import (
	"context"
	"io"
	"reflect"
	"unsafe"
)

func init() {
	Register(new(str))
}

type str struct {
}

func (that *str) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.String {
		return that, nil
	}
	return nil, nil
}
func (that *str) wire() wireType {
	return varlen
}

func (that *str) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*string)(pointer); v != "" || flags.has(wantzero) {
			return sizeOfVarLen(len(v))
		}
	}
	return 0
}

func (that *str) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*string)(pointer); v != "" || flags.has(wantzero) {
			n, err := encodeVarint(buff, uint64(len(v)))
			if nil != err {
				return n, err
			}
			c := copy(buff[n:], v)
			n += c
			if c < len(v) {
				err = io.ErrShortBuffer
			}
			return n, err
		}
	}
	return 0, nil
}

func (that *str) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeVarlen(buff)
	*(*string)(pointer) = string(v)
	return n, err
}
