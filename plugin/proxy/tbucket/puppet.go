/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tbucket

import (
	"math"
	"time"
)

func NewFreshRateLimiter(rate float64, capacity int64, minRate float64, minCap int64) FreshRateLimiter {
	return &puppet{
		rate:     rate,
		capacity: capacity,
		limiter:  NewBucketWithRate(math.Max(rate, minRate), int64(math.Max(float64(capacity), float64(minCap)))),
	}
}

var _ FreshRateLimiter = new(puppet)

type FreshRateLimiter interface {
	RateLimiter
	Refresh(rate float64, capacity int64, minRate float64, minCap int64)
}

type puppet struct {
	rate     float64
	capacity int64
	limiter  RateLimiter
}

func (that *puppet) Free() bool {
	return that.rate < 1 || that.capacity < 1
}

func (that *puppet) Refresh(rate float64, capacity int64, minRate float64, minCap int64) {
	that.rate = rate
	that.capacity = capacity
	that.limiter = NewBucketWithRate(math.Max(rate, minRate), int64(math.Max(float64(capacity), float64(minCap))))
}

func (that *puppet) Wait(count int64) {
	if that.Free() {
		return
	}
	that.limiter.Wait(count)
}

func (that *puppet) Rate() float64 {
	if that.Free() {
		return math.MaxFloat64
	}
	return that.limiter.Rate()
}

func (that *puppet) Capacity() int64 {
	if that.Free() {
		return math.MaxInt64
	}
	return that.limiter.Capacity()
}

func (that *puppet) Available() int64 {
	if that.Free() {
		return math.MaxInt64
	}
	return that.limiter.Available()
}

func (that *puppet) Take(count int64) time.Duration {
	if that.Free() {
		return 0
	}
	return that.limiter.Take(count)
}

func (that *puppet) TakeAvailable(count int64) int64 {
	if that.Free() {
		return math.MaxInt64
	}
	return that.limiter.TakeAvailable(count)
}

func (that *puppet) TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool) {
	if that.Free() {
		return 0, true
	}
	return that.limiter.TakeMaxDuration(count, maxWait)
}

func (that *puppet) WaitMaxDuration(count int64, maxWait time.Duration) bool {
	if that.Free() {
		return true
	}
	return that.limiter.WaitMaxDuration(count, maxWait)
}
