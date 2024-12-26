/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proto

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/macro"
)

func init() {
	var _ macro.Returns = new(PsiReturns)
}

type PsiReturns struct {
	Code    string              `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string              `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause        `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content *PsiExecuteResponse `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *PsiReturns) GetCode() string {
	return that.Code
}

func (that *PsiReturns) SetCode(code string) {
	that.Code = code
}

func (that *PsiReturns) GetMessage() string {
	return that.Message
}

func (that *PsiReturns) SetMessage(message string) {
	that.Message = message
}

func (that *PsiReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *PsiReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].(*PsiExecuteResponse)
		}
	}
}

func (that *PsiReturns) GetCause(ctx context.Context) *macro.Cause {
	if nil == that.Cause {
		return nil
	}
	cdc, ok := macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	if !ok {
		return macro.Errors(cause.Errorf("No codec named %s exist ", codec.JSON))
	}
	buff, err := cdc.Encode(that.Cause)
	if nil != err {
		return macro.Errors(err)
	}
	var c macro.Cause
	if _, err := cdc.Decode(buff, &c); nil != err {
		return macro.Errors(err)
	}
	return &c
}

func (that *PsiReturns) SetCause(ctx context.Context, cause *macro.Cause) {

}
