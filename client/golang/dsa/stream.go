/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package dsa

var _ Transformer[any, any] = new(streamTransformer[any, any])

type Transformer[R any, T any] interface {
	Map(rs []R, err error) ([]T, error)
}

type streamTransformer[R any, T any] struct {
	transformer func(r R) T
}

func (that *streamTransformer[R, T]) Map(rs []R, err error) ([]T, error) {
	if nil != err {
		return nil, err
	}
	if len(rs) < 1 {
		return nil, err
	}
	var ts []T
	for _, v := range rs {
		ts = append(ts, that.transformer(v))
	}
	return ts, nil
}

func Transform[R any, T any](transformer func(r R) T) Transformer[R, T] {
	return &streamTransformer[R, T]{transformer: transformer}
}
