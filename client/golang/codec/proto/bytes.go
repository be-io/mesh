/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proto

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"io"
	"reflect"
	"unsafe"
)

func init() {
	Register(__bytes)
}

var __bytes = new(bytes8)

type bytes8 struct {
}

func (that *bytes8) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Slice && kind.Elem().Kind() == reflect.Uint8 {
		return that, nil
	}
	return nil, nil
}

func (that *bytes8) wire() wireType {
	return varlen
}

func (that *bytes8) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*[]byte)(pointer); v != nil || flags.has(wantzero) {
			return sizeOfVarLen(len(v))
		}
	}
	return 0
}

func (that *bytes8) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*[]byte)(pointer); v != nil || flags.has(wantzero) {
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

func (that *bytes8) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeVarlen(buff)
	pb := (*[]byte)(pointer)
	if *pb == nil {
		*pb = make([]byte, 0, len(v))
	}
	*pb = append((*pb)[:0], v...)
	return n, err
}

func makeBytes(p unsafe.Pointer, n int) []byte {
	return *(*[]byte)(unsafe.Pointer(&sliceHeader{
		Data: p,
		Len:  n,
		Cap:  n,
	}))
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// isZeroBytes is an optimized version of this loop:
//
//	for i := range b {
//		if b[i] != 0 {
//			return false
//		}
//	}
//	return true
//
// This implementation significantly reduces the CPU footprint of checking for
// slices to be zero, especially when the length increases (these cases should
// be rare tho).
//
// name            old time/op  new time/op  delta
// IsZeroBytes0    1.78ns ± 1%  2.29ns ± 4%  +28.65%  (p=0.000 n=8+10)
// IsZeroBytes4    3.17ns ± 3%  2.37ns ± 3%  -25.21%  (p=0.000 n=10+10)
// IsZeroBytes7    3.97ns ± 4%  3.26ns ± 3%  -18.02%  (p=0.000 n=10+10)
// IsZeroBytes64K  14.8µs ± 3%   1.9µs ± 3%  -87.34%  (p=0.000 n=10+10)
func isZeroBytes(b []byte) bool {
	if n := len(b) / 8; n != 0 {
		if !isZeroUint64(*(*[]uint64)(unsafe.Pointer(&sliceHeader{
			Data: unsafe.Pointer(&b[0]),
			Len:  n,
			Cap:  n,
		}))) {
			return false
		}
		b = b[n*8:]
	}
	switch len(b) {
	case 7:
		return bto32(b) == 0 && bto16(b[4:]) == 0 && b[6] == 0
	case 6:
		return bto32(b) == 0 && bto16(b[4:]) == 0
	case 5:
		return bto32(b) == 0 && b[4] == 0
	case 4:
		return bto32(b) == 0
	case 3:
		return bto16(b) == 0 && b[2] == 0
	case 2:
		return bto16(b) == 0
	case 1:
		return b[0] == 0
	default:
		return true
	}
}

func bto32(b []byte) uint32 {
	return *(*uint32)(unsafe.Pointer(&b[0]))
}

func bto16(b []byte) uint16 {
	return *(*uint16)(unsafe.Pointer(&b[0]))
}

func isZeroUint64(b []uint64) bool {
	for i := range b {
		if b[i] != 0 {
			return false
		}
	}
	return true
}

func init() {
	Register(&bytearrays{arrays: map[reflect.Type]codec{}})
}

type bytearrays struct {
	arrays map[reflect.Type]codec
}

func (that *bytearrays) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() != reflect.Array || kind.Elem().Kind() != reflect.Uint8 {
		return nil, nil
	}
	if nil != that.arrays[kind] {
		that.arrays[kind] = &bytearray{length: kind.Len()}
	}
	return that.arrays[kind], nil
}

type bytearray struct {
	length int
}

func (that *bytearray) match(ctx context.Context, kind reflect.Type) bool {
	if kind.Kind() != reflect.Array || kind.Elem().Kind() != reflect.Uint8 {
		return false
	}
	return true
}

func (that *bytearray) wire() wireType {
	return varlen
}

func (that *bytearray) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	size := sizeOfVarLen(that.length)
	if pointer != nil && (flags.has(wantzero) || !isZeroBytes(makeBytes(pointer, that.length))) {
		return size
	}
	return 0
}

func (that *bytearray) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := makeBytes(pointer, that.length); flags.has(wantzero) || !isZeroBytes(v) {
			return __bytes.encode(ctx, buff, unsafe.Pointer(&v), noflags)
		}
	}
	return 0, nil
}

func (that *bytearray) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, r, err := decodeVarlen(buff)
	if err == nil {
		if copy(makeBytes(pointer, that.length), v) != that.length {
			err = cause.Errorf("cannot decode byte sequence of size %d into byte array of size %d", len(v), that.length)
		}
	}
	return r, err
}
