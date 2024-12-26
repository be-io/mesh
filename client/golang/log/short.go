/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package log

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"io"
	gog "log"
	"runtime/debug"
	"time"
)

func init() {
	gog.SetOutput(mog.Writer())
}

var mog = Load()
var mof Formatter = new(Stringify)

func Set(log Logger) {
	mog = log
	gog.SetOutput(mog.Writer())
}

func SetFormatter(formatter Formatter) {
	if nil != formatter {
		mof = formatter
	}
}

func Format(ctx Context, time time.Time, level Level, msg string, ref string) string {
	return mof.Format(ctx, time, level, msg, ref)
}

func Formatln(ctx Context, time time.Time, level Level, msg string, ref string) string {
	return mof.Formatln(ctx, time, level, msg, ref)
}

func With(key string, value string) Logger {
	return mog
}

func Info0(format string, args ...interface{}) {
	mog.Info0(format, args...)
}

func Warn0(format string, args ...interface{}) {
	mog.Warn0(format, args...)
}

func Error0(format string, args ...interface{}) {
	mog.Error0(format, args...)
}

func Debug0(format string, args ...interface{}) {
	mog.Debug0(format, args...)
}

func Fatal0(format string, args ...interface{}) {
	mog.Fatal0(format, args...)
}

func Stack0(format string, args ...interface{}) {
	mog.Stack0(format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	mog.Info(ctx, format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	mog.Warn(ctx, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	mog.Error(ctx, format, args...)
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	mog.Debug(ctx, format, args...)
}

func Fatal(ctx context.Context, format string, args ...interface{}) {
	mog.Fatal(ctx, format, args...)
}

func Stack(ctx context.Context, format string, args ...interface{}) {
	mog.Stack(ctx, format, args...)
}

func Print(ctx context.Context, level Level, format string, args ...interface{}) {
	mog.Print(ctx, level, format, args...)
}

func Catch(err error) {
	mog.Catch(err)
}

func Panic(err error) {
	mog.Panic(err)
}

func Writer() io.Writer {
	return mog.Writer()
}

// Devour wraps a `go func()` with recover()
func Devour(handler func()) {
	Recover(handler, nil)
}

// Recover wraps a `go func()` with recover()
func Recover(handler func(), recoverHandler func(r interface{})) {
	defer func() {
		if err := recover(); nil != err {
			Error0("%v", err)
			if nil != recoverHandler {
				go func() {
					defer func() {
						if err := recover(); nil != err {
							Error0("%v", err)
						}
					}()
					recoverHandler(err)
				}()
			}
		}
	}()
	handler()
}

// PError recover panic as error
func PError(ctx context.Context, fn func() error) error {
	_, err := PRError(ctx, func() (interface{}, error) {
		return nil, fn()
	})
	return err
}

// PRError recover panic as error
func PRError[T any](ctx context.Context, fn func() (T, error)) (t T, err error) {
	defer func() {
		if ca := recover(); nil != ca {
			err = cause.Errorf("%v", ca)
			Error(ctx, err.Error())
			Error(ctx, string(debug.Stack()))
		}
	}()
	return fn()
}

func SetLevel(level Level) {
	mog.Level(level)
}
