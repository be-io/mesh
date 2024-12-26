/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proto

import (
	"context"
	"encoding/binary"
	"github.com/opendatav/mesh/client/golang/cause"
	"io"
	"reflect"
	"unsafe"
)

func Register(matcher matcher) {
	repository = append(repository, matcher)
}

func infer(ctx context.Context, kind reflect.Type) (codec, error) {
	cdc, err := repository.match(ctx, kind)
	if nil != err {
		return nil, cause.Error(err)
	}
	if nil != cdc {
		return cdc, nil
	}
	if kind.Kind() == reflect.Ptr {
		cdc, err = repository.match(ctx, kind.Elem())
		if nil != err {
			return nil, cause.Error(err)
		}
		return cdc, nil
	}
	cdc, err = repository.match(ctx, reflect.PtrTo(kind))
	if nil != err {
		return nil, cause.Error(err)
	}
	return cdc, nil
}

func getOriginType(kind reflect.Type) reflect.Type {
	for kind.Kind() == reflect.Ptr {
		kind = kind.Elem()
	}
	return kind
}

const TAG = "index"

var repository matchers

type matcher interface {
	match(ctx context.Context, kind reflect.Type) (codec, error)
}

type matchers []matcher

func (that matchers) match(ctx context.Context, kind reflect.Type) (codec, error) {
	for _, mat := range that {
		cdc, err := mat.match(ctx, kind)
		if nil != err {
			return nil, cause.Error(err)
		}
		if nil != cdc {
			return cdc, nil
		}
	}
	return nil, cause.Errorf("No codec for %s", kind.Name())
}

// codec：
//
//		          C++        Java       Python      Go          C#         PHP            Ruby
//	 double    double     double     float       float64     double     float          Float
//	 float     float      float      float       float32     float      float          Float
//	 int32     int32      int        int         int32       int        integer        Fixnum or Bignum (as required)  使用变长编码，对负数编码效率低，如果你的变量可能是负数，可以使用sint32
//	 int64     int64      long       int/long    int64       long       integer/string Bignum                          使用变长编码，对负数编码效率低，如果你的变量可能是负数，可以使用sint64
//	 uint32    uint32     int        int/long    uint32      uint       integer        Fixnum or Bignum (as required)  使用变长编码
//	 uint64    uint64     long       int/long    uint64      ulong      integer/string Bignum                          使用变长编码
//	 sint32    int32      int        int         int32       int        integer        Fixnum or Bignum (as required)  使用变长编码，带符号的int类型，对负数编码比int32高效
//	 sint64    int64      long       int/long    int64       long       integer/string Bignum                          使用变长编码，带符号的int类型，对负数编码比int64高效
//	 fixed32   uint32     int        int         int32       uint       integer        Fixnum or Bignum (as required)  4字节编码， 如果变量经常大于228
//	 fixed64   uint64     long       int/long    uint64      ulong      integer/string Bignum                          8字节编码， 如果变量经常大于256
//	 sfixed32  int32      int        int         int32       int        integer        Fixnum or Bignum (as required)  4字节编码
//	 sfixed64  int64      long       int/long    int64       long       integer/string Bignum                          8字节编码
//	 bool      bool       boolean    bool        bool        bool       boolean        TrueClass/FalseClass
//	 string    string     String     str/unicode string      string     string         String (UTF-8)                  必须包含utf-8编码或者7-bit ASCII text
//	 bytes     string     ByteString str         []byte      ByteString string         String(ASCII-8BIT)              任意的字节序列
type codec interface {
	wire() wireType
	size(ctx context.Context, pointer unsafe.Pointer, flags flags) int
	encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error)
	decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error)
}

// EncodeTag encodes a pair of field number and wire type into a protobuf tag.
func EncodeTag(f FieldNumber, t WireType) uint64 {
	return uint64(f)<<3 | uint64(t)
}

// EncodeZigZag returns v as a zig-zag encoded value.
func EncodeZigZag(v int64) uint64 {
	return encodeZigZag64(v)
}

func encodeZigZag64(v int64) uint64 {
	return (uint64(v) << 1) ^ uint64(v>>63)
}

func encodeZigZag32(v int32) uint32 {
	return (uint32(v) << 1) ^ uint32(v>>31)
}

func encodeVarint(b []byte, v uint64) (int, error) {
	n := sizeOfVarint(v)

	if len(b) < n {
		return 0, io.ErrShortBuffer
	}

	switch n {
	case 1:
		b[0] = byte(v)

	case 2:
		b[0] = byte(v) | 0x80
		b[1] = byte(v >> 7)

	case 3:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v >> 14)

	case 4:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v>>14) | 0x80
		b[3] = byte(v >> 21)

	case 5:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v>>14) | 0x80
		b[3] = byte(v>>21) | 0x80
		b[4] = byte(v >> 28)

	case 6:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v>>14) | 0x80
		b[3] = byte(v>>21) | 0x80
		b[4] = byte(v>>28) | 0x80
		b[5] = byte(v >> 35)

	case 7:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v>>14) | 0x80
		b[3] = byte(v>>21) | 0x80
		b[4] = byte(v>>28) | 0x80
		b[5] = byte(v>>35) | 0x80
		b[6] = byte(v >> 42)

	case 8:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v>>14) | 0x80
		b[3] = byte(v>>21) | 0x80
		b[4] = byte(v>>28) | 0x80
		b[5] = byte(v>>35) | 0x80
		b[6] = byte(v>>42) | 0x80
		b[7] = byte(v >> 49)

	case 9:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v>>14) | 0x80
		b[3] = byte(v>>21) | 0x80
		b[4] = byte(v>>28) | 0x80
		b[5] = byte(v>>35) | 0x80
		b[6] = byte(v>>42) | 0x80
		b[7] = byte(v>>49) | 0x80
		b[8] = byte(v >> 56)

	case 10:
		b[0] = byte(v) | 0x80
		b[1] = byte(v>>7) | 0x80
		b[2] = byte(v>>14) | 0x80
		b[3] = byte(v>>21) | 0x80
		b[4] = byte(v>>28) | 0x80
		b[5] = byte(v>>35) | 0x80
		b[6] = byte(v>>42) | 0x80
		b[7] = byte(v>>49) | 0x80
		b[8] = byte(v>>56) | 0x80
		b[9] = byte(v >> 63)
	}

	return n, nil
}

func encodeVarintZigZag(b []byte, v int64) (int, error) {
	return encodeVarint(b, encodeZigZag64(v))
}

func encodeLE32(b []byte, v uint32) (int, error) {
	if len(b) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(b, v)
	return 4, nil
}

func encodeLE64(b []byte, v uint64) (int, error) {
	if len(b) < 8 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint64(b, v)
	return 8, nil
}

func encodeTag(b []byte, f int, t wireType) (int, error) {
	return encodeVarint(b, uint64(f)<<3|uint64(t))
}

//

// DecodeTag reverses the encoding applied by EncodeTag.
func DecodeTag(tag uint64) (FieldNumber, WireType) {
	return FieldNumber(tag >> 3), WireType(tag & 7)
}

// DecodeZigZag reverses the encoding applied by EncodeZigZag.
func DecodeZigZag(v uint64) int64 {
	return decodeZigZag64(v)
}

func decodeZigZag64(v uint64) int64 {
	return int64(v>>1) ^ -(int64(v) & 1)
}

func decodeZigZag32(v uint32) int32 {
	return int32(v>>1) ^ -(int32(v) & 1)
}

var (
	errVarintOverflow = cause.Errorf("varint overflowed 64 bits intcodec")
)

func decodeVarint(b []byte) (uint64, int, error) {
	if len(b) != 0 && b[0] < 0x80 {
		// Fast-path for decoding the common case of varints that fit on a
		// single byte.
		//
		// This path is ~60% faster than calling binary.Uvarint.
		return uint64(b[0]), 1, nil
	}

	var x uint64
	var s uint

	for i, c := range b {
		if c < 0x80 {
			if i > 9 || i == 9 && c > 1 {
				return 0, i, errVarintOverflow
			}
			return x | uint64(c)<<s, i + 1, nil
		}
		x |= uint64(c&0x7f) << s
		s += 7
	}

	return x, len(b), io.ErrUnexpectedEOF
}

func decodeVarintZigZag(b []byte) (int64, int, error) {
	v, n, err := decodeVarint(b)
	return decodeZigZag64(v), n, err
}

func decodeLE32(b []byte) (uint32, int, error) {
	if len(b) < 4 {
		return 0, 0, io.ErrUnexpectedEOF
	}
	return binary.LittleEndian.Uint32(b), 4, nil
}

func decodeLE64(b []byte) (uint64, int, error) {
	if len(b) < 8 {
		return 0, 0, io.ErrUnexpectedEOF
	}
	return binary.LittleEndian.Uint64(b), 8, nil
}

func decodeTag(b []byte) (f fieldNumber, t wireType, n int, err error) {
	v, n, err := decodeVarint(b)
	return fieldNumber(v >> 3), wireType(v & 7), n, err
}

func decodeVarlen(b []byte) ([]byte, int, error) {
	v, n, err := decodeVarint(b)
	if nil != err {
		return nil, n, err
	}
	if v > uint64(len(b)-n) {
		return nil, n, io.ErrUnexpectedEOF
	}
	return b[n : n+int(v)], n + int(v), nil
}
