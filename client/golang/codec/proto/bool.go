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
	Register(new(boolean))
}

type boolean struct {
}

func (that *boolean) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Bool {
		return that, nil
	}
	return nil, nil
}

func (that *boolean) wire() wireType {
	return varint
}

func (that *boolean) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil && *(*bool)(pointer) || flags.has(wantzero) {
		return 1
	}
	return 0
}

func (that *boolean) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil && *(*bool)(pointer) || flags.has(wantzero) {
		if len(buff) == 0 {
			return 0, io.ErrShortBuffer
		}
		buff[0] = 1
		return 1, nil
	}
	return 0, nil
}

func (that *boolean) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if len(buff) == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	*(*bool)(pointer) = buff[0] != 0
	return 1, nil
}
