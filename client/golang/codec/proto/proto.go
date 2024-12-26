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
	"google.golang.org/protobuf/proto"
	"math/bits"
	"reflect"
	"unsafe"
)

func Size(ctx context.Context, v interface{}) (int, error) {
	kind, p := inspect(v)
	cdc, err := infer(ctx, kind)
	if nil != err {
		return 0, cause.Error(err)
	}
	return cdc.size(ctx, p, inline|toplevel), nil
}

func Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	if m, ok := v.(proto.Message); ok {
		return proto.Marshal(m)
	}
	kind, p := inspect(v)
	cdc, err := infer(ctx, kind)
	if nil != err {
		return nil, cause.Error(err)
	}
	buff := make([]byte, cdc.size(ctx, p, inline|toplevel))
	_, err = cdc.encode(ctx, buff, p, inline|toplevel)
	if nil != err {
		return nil, cause.Error(err)
	}
	return buff, nil
}

func Unmarshal(ctx context.Context, buff []byte, ptr interface{}) error {
	kind, pt := inspect(ptr)
	if kind.Implements(pureMessageType) {
		return proto.Unmarshal(buff, ptr.(proto.Message))
	}
	if len(buff) == 0 {
		// An empty input is a valid protobuf message with all fields set to the
		// zero-value.
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
		return nil
	}
	kind = kind.Elem() // Unmarshal must be passed a pointer
	cdc, err := infer(ctx, kind)
	if nil != err {
		return cause.Error(err)
	}

	n, err := cdc.decode(ctx, buff, pt, toplevel)
	if nil != err {
		return err
	}
	if n < len(buff) {
		return cause.Errorf("proto.Unmarshal(%T): read=%d < buffer=%d", ptr, n, len(buff))
	}
	return nil
}

type flags uintptr

const (
	noflags  flags = 0
	inline   flags = 1 << 0
	wantzero flags = 1 << 1
	// Shared with structField.flags in struct.go:
	// zigzag flags = 1 << 2
	toplevel flags = 1 << 3
)

func (f flags) has(x flags) bool {
	return (f & x) != 0
}

func (f flags) with(x flags) flags {
	return f | x
}

func (f flags) without(x flags) flags {
	return f & ^x
}

func (f flags) uint64(i int64) uint64 {
	if f.has(zigzag) {
		return encodeZigZag64(i)
	} else {
		return uint64(i)
	}
}

func (f flags) int64(u uint64) int64 {
	if f.has(zigzag) {
		return decodeZigZag64(u)
	} else {
		return int64(u)
	}
}

type iface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func inspect(v interface{}) (reflect.Type, unsafe.Pointer) {
	return reflect.TypeOf(v), pointer(v)
}

func pointer(v interface{}) unsafe.Pointer {
	return (*iface)(unsafe.Pointer(&v)).ptr
}

func inlined(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Ptr:
		return true
	case reflect.Map:
		return true
	case reflect.Struct:
		return t.NumField() == 1 && inlined(t.Field(0).Type)
	default:
		return false
	}
}

type fieldNumber uint

// wireType
//
//	字段类型                                                                    二进制类型            二进制编码值
//	int32, int64, uint32, uint64, sint32, sint64, bool, enum                   Varint              0
//	fixed64, sfixed64, double                                                  64bit               1
//	string, bytes, embedded messages, packed repeated fields                   Length-delimited    2
//	groups(deprecated)                                                         Start group         3
//	groups(deprecated)                                                         End group           4
//	fixed32, sfixed32, float                                                   32bit               5
type wireType uint

const (
	varint  wireType = 0
	fixed64 wireType = 1
	varlen  wireType = 2
	fixed32 wireType = 5
)

func (wt wireType) String() string {
	switch wt {
	case varint:
		return "varint"
	case varlen:
		return "varlen"
	case fixed32:
		return "fixed32"
	case fixed64:
		return "fixed64"
	default:
		return "unknown"
	}
}

// backward compatibility with gogoproto custom types.
type message interface {
	Size() int
	MarshalTo([]byte) (int, error)
	Unmarshal([]byte) error
}

type protoMessage interface {
	ProtoMessage()
}

var (
	messageType       = reflect.TypeOf((*Message)(nil)).Elem()
	customMessageType = reflect.TypeOf((*message)(nil)).Elem()
	protoMessageType  = reflect.TypeOf((*protoMessage)(nil)).Elem()
	pureMessageType   = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

func implements(t, iface reflect.Type) bool {
	return t.Implements(iface) || reflect.PtrTo(t).Implements(iface)
}

func sizeOfVarint(v uint64) int {
	return (bits.Len64(v|1) + 6) / 7
}

func sizeOfVarintZigZag(v int64) int {
	return sizeOfVarint((uint64(v) << 1) ^ uint64(v>>63))
}

func sizeOfVarLen(n int) int {
	return sizeOfVarint(uint64(n)) + n
}

func sizeOfTag(f int, t wireType) int {
	return sizeOfVarint(uint64(f)<<3 | uint64(t))
}
