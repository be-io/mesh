/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
)

var _ macro.Returns = new(Outbound)

type Outbound struct {
	Code    string       `index:"0" json:"code,omitempty" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message,omitempty" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause,omitempty" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content interface{}  `index:"15" json:"content,omitempty" xml:"content" yaml:"content" comment:"Service response content"`
}

func (that *Outbound) GetCode() string {
	return that.Code
}

func (that *Outbound) SetCode(code string) {
	that.Code = code
}

func (that *Outbound) GetMessage() string {
	return that.Message
}

func (that *Outbound) SetMessage(message string) {
	that.Message = message
}

func (that *Outbound) GetContent(ctx context.Context) []interface{} {
	return []interface{}{that.Content}
}

func (that *Outbound) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		that.Content = arguments[0]
	}
}

func (that *Outbound) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *Outbound) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}
