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
	"io"
	"time"
)

func init() {
	var _ Context = new(TraceContext)
	var _ Formatter = new(Stringify)
	var _ Formatter = new(Jsonify)
}

const (
	DateFormat   = "2006-01-02 15:04:05"
	DateFormat23 = "2006-01-02 15:04:05.000"
)

// Logger is topest dependencies
type Logger interface {
	Name() string

	Info0(format string, args ...interface{})

	Warn0(format string, args ...interface{})

	Error0(format string, args ...interface{})

	Debug0(format string, args ...interface{})

	Fatal0(format string, args ...interface{})

	Stack0(format string, args ...interface{})

	Info(ctx context.Context, format string, args ...interface{})

	Warn(ctx context.Context, format string, args ...interface{})

	Error(ctx context.Context, format string, args ...interface{})

	Debug(ctx context.Context, format string, args ...interface{})

	Fatal(ctx context.Context, format string, args ...interface{})

	Stack(ctx context.Context, format string, args ...interface{})

	Print(ctx context.Context, level Level, format string, args ...interface{})

	Catch(err error)

	Panic(err error)

	Writer() io.Writer

	Level(level Level)
}

type Context interface {

	// GetTraceId the request trace id.
	GetTraceId() string

	// GetSpanId the request span id.
	GetSpanId() string

	// GetTimestamp the request create time.
	GetTimestamp() int64
}

type Appender interface {
	Name() string
}

type Pair struct {
	Key   string
	Value interface{}
}

type TraceContext struct {
	TraceId   string
	SpanId    string
	Timestamp int64
}

func (that *TraceContext) GetTraceId() string {
	return that.TraceId
}

func (that *TraceContext) GetSpanId() string {
	return that.SpanId
}

func (that *TraceContext) GetTimestamp() int64 {
	return that.Timestamp
}

type Formatter interface {

	// Format the log
	Format(ctx Context, time time.Time, level Level, msg string, ref string) string

	// Formatln the log with line
	Formatln(ctx Context, time time.Time, level Level, msg string, ref string) string
}

type Stringify struct {
}

func (that *Stringify) Formatln(ctx Context, time time.Time, level Level, msg string, ref string) string {
	return fmt.Sprintf("[mesh] %s [%s] %s#%s %s %s\n", time.Format(DateFormat), level.String(), ctx.GetTraceId(), ctx.GetSpanId(), msg, ref)
}

func (that *Stringify) Format(ctx Context, time time.Time, level Level, msg string, ref string) string {
	return fmt.Sprintf("[mesh] %s [%s] %s#%s %s %s", time.Format(DateFormat), level.String(), ctx.GetTraceId(), ctx.GetSpanId(), msg, ref)
}

type Jsonify struct {
}

func (that *Jsonify) Formatln(ctx Context, time time.Time, level Level, msg string, ref string) string {
	return fmt.Sprintf("{\"name\":\"mesh\",\"timestamp\":\"%s\",\"level\":\"%s\",\"trace_id\":\"%s\",\"span_id\":\"%s\",\"msg\":\"%s\",\"ref\":\"%s\"}\n", time.Format(DateFormat), level.String(), ctx.GetTraceId(), ctx.GetSpanId(), msg, ref)
}

func (that *Jsonify) Format(ctx Context, time time.Time, level Level, msg string, ref string) string {
	return fmt.Sprintf("{\"name\":\"mesh\",\"timestamp\":\"%s\",\"level\":\"%s\",\"trace_id\":\"%s\",\"span_id\":\"%s\",\"msg\":\"%s\",\"ref\":\"%s\"}", time.Format(DateFormat), level.String(), ctx.GetTraceId(), ctx.GetSpanId(), msg, ref)
}
