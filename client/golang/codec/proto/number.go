/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proto

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

func init() {
	Register(_int)
	Register(_uint)
	Register(_uint64)
	Register(_fixed64)
	Register(_uint32)
	Register(_fixed32)
	Register(_int32)
	Register(_int64)
	Register(_float32)
	Register(_float64)
}

var (
	_int     = new(integer)
	_uint    = new(uinteger)
	_uint64  = new(uint64bit)
	_fixed64 = new(fixed64bit)
	_uint32  = new(uint32bit)
	_fixed32 = new(fixed32bit)
	_int32   = new(int32bit)
	_int64   = new(int64bit)
	_float32 = new(float32bit)
	_float64 = new(float64bit)
)

type integer struct {
}

func (that *integer) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Int {
		return that, nil
	}
	return nil, nil
}

func (that *integer) wire() wireType {
	return varint
}

func (that *integer) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*int)(pointer); v != 0 || flags.has(wantzero) {
			return sizeOfVarint(flags.uint64(int64(v)))
		}
	}
	return 0
}

func (that *integer) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*int)(pointer); v != 0 || flags.has(wantzero) {
			return encodeVarint(buff, flags.uint64(int64(v)))
		}
	}
	return 0, nil
}

func (that *integer) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeVarint(buff)
	*(*int)(pointer) = int(flags.int64(v))
	return n, err
}

type uinteger struct {
}

func (that *uinteger) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Uint {
		return that, nil
	}
	return nil, nil
}

func (that *uinteger) wire() wireType {
	return varint
}

func (that *uinteger) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*uint)(pointer); v != 0 || flags.has(wantzero) {
			return sizeOfVarint(uint64(v))
		}
	}
	return 0
}

func (that *uinteger) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*uint)(pointer); v != 0 || flags.has(wantzero) {
			return encodeVarint(buff, uint64(v))
		}
	}
	return 0, nil
}

func (that *uinteger) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeVarint(buff)
	*(*uint)(pointer) = uint(v)
	return n, err
}

type uint64bit struct {
}

func (that *uint64bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Uint64 {
		return that, nil
	}
	return nil, nil
}

func (that *uint64bit) wire() wireType {
	return varint
}

func (that *uint64bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*uint64)(pointer); v != 0 || flags.has(wantzero) {
			return sizeOfVarint(v)
		}
	}
	return 0
}

func (that *uint64bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*uint64)(pointer); v != 0 || flags.has(wantzero) {
			return encodeVarint(buff, v)
		}
	}
	return 0, nil
}

func (that *uint64bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeVarint(buff)
	*(*uint64)(pointer) = uint64(v)
	return n, err
}

type fixed64bit struct {
}

func (that *fixed64bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	return nil, nil
}

func (that *fixed64bit) wire() wireType {
	return fixed64
}

func (that *fixed64bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*uint64)(pointer); v != 0 || flags.has(wantzero) {
			return 8
		}
	}
	return 0
}

func (that *fixed64bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*uint64)(pointer); v != 0 || flags.has(wantzero) {
			return encodeLE64(buff, v)
		}
	}
	return 0, nil
}

func (that *fixed64bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeLE64(buff)
	*(*uint64)(pointer) = v
	return n, err
}

type uint32bit struct {
}

func (that *uint32bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Uint32 {
		return that, nil
	}
	return nil, nil
}

func (that *uint32bit) wire() wireType {
	return varint
}

func (that *uint32bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*uint32)(pointer); v != 0 || flags.has(wantzero) {
			return sizeOfVarint(uint64(v))
		}
	}
	return 0
}

func (that *uint32bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*uint32)(pointer); v != 0 || flags.has(wantzero) {
			return encodeVarint(buff, uint64(v))
		}
	}
	return 0, nil
}

func (that *uint32bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeVarint(buff)
	if v > math.MaxUint32 {
		return n, fmt.Errorf("intcodec overflow decoding %v into uint32", v)
	}
	*(*uint32)(pointer) = uint32(v)
	return n, err
}

type fixed32bit struct {
}

func (that *fixed32bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	return nil, nil
}

func (that *fixed32bit) wire() wireType {
	return fixed32
}

func (that *fixed32bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*uint32)(pointer); v != 0 || flags.has(wantzero) {
			return 4
		}
	}
	return 0
}

func (that *fixed32bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*uint32)(pointer); v != 0 || flags.has(wantzero) {
			return encodeLE32(buff, v)
		}
	}
	return 0, nil
}

func (that *fixed32bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeLE32(buff)
	*(*uint32)(pointer) = v
	return n, err
}

type int32bit struct {
}

func (that *int32bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Int32 {
		return that, nil
	}
	return nil, nil
}

func (that *int32bit) wire() wireType {
	return varint
}

func (that *int32bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*int32)(pointer); v != 0 || flags.has(wantzero) {
			return sizeOfVarint(flags.uint64(int64(v)))
		}
	}
	return 0
}

func (that *int32bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*int32)(pointer); v != 0 || flags.has(wantzero) {
			return encodeVarint(buff, flags.uint64(int64(v)))
		}
	}
	return 0, nil
}

func (that *int32bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	u, n, err := decodeVarint(buff)
	v := flags.int64(u)
	if v < math.MinInt32 || v > math.MaxInt32 {
		return n, fmt.Errorf("intcodec overflow decoding %v into int32", v)
	}
	*(*int32)(pointer) = int32(v)
	return n, err
}

type int64bit struct {
}

func (that *int64bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Int64 {
		return that, nil
	}
	return nil, nil
}

func (that *int64bit) wire() wireType {
	return varint
}

func (that *int64bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*int64)(pointer); v != 0 || flags.has(wantzero) {
			return sizeOfVarint(flags.uint64(v))
		}
	}
	return 0
}

func (that *int64bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*int64)(pointer); v != 0 || flags.has(wantzero) {
			return encodeVarint(buff, flags.uint64(v))
		}
	}
	return 0, nil
}

func (that *int64bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeVarint(buff)
	*(*int64)(pointer) = flags.int64(v)
	return n, err
}

type float32bit struct {
}

func (that *float32bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Float32 {
		return that, nil
	}
	return nil, nil
}

func (that *float32bit) wire() wireType {
	return fixed32
}

func (that *float32bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*float32)(pointer); v != 0 || flags.has(wantzero) || math.Signbit(float64(v)) {
			return 4
		}
	}
	return 0
}

func (that *float32bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*float32)(pointer); v != 0 || flags.has(wantzero) || math.Signbit(float64(v)) {
			return encodeLE32(buff, math.Float32bits(v))
		}
	}
	return 0, nil
}

func (that *float32bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeLE32(buff)
	*(*float32)(pointer) = math.Float32frombits(v)
	return n, err
}

type float64bit struct {
}

func (that *float64bit) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() == reflect.Float64 {
		return that, nil
	}
	return nil, nil
}

func (that *float64bit) wire() wireType {
	return fixed64
}

func (that *float64bit) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	if pointer != nil {
		if v := *(*float64)(pointer); v != 0 || flags.has(wantzero) || math.Signbit(v) {
			return 8
		}
	}
	return 0
}

func (that *float64bit) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	if pointer != nil {
		if v := *(*float64)(pointer); v != 0 || flags.has(wantzero) || math.Signbit(v) {
			return encodeLE64(buff, math.Float64bits(v))
		}
	}
	return 0, nil
}

func (that *float64bit) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	v, n, err := decodeLE64(buff)
	*(*float64)(pointer) = math.Float64frombits(v)
	return n, err
}
