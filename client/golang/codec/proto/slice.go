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
	"io"
	"reflect"
	"unsafe"
)

func init() {
	Register(&slices{slices: map[reflect.Type]codec{}})
}

type repeatedField struct {
}

type slices struct {
	slices map[reflect.Type]codec
}

func (that *slices) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if reflect.Slice != kind.Kind() {
		return nil, nil
	}
	cdc, err := that.construct(ctx, kind)
	if nil != err {
		return nil, cause.Error(err)
	}
	if nil == that.slices[kind] {
		that.slices[kind] = cdc
	}
	return that.slices[kind], nil
}

func (that *slices) construct(ctx context.Context, kind reflect.Type) (codec, error) {
	elem := kind.Elem()

	if elem.Kind() == reflect.Uint8 { // []byte
		return infer(ctx, kind)
	} else {
		cdc, err := infer(ctx, elem)
		if nil != err {
			return nil, cause.Error(err)
		}
		if getOriginType(elem).Kind() == reflect.Struct {
			return &slice{codec: cdc, kind: kind, index: 0, embedded: true}, nil
		}
		return &slice{codec: cdc, kind: kind, index: 0, embedded: false}, nil
	}
}

type slice struct {
	codec    codec
	kind     reflect.Type
	index    int
	embedded bool
}

func (that *slice) wire() wireType {
	return that.codec.wire()
}

func (that *slice) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	elemSize := alignedSize(that.kind.Elem())
	tagSize := sizeOfTag(that.index, that.codec.wire())
	n := 0

	if v := (*Slice)(pointer); v != nil {
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i, elemSize)
			size := that.codec.size(ctx, elem, wantzero)
			n += tagSize + size
			if that.embedded {
				n += sizeOfVarint(uint64(size))
			}
		}
	}

	return n
}

func (that *slice) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	elemSize := alignedSize(that.kind.Elem())
	tagSize := sizeOfTag(that.index, that.codec.wire())
	tagData := make([]byte, tagSize)
	encodeTag(tagData, that.index, that.codec.wire())
	offset := 0

	if s := (*Slice)(pointer); s != nil {
		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i, elemSize)
			size := that.codec.size(ctx, elem, wantzero)

			n := copy(buff[offset:], tagData)
			offset += n
			if n < len(tagData) {
				return offset, io.ErrShortBuffer
			}

			if that.embedded {
				n, err := encodeVarint(buff[offset:], uint64(size))
				offset += n
				if nil != err {
					return offset, err
				}
			}

			if (len(buff) - offset) < size {
				return len(buff), io.ErrShortBuffer
			}

			n, err := that.codec.encode(ctx, buff[offset:offset+size], elem, wantzero)
			offset += n
			if nil != err {
				return offset, err
			}
		}
	}

	return offset, nil
}

func (that *slice) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	elemType := that.kind.Elem()
	elemSize := alignedSize(elemType)
	s := (*Slice)(pointer)
	i := s.Len()

	if i == s.Cap() {
		*s = growSlice(elemType, s)
	}

	n, err := that.codec.decode(ctx, buff, s.Index(i, elemSize), noflags)
	if err == nil {
		s.SetLen(i + 1)
	}
	return n, err
}

func alignedSize(t reflect.Type) uintptr {
	a := t.Align()
	s := t.Size()
	return align(uintptr(a), uintptr(s))
}

func align(align, size uintptr) uintptr {
	if align != 0 && (size%align) != 0 {
		size = ((size / align) + 1) * align
	}
	return size
}

func growSlice(t reflect.Type, s *Slice) Slice {
	capacity := 2 * s.Cap()
	if capacity == 0 {
		capacity = 10
	}
	p := pointer(t)
	d := MakeSlice(p, s.Len(), capacity)
	CopySlice(p, d, *s)
	return d
}
