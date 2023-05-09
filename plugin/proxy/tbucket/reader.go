/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tbucket

import "io"

type reader struct {
	r      io.Reader
	bucket RateLimiter
}

// Reader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Reader(r io.Reader, bucket RateLimiter) io.Reader {
	return &reader{
		r:      r,
		bucket: bucket,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.bucket.Wait(int64(n))
	return n, err
}

type writer struct {
	w      io.Writer
	bucket RateLimiter
}

// Writer returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Writer(w io.Writer, bucket RateLimiter) io.Writer {
	return &writer{
		w:      w,
		bucket: bucket,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}
