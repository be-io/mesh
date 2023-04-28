/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package log

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	var _ io.Closer = new(zero)
	var _ Logger = new(zero)
	Provide(new(zero))
}

var zer = func() Logger {
	return &zero{ctx: macro.Context()}
}()

type zero struct {
	ctx     context.Context
	once    sync.Once
	writers []io.WriteCloser
	name    string
}

func (that *zero) GetTraceId() string {
	return ""
}

func (that *zero) GetSpanId() string {
	return ""
}

func (that *zero) GetTimestamp() int64 {
	return 0
}

func (that *zero) Name() string {
	return "zero"
}

func (that *zero) Info0(format string, args ...interface{}) {
	that.Print(that.ctx, INFO, format, args...)
}

func (that *zero) Warn0(format string, args ...interface{}) {
	that.Print(that.ctx, WARN, format, args...)
}

func (that *zero) Error0(format string, args ...interface{}) {
	that.Print(that.ctx, ERROR, format, args...)
}

func (that *zero) Debug0(format string, args ...interface{}) {
	that.Print(that.ctx, DEBUG, format, args...)
}

func (that *zero) Fatal0(format string, args ...interface{}) {
	that.Print(that.ctx, FATAL, format, args...)
}

func (that *zero) Stack0(format string, args ...interface{}) {
	that.Print(that.ctx, STACK, format, args...)
}

func (that *zero) Info(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, INFO, format, args...)
}

func (that *zero) Warn(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, WARN, format, args...)
}

func (that *zero) Error(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, ERROR, format, args...)
}

func (that *zero) Debug(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, DEBUG, format, args...)
}

func (that *zero) Fatal(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, FATAL, format, args...)
}

func (that *zero) Stack(ctx context.Context, format string, args ...interface{}) {
	that.Print(ctx, STACK, format, args...)
}

func (that *zero) Print(ctx context.Context, level Level, format string, args ...interface{}) {
	that.once.Do(func() { that.Init() })
	ntx, ok := ctx.(Context)
	if !ok {
		ntx = that
	}
	params, pairs := Split0(ctx, args)
	ref := Caller(4)
	that.Event(level).
		Str("name", that.name).
		Str("ref", ref).
		Str("trace_id", ntx.GetTraceId()).
		Str("span_id", ntx.GetSpanId()).
		Fields(pairs).
		Msgf(format, params...)
}

func (that *zero) Catch(err error) {
	if nil != err {
		that.Print(that.ctx, ERROR, err.Error())
	}
}

func (that *zero) Panic(err error) {
	if nil != err {
		that.Print(that.ctx, FATAL, err.Error())
	}
}

func (that *zero) Writer() io.Writer {
	return that
}

func (that *zero) Write(p []byte) (n int, err error) {
	if len(that.writers) < 1 {
		return nop.Writer().Write(p)
	}
	for _, writer := range that.writers {
		n, err = writer.Write(p)
		if nil != err {
			_, _ = nop.Writer().Write([]byte(err.Error()))
		}
	}
	return
}

func (that *zero) Close() error {
	for _, writer := range that.writers {
		if err := writer.Close(); nil != err {
			_, _ = nop.Writer().Write([]byte(err.Error()))
		}
	}
	return nil
}

func (that *zero) Level(level Level) {
	if ALL.Is(int(level)) {
		log.Level(zerolog.DebugLevel)
		return
	}
	if STACK.Is(int(level)) {
		log.Level(zerolog.DebugLevel)
		return
	}
	if DEBUG.Is(int(level)) {
		log.Level(zerolog.DebugLevel)
		return
	}
	if INFO.Is(int(level)) {
		log.Level(zerolog.InfoLevel)
		return
	}
	if WARN.Is(int(level)) {
		log.Level(zerolog.WarnLevel)
		return
	}
	if ERROR.Is(int(level)) {
		log.Level(zerolog.ErrorLevel)
		return
	}
	if FATAL.Is(int(level)) {
		log.Level(zerolog.FatalLevel)
		return
	}
	log.Level(zerolog.InfoLevel)
}

func (that *zero) Event(level Level) *zerolog.Event {
	switch level {
	case FATAL:
		return log.Fatal()
	case ERROR:
		return log.Error()
	case WARN:
		return log.Warn()
	case INFO:
		return log.Info()
	case DEBUG:
		return log.Debug()
	case STACK:
		return log.Debug()
	case ALL:
		return log.Debug()
	default:
		return log.Info()
	}
}

// Init
// 11G, 120 backups, 100M
func (that *zero) Init() {
	uri, err := func() (*url.URL, error) {
		home, _ := os.UserHomeDir()
		for _, vp := range []string{
			macro.LHome(),
			fmt.Sprintf("file:///var/log/be/%s/%s.log?size=100&backups=120&age=28&compress=1", macro.Name(), macro.Name()),
			fmt.Sprintf("file:///%s", filepath.Join(home, "logs", fmt.Sprintf("%s.log?size=100&backups=120&age=28&compress=1", macro.Name()))),
		} {
			if "" == vp {
				continue
			}
			uri, err := url.Parse(vp)
			if nil != err {
				nop.Warn(that.ctx, err.Error())
				continue
			}
			if err = that.Make(uri.Path); nil != err {
				nop.Warn(that.ctx, err.Error())
				continue
			}
			nop.Info(that.ctx, uri.String())
			return uri, nil
		}
		panic("Unexpected log path. ")
	}()
	if nil != err {
		nop.Error(that.ctx, err.Error())
	}
	parse := func(name string, dft int) int {
		sv := uri.Query().Get(name)
		if "" == sv {
			return dft
		}
		iv, e := strconv.Atoi(sv)
		if nil != e {
			nop.Error(that.ctx, err.Error())
			return dft
		}
		return iv
	}
	// 10G = 100 * 100MB
	fw := &lumberjack.Logger{
		Filename:   uri.Path,
		MaxSize:    parse("size", 100),        // megabytes
		MaxBackups: parse("backups", 100),     // size
		MaxAge:     parse("age", 28),          //days
		Compress:   0 != parse("compress", 0), // disabled by default
	}
	that.name = macro.Name()
	if macro.JsonLogFormat.Enable() {
		cw := zerolog.ConsoleWriter{
			Out:             os.Stdout,
			TimeFormat:      DateFormat23,
			FormatTimestamp: that.FormatTimestamp,
		}
		if macro.NoStdColor.Enable() {
			cw.NoColor = true
			cw.FormatLevel = that.FormatLevel
		}
		zerolog.LevelTraceValue = strings.ToUpper(zerolog.LevelTraceValue)
		zerolog.LevelDebugValue = strings.ToUpper(zerolog.LevelDebugValue)
		zerolog.LevelInfoValue = strings.ToUpper(zerolog.LevelInfoValue)
		zerolog.LevelWarnValue = strings.ToUpper(zerolog.LevelWarnValue)
		zerolog.LevelErrorValue = strings.ToUpper(zerolog.LevelErrorValue)
		zerolog.LevelFatalValue = strings.ToUpper(zerolog.LevelFatalValue)
		zerolog.LevelPanicValue = strings.ToUpper(zerolog.LevelPanicValue)
		zerolog.LevelFieldMarshalFunc = that.FormatZLevel
		that.writers = append(that.writers, &WriterCloser{writer: cw}, fw)
	} else {
		cw := zerolog.ConsoleWriter{
			Out:             &Writers{writers: []io.Writer{os.Stdout, fw}},
			TimeFormat:      DateFormat23,
			FormatTimestamp: that.FormatTimestamp,
		}
		if macro.NoStdColor.Enable() {
			cw.NoColor = true
			cw.FormatLevel = that.FormatLevel
		}
		that.writers = append(that.writers, &WriterCloser{writer: cw, closer: fw})
	}
	zerolog.TimeFieldFormat = DateFormat23
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	log.Logger = zerolog.New(that).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &log.Logger
	AddProcessShutdownHook(func() error { return that.Close() })
}

func (that *zero) Make(path string) error {
	if err := that.makeDir(filepath.Dir(path)); nil != err {
		return err
	}
	return that.makeFile(path)
}

func (that *zero) makeDir(path string) error {
	_, err := os.Stat(path)
	if nil == err || !os.IsNotExist(err) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

func (that *zero) makeFile(path string) error {
	_, err := os.Stat(path)
	if nil == err || !os.IsNotExist(err) {
		return nil
	}
	fd, err := os.Create(path)
	if nil != err {
		return err
	}
	return fd.Close()
}

func (that *zero) FormatLevel(i interface{}) string {
	if v, ok := i.(string); ok {
		switch v {
		case zerolog.LevelTraceValue, "TRACE":
			return "TRACE"
		case zerolog.LevelDebugValue, "DEBUG":
			return "DEBUG"
		case zerolog.LevelInfoValue, "INFO":
			return "INFO"
		case zerolog.LevelWarnValue, "WARN":
			return "WARN"
		case zerolog.LevelErrorValue, "ERROR":
			return "ERROR"
		case zerolog.LevelFatalValue, "FATAL":
			return "FATAL"
		case zerolog.LevelPanicValue, "PANIC":
			return "PANIC"
		default:
			return "???"
		}
	} else {
		if i == nil {
			return "???"
		} else {
			return strings.ToUpper(fmt.Sprintf("%v", i))
		}
	}
}

func (that *zero) FormatZLevel(level zerolog.Level) string {
	switch level {
	case zerolog.TraceLevel:
		return "TRACE"
	case zerolog.DebugLevel:
		return "DEBUG"
	case zerolog.InfoLevel:
		return "INFO"
	case zerolog.WarnLevel:
		return "WARN"
	case zerolog.ErrorLevel:
		return "ERROR"
	case zerolog.FatalLevel:
		return "FATAL"
	case zerolog.PanicLevel:
		return "PANIC"
	case zerolog.Disabled:
		return "DISABLED"
	case zerolog.NoLevel:
		return ""
	}
	return strconv.Itoa(int(level))
}

func (that *zero) FormatTimestamp(i interface{}) string {
	switch tt := i.(type) {
	case string:
		return tt
	case json.Number:
		i64, err := tt.Int64()
		if nil != err {
			return tt.String()
		}
		var sec, nsec int64 = i64, 0
		switch zerolog.TimeFieldFormat {
		case zerolog.TimeFormatUnixMs:
			nsec = int64(time.Duration(i64) * time.Millisecond)
			sec = 0
		case zerolog.TimeFormatUnixMicro:
			nsec = int64(time.Duration(i64) * time.Microsecond)
			sec = 0
		}
		ts := time.Unix(sec, nsec)
		return ts.Format(DateFormat23)
	default:
		return fmt.Sprintf("%v", i)
	}
}

type WriterCloser struct {
	writer io.Writer
	closer io.Closer
}

func (that *WriterCloser) Write(p []byte) (n int, err error) {
	return that.writer.Write(p)
}

func (that *WriterCloser) Close() error {
	if nil != that.closer {
		return that.closer.Close()
	}
	return nil
}

type Writers struct {
	writers []io.Writer
}

func (that *Writers) Write(p []byte) (n int, err error) {
	if len(that.writers) < 1 {
		return nop.Writer().Write(p)
	}
	for _, writer := range that.writers {
		n, err = writer.Write(p)
		if nil != err {
			_, _ = nop.Writer().Write([]byte(err.Error()))
		}
	}
	return
}
