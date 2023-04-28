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
	Register(new(providers))
}

type providers struct {
}

func (that *providers) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if !implements(kind, customMessageType) || implements(kind, protoMessageType) {
		return nil, nil
	}
	return &provider{kind: kind}, nil
}

type provider struct {
	kind reflect.Type
}

func (that *provider) wire() wireType {
	return varlen
}

func (that *provider) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if m := reflect.NewAt(that.kind, pointer).Interface().(message); m != nil {
			size := m.Size()
			if flags.has(toplevel) {
				return size
			}
			return sizeOfVarLen(size)
		}
	}
	return 0
}

func (that *provider) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if m := reflect.NewAt(that.kind, pointer).Interface().(message); m != nil {
			size := m.Size()

			if flags.has(toplevel) {
				if len(buff) < size {
					return 0, io.ErrShortBuffer
				}
				return m.MarshalTo(buff)
			}

			vlen := sizeOfVarLen(size)
			if len(buff) < vlen {
				return 0, io.ErrShortBuffer
			}

			n1, err := encodeVarint(buff, uint64(size))
			if nil != err {
				return n1, err
			}

			n2, err := m.MarshalTo(buff[n1:])
			return n1 + n2, err
		}
	}
	return 0, nil
}

func (that *provider) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	m := reflect.NewAt(that.kind, pointer).Interface().(message)

	if flags.has(toplevel) {
		return len(buff), m.Unmarshal(buff)
	}

	v, n, err := decodeVarlen(buff)
	if nil != err {
		return n, err
	}

	return n + len(v), m.Unmarshal(v)
}
