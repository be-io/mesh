/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tbucket

import "time"

type RateLimiter interface {

	// Wait takes count tokens from the bucket, waiting until they are
	// available.
	Wait(count int64)

	// Rate returns the fill rate of the bucket, in tokens per second.
	Rate() float64

	// Capacity returns the capacity that the bucket was created with.
	Capacity() int64

	// Available returns the number of available tokens. It will be negative
	// when there are consumers waiting for tokens. Note that if this
	// returns greater than zero, it does not guarantee that calls that take
	// tokens from the buffer will succeed, as the number of available
	// tokens could have changed in the meantime. This method is intended
	// primarily for metrics reporting and debugging.
	Available() int64

	// Take takes count tokens from the bucket without blocking. It returns
	// the time that the caller should wait until the tokens are actually
	// available.
	//
	// Note that if the request is irrevocable - there is no way to return
	// tokens to the bucket once this method commits us to taking them.
	Take(count int64) time.Duration

	// TakeAvailable takes up to count immediately available tokens from the
	// bucket. It returns the number of tokens removed, or zero if there are
	// no available tokens. It does not block.
	TakeAvailable(count int64) int64

	// TakeMaxDuration is like Take, except that
	// it will only take tokens from the bucket if the wait
	// time for the tokens is no greater than maxWait.
	//
	// If it would take longer than maxWait for the tokens
	// to become available, it does nothing and reports false,
	// otherwise it returns the time that the caller should
	// wait until the tokens are actually available, and reports
	// true.
	TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool)

	// WaitMaxDuration is like Wait except that it will
	// only take tokens from the bucket if it needs to wait
	// for no greater than maxWait. It reports whether
	// any tokens have been removed from the bucket
	// If no tokens have been removed, it returns immediately.
	WaitMaxDuration(count int64, maxWait time.Duration) bool
}
