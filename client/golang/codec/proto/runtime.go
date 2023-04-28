/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

// proto exposes internal APIs of the Go runtime.
//
// This package is internal so it doesn't become part of the exported APIs that
// users of this package can take dependencies on. There is a risk that these
// APIs will be implicitly changed by Go, in which case packages that depend on
// it will break. We use these APIs to have access to optimziations that aren't
// possible today via the reflect package. Ideally, the reflect package evolves
// to expose APIs that are efficient enough that we can drop the need for this
// package, but until then we will be maintaining bridges to these Go runtime
// functions and types.

package proto

import "unsafe"

func Assign(typ, dst, src unsafe.Pointer) {
	typedmemmove(typ, dst, src)
}

func MapAssign(t, m, k unsafe.Pointer) unsafe.Pointer {
	return mapassign(t, m, k)
}

func MakeMap(t unsafe.Pointer, cap int) unsafe.Pointer {
	return makemap(t, cap)
}

type MapIter struct{ hiter }

func (it *MapIter) Init(t unsafe.Pointer, m unsafe.Pointer) {
	mapiterinit(t, m, &it.hiter)
}

func (it *MapIter) Done() {
	if it.h != nil {
		it.key = nil
		mapiternext(&it.hiter)
	}
}

func (it *MapIter) Next() {
	mapiternext(&it.hiter)
}

func (it *MapIter) HasNext() bool {
	return it.key != nil
}

func (it *MapIter) Key() unsafe.Pointer { return it.key }

func (it *MapIter) Value() unsafe.Pointer { return it.value }

// copied from src/runtime/map.go, all pointer types replaced with
// unsafe.Pointer.
//
// Alternatively we could get away with a heap allocation and only
// defining key and val if we were using reflect.mapiterinit instead,
// which returns a heap-allocated *hiter.
type hiter struct {
	key         unsafe.Pointer // nil when iteration is done
	value       unsafe.Pointer
	t           unsafe.Pointer
	h           unsafe.Pointer
	buckets     unsafe.Pointer // bucket ptr at hash_iter initialization time
	bptr        unsafe.Pointer // current bucket
	overflow    unsafe.Pointer // keeps overflow buckets of hmap.buckets alive
	oldoverflow unsafe.Pointer // keeps overflow buckets of hmap.oldbuckets alive
	startBucket uintptr        // bucket iteration started at
	offset      uint8          // intra-bucket offset to start from during iteration (should be big enough to hold bucketCnt-1)
	wrapped     bool           // already wrapped around from end of bucket array to beginning
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

//go:noescape
//go:linkname makemap reflect.makemap
func makemap(t unsafe.Pointer, cap int) unsafe.Pointer

// m escapes into the return value, but the caller of mapiterinit
// doesn't let the return value escape.
//go:noescape
//go:linkname mapiterinit runtime.mapiterinit
func mapiterinit(t unsafe.Pointer, m unsafe.Pointer, it *hiter)

//go:noescape
//go:linkname mapiternext runtime.mapiternext
func mapiternext(it *hiter)

//go:noescape
//go:linkname mapassign runtime.mapassign
func mapassign(t, m, k unsafe.Pointer) unsafe.Pointer

//go:nosplit
//go:noescape
//go:linkname typedmemmove runtime.typedmemmove
func typedmemmove(typ, dst, src unsafe.Pointer)

type Slice struct {
	data unsafe.Pointer
	len  int
	cap  int
}

func (s *Slice) Cap() int {
	return s.cap
}

func (s *Slice) Len() int {
	return s.len
}

func (s *Slice) SetLen(n int) {
	s.len = n
}

func (s *Slice) Index(i int, elemSize uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(s.data) + (uintptr(i) * elemSize))
}

func MakeSlice(elemType unsafe.Pointer, len, cap int) Slice {
	return Slice{
		data: newarray(elemType, cap),
		len:  len,
		cap:  cap,
	}
}

func CopySlice(elemType unsafe.Pointer, dst, src Slice) int {
	return typedslicecopy(elemType, dst, src)
}

//go:linkname newarray runtime.newarray
func newarray(t unsafe.Pointer, n int) unsafe.Pointer

//go:linkname typedslicecopy runtime.typedslicecopy
//go:noescape
func typedslicecopy(t unsafe.Pointer, dst, src Slice) int
