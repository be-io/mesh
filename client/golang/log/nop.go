/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package log

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/macro"
	"io"
	"time"
)

func init() {
	Provide(nop)
}

var nop = func() Logger {
	return &std{context: macro.Context(), level: INFO}
}()

type std struct {
	context context.Context
	level   Level
}

func (that *std) GetTraceId() string {
	return ""
}

func (that *std) GetSpanId() string {
	return ""
}

func (that *std) GetTimestamp() int64 {
	return 0
}

func (that *std) Name() string {
	return "std"
}

func (that *std) Info0(format string, args ...interface{}) {
	that.Print(that.context, INFO, format, args...)
}

func (that *std) Warn0(format string, args ...interface{}) {
	that.Print(that.context, WARN, format, args...)
}

func (that *std) Error0(format string, args ...interface{}) {
	that.Print(that.context, ERROR, format, args...)
}

func (that *std) Debug0(format string, args ...interface{}) {
	that.Print(that.context, DEBUG, format, args...)
}

func (that *std) Fatal0(format string, args ...interface{}) {
	that.Print(that.context, FATAL, format, args...)
}

func (that *std) Stack0(format string, args ...interface{}) {
	that.Print(that.context, STACK, format, args...)
}

func (that *std) Info(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, INFO, format, args...)
}

func (that *std) Warn(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, WARN, format, args...)
}

func (that *std) Error(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, ERROR, format, args...)
}

func (that *std) Debug(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, DEBUG, format, args...)
}

func (that *std) Fatal(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, FATAL, format, args...)
}

func (that *std) Stack(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, STACK, format, args...)
}

func (that *std) Catch(err error) {
	if nil != err {
		that.Print(that.context, ALL, err.Error())
	}
}

func (that *std) Panic(err error) {
	if nil != err {
		that.Print(that.context, ALL, err.Error())
		panic(err)
	}
}

func (that *std) Writer() io.Writer {
	return that
}

func (that *std) Level(level Level) {
	that.level = level
}

func (that *std) Print(ctx context.Context, level Level, format string, args ...interface{}) {
	ltx, ok := ctx.(Context)
	if !ok {
		ltx = that
	}
	if !level.Is(int(that.level)) {
		return
	}
	params, _ := Split0(ctx, args)
	msg := fmt.Sprintf(format, params...)
	ref := Caller(4)
	fmt.Println(mof.Format(ltx, time.Now(), level, msg, ref))
}

func (that *std) Write(text []byte) (n int, err error) {
	fmt.Print(string(text))
	return len(text), nil
}
