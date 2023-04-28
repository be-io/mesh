/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package dsa

import "sort"

var _ List[any] = new(list[any])

type List[T any] interface {
	Gt(index int) T
	Ad(v T)
	Rm(index int)
	Len() int
	Peek() T
	Clone() []T
	Group(fn func(v T) string) Map[string, List[T]]
}

func NewSortList[T any](less func(i, j T) bool) List[T] {
	return &list[T]{less: less}
}

func NewList[T any]() List[T] {
	return &list[T]{}
}

func NewArrayList[T any](array []T) List[T] {
	li := &list[T]{}
	for _, v := range array {
		li.Ad(v)
	}
	return li
}

type list[T any] struct {
	elements []T
	less     func(i, j T) bool
}

func (that *list[T]) Gt(index int) T {
	return that.elements[index]
}

func (that *list[T]) Ad(v T) {
	that.elements = append(that.elements, v)
	if nil != that.less {
		sort.SliceStable(that.elements, func(i, j int) bool {
			return that.less(that.elements[i], that.elements[j])
		})
	}
}

func (that *list[T]) Rm(index int) {
	that.elements = append(that.elements[:index], that.elements[index+1:]...)
}

func (that *list[T]) Len() int {
	return len(that.elements)
}

func (that *list[T]) Peek() T {
	if len(that.elements) < 1 {
		var t T
		return t
	}
	return that.elements[0]
}

func (that *list[T]) Clone() []T {
	var elements []T
	copy(that.elements, elements)
	return elements
}

func (that *list[T]) Group(fn func(v T) string) Map[string, List[T]] {
	m := NewStringMap[List[T]]()
	for _, e := range that.elements {
		m.PutIfy(fn(e), func(k string) List[T] { return NewList[T]() }).Ad(e)
	}
	return m
}
