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
	"sync"
	"unsafe"
)

const (
	zeroSize = 1 // sizeOfVarint(0)
)

func init() {
	Register(new(dicts))
}

type dicts struct {
}

func (that *dicts) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if reflect.Map != kind.Kind() {
		return nil, nil
	}
	key, val := kind.Key(), kind.Elem()
	k, err := infer(ctx, key)
	if nil != err {
		return nil, cause.Error(err)
	}
	v, err := infer(ctx, val)
	if nil != err {
		return nil, cause.Error(err)
	}
	cdc := &dict{
		kind:     kind,
		number:   0,
		keyFlags: 0,
		valFlags: embedded | repeated,
		keyCodec: k,
		valCodec: v,
	}
	return cdc, nil
}

type dict struct {
	kind     reflect.Type
	number   uint16
	keyFlags uint8
	valFlags uint8
	keyCodec codec
	valCodec codec
}

func (that *dict) wire() wireType {
	return varlen
}

func (that *dict) size(ctx context.Context, ptr unsafe.Pointer, flags flags) int {
	mapTagSize := sizeOfTag(int(that.number), varlen)
	keyTagSize := sizeOfTag(1, that.keyCodec.wire())
	valTagSize := sizeOfTag(2, that.valCodec.wire())
	if ptr == nil {
		return 0
	}

	if !flags.has(inline) {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	n := 0
	m := MapIter{}
	defer m.Done()

	for m.Init(pointer(that.kind), ptr); m.HasNext(); m.Next() {
		keySize := that.keyCodec.size(ctx, m.Key(), wantzero)
		valSize := that.valCodec.size(ctx, m.Value(), wantzero)

		if keySize > 0 {
			n += keyTagSize + keySize
			if (that.keyFlags & embedded) != 0 {
				n += sizeOfVarint(uint64(keySize))
			}
		}

		if valSize > 0 {
			n += valTagSize + valSize
			if (that.valFlags & embedded) != 0 {
				n += sizeOfVarint(uint64(valSize))
			}
		}

		n += mapTagSize + sizeOfVarint(uint64(keySize+valSize))
	}

	if n == 0 {
		n = mapTagSize + zeroSize
	}

	return n
}

func (that *dict) encode(ctx context.Context, buff []byte, p unsafe.Pointer, flags flags) (int, error) {
	keyTag := [1]byte{}
	valTag := [1]byte{}
	encodeTag(keyTag[:], 1, that.keyCodec.wire())
	encodeTag(valTag[:], 2, that.valCodec.wire())

	number := fieldNumber(that.number)
	mapTag := make([]byte, sizeOfTag(int(number), varlen)+zeroSize)
	encodeTag(mapTag, int(number), varlen)

	zero := mapTag
	mapTag = mapTag[:len(mapTag)-1]

	if p == nil {
		return 0, nil
	}

	if !flags.has(inline) {
		p = *(*unsafe.Pointer)(p)
	}

	offset := 0
	m := MapIter{}
	defer m.Done()

	for m.Init(pointer(that.kind), p); m.HasNext(); m.Next() {
		key := m.Key()
		val := m.Value()

		keySize := that.keyCodec.size(ctx, key, wantzero)
		valSize := that.valCodec.size(ctx, val, wantzero)
		elemSize := keySize + valSize

		if keySize > 0 {
			elemSize += len(keyTag)
			if (that.keyFlags & embedded) != 0 {
				elemSize += sizeOfVarint(uint64(keySize))
			}
		}

		if valSize > 0 {
			elemSize += len(valTag)
			if (that.valFlags & embedded) != 0 {
				elemSize += sizeOfVarint(uint64(valSize))
			}
		}

		n := copy(buff[offset:], mapTag)
		offset += n
		if n < len(mapTag) {
			return offset, io.ErrShortBuffer
		}
		n, err := encodeVarint(buff[offset:], uint64(elemSize))
		offset += n
		if nil != err {
			return offset, err
		}

		if keySize > 0 {
			n := copy(buff[offset:], keyTag[:])
			offset += n
			if n < len(keyTag) {
				return offset, io.ErrShortBuffer
			}

			if (that.keyFlags & embedded) != 0 {
				n, err := encodeVarint(buff[offset:], uint64(keySize))
				offset += n
				if nil != err {
					return offset, err
				}
			}

			if (len(buff) - offset) < keySize {
				return len(buff), io.ErrShortBuffer
			}

			n, err := that.keyCodec.encode(ctx, buff[offset:offset+keySize], key, wantzero)
			offset += n
			if nil != err {
				return offset, err
			}
		}

		if valSize > 0 {
			n := copy(buff[offset:], valTag[:])
			offset += n
			if n < len(valTag) {
				return n, io.ErrShortBuffer
			}

			if (that.valFlags & embedded) != 0 {
				n, err := encodeVarint(buff[offset:], uint64(valSize))
				offset += n
				if nil != err {
					return offset, err
				}
			}

			if (len(buff) - offset) < valSize {
				return len(buff), io.ErrShortBuffer
			}

			n, err := that.valCodec.encode(ctx, buff[offset:offset+valSize], val, wantzero)
			offset += n
			if nil != err {
				return offset, err
			}
		}
	}

	if offset == 0 {
		if offset = copy(buff, zero); offset < len(zero) {
			return offset, io.ErrShortBuffer
		}
	}

	return offset, nil
}

func (that *dict) decode(ctx context.Context, buff []byte, p unsafe.Pointer, flags flags) (int, error) {
	st := reflect.StructOf([]reflect.StructField{
		{Name: "Key", Type: that.kind.Key()},
		{Name: "Elem", Type: that.kind.Elem()},
	})

	structCodec, err := infer(ctx, st)
	if nil != err {
		return 0, cause.Error(err)
	}
	structPool := new(sync.Pool)
	structZero := pointer(reflect.Zero(st).Interface())

	valueType := that.kind.Elem()
	valueOffset := st.Field(1).Offset

	mtype := pointer(that.kind)
	stype := pointer(st)
	vtype := pointer(valueType)

	m := (*unsafe.Pointer)(p)
	if *m == nil {
		*m = MakeMap(mtype, 10)
	}
	if len(buff) == 0 {
		return 0, nil
	}

	s := pointer(structPool.Get())
	if s == nil {
		s = unsafe.Pointer(reflect.New(st).Pointer())
	}

	n, err := structCodec.decode(ctx, buff, s, noflags)
	if err == nil {
		v := MapAssign(mtype, *m, s)
		Assign(vtype, v, unsafe.Pointer(uintptr(s)+valueOffset))
	}

	Assign(stype, s, structZero)
	structPool.Put(s)
	return n, err
}
