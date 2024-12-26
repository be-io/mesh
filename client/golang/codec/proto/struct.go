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
	"github.com/opendatav/mesh/client/golang/cause"
	"io"
	"reflect"
	"strconv"
	"unsafe"
)

const (
	embedded = 1 << 0
	repeated = 1 << 1
	zigzag   = 1 << 2
)

func init() {
	Register(&structs{objects: map[reflect.Type]codec{}})
}

type structs struct {
	objects map[reflect.Type]codec
}

func (that *structs) match(ctx context.Context, kind reflect.Type) (codec, error) {
	if kind.Kind() != reflect.Struct {
		return nil, nil
	}
	cdc, err := that.construct(ctx, kind)
	if nil != err {
		return nil, cause.Error(err)
	}
	if nil == that.objects[kind] {
		that.objects[kind] = cdc
	}
	return that.objects[kind], nil
}

func FindIndex(field reflect.StructField, dft int) (int, error) {
	if tag, ok := field.Tag.Lookup(TAG); ok {
		idx, err := strconv.Atoi(tag)
		if nil != err {
			return dft, cause.Error(err)
		}
		return idx, nil
	}
	return dft, nil
}

func (that *structs) construct(ctx context.Context, kind reflect.Type) (codec, error) {
	numField := kind.NumField()
	fields := make([]*structField, 0, numField)
	for index := 0; index < numField; index++ {
		field := kind.Field(index)
		if field.PkgPath != "" {
			continue // unexported
		}
		idx, err := FindIndex(field, index)
		if nil != err {
			return nil, cause.Error(err)
		}
		cdc, err := infer(ctx, getOriginType(field.Type))
		if nil != err {
			return nil, cause.Error(err)
		}
		fields = append(fields, &structField{
			codec:  cdc,
			index:  idx,
			offset: field.Offset,
			flags:  0,
		})
	}
	return &object{kind: kind, fields: fields}, nil
}

type object struct {
	kind   reflect.Type
	fields []*structField
}

func (that *object) wire() wireType {
	return varlen
}

func (that *object) size(ctx context.Context, pointer unsafe.Pointer, flags flags) int {
	var inlined = inlined(that.kind)
	var unique, repeated []*structField

	for _, field := range that.fields {
		if field.repeated() {
			repeated = append(repeated, field)
		} else {
			unique = append(unique, field)
		}
	}
	if pointer == nil {
		return 0
	}

	if !inlined {
		flags = flags.without(inline | toplevel)
	} else {
		flags = flags.without(toplevel)
	}
	n := 0

	for _, f := range unique {
		size := f.codec.size(ctx, f.pointer(pointer), f.makeFlags(flags))
		n += size
		if size > 0 {
			if f.embedded() {
				n += sizeOfVarint(uint64(size))
			}
			flags = flags.without(wantzero)
		}
	}

	for _, f := range repeated {
		size := f.codec.size(ctx, f.pointer(pointer), f.makeFlags(flags))
		if size > 0 {
			n += size
			flags = flags.without(wantzero)
		}
	}

	return n
}

func (that *object) encode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	var inlined = inlined(that.kind)
	var unique, repeated []*structField

	for _, field := range that.fields {
		if field.repeated() {
			repeated = append(repeated, field)
		} else {
			unique = append(unique, field)
		}
	}
	if pointer == nil {
		return 0, nil
	}

	if !inlined {
		flags = flags.without(inline | toplevel)
	} else {
		flags = flags.without(toplevel)
	}
	offset := 0

	for _, f := range unique {
		fieldFlags := f.makeFlags(flags)
		elem := f.pointer(pointer)
		size := f.codec.size(ctx, elem, fieldFlags)

		if size > 0 {
			n, err := encodeTag(buff[offset:], f.index, f.codec.wire())
			offset += n
			if nil != err {
				return offset, err
			}

			if f.embedded() {
				n, err := encodeVarint(buff[offset:], uint64(size))
				offset += n
				if nil != err {
					return offset, err
				}
			}

			if (len(buff) - offset) < size {
				return len(buff), io.ErrShortBuffer
			}

			n, err = f.codec.encode(ctx, buff[offset:offset+size], elem, fieldFlags)
			offset += n
			if nil != err {
				return offset, err
			}

			flags = flags.without(wantzero)
		}
	}

	for _, f := range repeated {
		n, err := f.codec.encode(ctx, buff[offset:], f.pointer(pointer), f.makeFlags(flags))
		offset += n
		if nil != err {
			return offset, err
		}
		if n > 0 {
			flags = flags.without(wantzero)
		}
	}

	return offset, nil
}

func (that *object) decode(ctx context.Context, buff []byte, pointer unsafe.Pointer, flags flags) (int, error) {
	maxFieldNumber := 0

	for _, field := range that.fields {
		if n := field.index; n > maxFieldNumber {
			maxFieldNumber = n
		}
	}

	fieldIndex := make([]*structField, maxFieldNumber+1)

	for _, field := range that.fields {
		fieldIndex[field.index] = field
	}
	flags = flags.without(toplevel)
	offset := 0

	for offset < len(buff) {
		finx, wt, n, err := decodeTag(buff[offset:])
		offset += n
		if nil != err {
			return offset, err
		}

		i := int(finx)
		f := (*structField)(nil)

		if i >= 0 && i < len(fieldIndex) {
			f = fieldIndex[i]
		}

		if f == nil {
			skip := 0
			size := uint64(0)
			switch wt {
			case varint:
				_, skip, err = decodeVarint(buff[offset:])
			case varlen:
				size, skip, err = decodeVarint(buff[offset:])
				if err == nil {
					if size > uint64(len(buff)-skip) {
						err = io.ErrUnexpectedEOF
					} else {
						skip += int(size)
					}
				}
			case fixed32:
				_, skip, err = decodeLE32(buff[offset:])
			case fixed64:
				_, skip, err = decodeLE64(buff[offset:])
			default:
				err = cause.Errorf("unknown type %v", wt)
			}
			if (offset + skip) <= len(buff) {
				offset += skip
			} else {
				offset, err = len(buff), io.ErrUnexpectedEOF
			}
			if nil != err {
				return offset, cause.Errorf("%d %v %v", finx, wt, err)
			}
			continue
		}

		if wt != f.codec.wire() {
			return offset, cause.Errorf("expected wire type %d, %d %v", f.codec.wire(), finx, wt)
		}

		// `data` will only contain the section of the input buffer where
		// the data for the next field is available. This is necessary to
		// limit how many bytes will be consumed by embedded messages.
		var data []byte
		switch wt {
		case varint:
			_, n, err := decodeVarint(buff[offset:])
			if nil != err {
				return offset, cause.Errorf("%d %v %v", finx, wt, err)
			}
			data = buff[offset : offset+n]

		case varlen:
			l, n, err := decodeVarint(buff[offset:])
			if nil != err {
				return offset + n, cause.Errorf("%d %v %v", finx, wt, err)
			}
			if l > uint64(len(buff)-(offset+n)) {
				return len(buff), cause.Errorf("%d %v %v", finx, wt, io.ErrUnexpectedEOF)
			}
			if f.embedded() {
				offset += n
				data = buff[offset : offset+int(l)]
			} else {
				data = buff[offset : offset+n+int(l)]
			}

		case fixed32:
			if (offset + 4) > len(buff) {
				return len(buff), cause.Errorf("%d %v %v", finx, wt, io.ErrUnexpectedEOF)
			}
			data = buff[offset : offset+4]

		case fixed64:
			if (offset + 8) > len(buff) {
				return len(buff), cause.Errorf("%d %v %v", finx, wt, io.ErrUnexpectedEOF)
			}
			data = buff[offset : offset+8]

		default:
			return offset, cause.Errorf("%d %v unknown", finx, wt)
		}

		n, err = f.codec.decode(ctx, data, f.pointer(pointer), f.makeFlags(flags))
		offset += n
		if nil != err {
			return offset, cause.Errorf("%d %v %v", finx, wt, err)
		}
	}

	return offset, nil
}

type structField struct {
	codec  codec
	index  int
	offset uintptr
	flags  int
}

func (f *structField) String() string {
	return fmt.Sprintf("[%d,%s]", f.index, f.codec.wire())
}

func (f *structField) embedded() bool {
	return (f.flags & embedded) != 0
}

func (f *structField) repeated() bool {
	return (f.flags & repeated) != 0
}

func (f *structField) pointer(p unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + uintptr(f.offset))
}

func (f *structField) makeFlags(base flags) flags {
	return base | flags(f.flags&zigzag)
}
