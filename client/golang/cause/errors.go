/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cause

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	var _ Codeable = new(Cause)
}

type Cause struct {
	code string
	at   string
	err  error
}

func (that *Cause) Error() string {
	return fmt.Sprintf("%s:%s", that.GetMessage(), that.at)
}

func (that *Cause) GetCode() string {
	return that.code
}

func (that *Cause) GetMessage() string {
	if nil == that.err {
		return "Unknown"
	}
	return that.err.Error()
}

func Error(err error) error {
	if nil == err {
		return err
	}
	if cause, ok := err.(*Cause); ok {
		return cause
	}
	return &Cause{
		code: SystemError.Code,
		at:   Caller(2),
		err:  err,
	}
}

func Errorf(format string, args ...interface{}) error {
	return &Cause{
		code: SystemError.Code,
		at:   Caller(2),
		err:  fmt.Errorf(format, args...),
	}
}

func CompatibleError(format string, args ...interface{}) error {
	return Errorcf(Compatible, format, args...)
}

func ValidateError(err error) error {
	return &Cause{
		code: Validate.Code,
		at:   Caller(2),
		err:  err,
	}
}

func ValidateErrorf(format string, args ...interface{}) error {
	return &Cause{
		code: Validate.Code,
		at:   Caller(2),
		err:  fmt.Errorf(format, args...),
	}
}

func TimeoutError(err error) error {
	return &Cause{
		code: Timeout.Code,
		at:   Caller(2),
		err:  err,
	}
}

func TimeoutErrorf(format string, args ...interface{}) error {
	return &Cause{
		code: Timeout.Code,
		at:   Caller(2),
		err:  fmt.Errorf(format, args...),
	}
}

func NotFoundError(err error) error {
	return &Cause{
		code: NotFound.Code,
		at:   Caller(2),
		err:  err,
	}
}

func NotFoundErrorf(format string, args ...interface{}) error {
	return &Cause{
		code: NotFound.Code,
		at:   Caller(2),
		err:  fmt.Errorf(format, args...),
	}
}

func UnauthorizedError(err error) error {
	return &Cause{
		code: Unauthorized.Code,
		at:   Caller(2),
		err:  err,
	}
}

func UnauthorizedErrorf(format string, args ...interface{}) error {
	return &Cause{
		code: Unauthorized.Code,
		at:   Caller(2),
		err:  fmt.Errorf(format, args...),
	}
}

func NotImplementError() error {
	return &Cause{
		code: Compatible.Code,
		at:   Caller(2),
		err:  fmt.Errorf("API not implement now! "),
	}
}

func NoImplement(name string) error {
	return &Cause{
		code: Compatible.Code,
		at:   Caller(2),
		err:  fmt.Errorf("%s not present. ", name),
	}
}

func Errorable(code Codeable) error {
	return &Cause{
		code: code.GetCode(),
		at:   Caller(2),
		err:  fmt.Errorf(code.GetMessage()),
	}
}

func Errorc(code Codeable, err error) error {
	return &Cause{
		code: code.GetCode(),
		at:   Caller(2),
		err:  err,
	}
}

func Errorcf(code Codeable, format string, args ...interface{}) error {
	return &Cause{
		code: code.GetCode(),
		at:   Caller(2),
		err:  fmt.Errorf(format, args...),
	}
}

func Errorh(code int, format string) error {
	return &Cause{
		code: strconv.Itoa(code),
		at:   Caller(2),
		err:  errors.New(format),
	}
}

func Errorm(code string, message string) error {
	return &Cause{
		code: code,
		at:   Caller(2),
		err:  errors.New(message),
	}
}

func DeError(err error) error {
	if nil == err {
		return nil
	}
	if e, ok := err.(*Cause); ok {
		return DeError(e.err)
	}
	return err
}

func Coder(err error) string {
	if c, ok := err.(Codeable); ok {
		return c.GetCode()
	}
	return SystemError.Code
}

func Caller(skip int) string {
	_, name, line, _ := runtime.Caller(skip)
	if isTest() {
		return fmt.Sprintf("%s:%d", name, line)
	}
	return fmt.Sprintf("%s:%d", name[strings.LastIndex(name, "/")+1:], line)
}

// IsTest is golang testing.
func isTest() bool {
	return nil != flag.Lookup("test.v")
}

func Match(err error, causes ...Codeable) bool {
	if nil == err || len(causes) < 1 {
		return false
	}
	if e, ok := err.(*Cause); ok {
		for _, c := range causes {
			if e.GetCode() == c.GetCode() {
				return true
			}
		}
		return Match(e.err, causes...)
	}
	return false
}

func Parse(except error) (string, string) {
	if cable, ok := except.(Codeable); ok {
		return cable.GetCode(), cable.GetMessage()
	}
	rrr := DeError(except)
	if cable, ok := rrr.(Codeable); ok {
		return cable.GetCode(), cable.GetMessage()
	}
	return SystemError.GetCode(), rrr.Error()
}

func WriteHTTPError(err error, writer http.ResponseWriter) bool {
	if nil == err {
		return false
	}
	c, ok := err.(Codeable)
	if !ok {
		return false
	}
	if len(c.GetCode()) != 3 {
		return false
	}
	cc, err := strconv.Atoi(c.GetCode())
	if nil != err {
		return false
	}
	http.Error(writer, c.GetMessage(), cc)
	return true
}
