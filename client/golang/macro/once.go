/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import "reflect"

type Once[T any] struct {
	value    T
	producer func() T
}

func (that *Once[T]) With(producer func() T) *Once[T] {
	that.producer = producer
	return that
}

func (that *Once[T]) GetWith(producer func() T) T {
	that.producer = producer
	return that.Get()
}

func (that *Once[T]) Get() T {
	if !reflect.ValueOf(&that.value).Elem().IsZero() {
		return that.value
	}
	if nil != that.producer {
		that.value = that.producer()
	}
	return that.value
}

func (that *Once[T]) Set(v T) {
	that.value = v
}
