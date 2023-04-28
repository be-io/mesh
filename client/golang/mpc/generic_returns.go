/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"strings"
)

func init() {
	var _ macro.Returns = new(GenericReturns)
}

type GenericReturns map[string]interface{}

func (that GenericReturns) GetCode() string {
	if nil == that["code"] {
		return ""
	}
	return fmt.Sprintf("%v", that["code"])
}

func (that GenericReturns) SetCode(code string) {
	that["code"] = code
}

func (that GenericReturns) GetMessage() string {
	if nil == that["message"] {
		return ""
	}
	return fmt.Sprintf("%v", that["message"])
}

func (that GenericReturns) SetMessage(message string) {
	that["message"] = message
}

func (that GenericReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	for key, value := range that {
		if strings.Index(key, "content") > -1 {
			arguments = append(arguments, value)
		}
	}
	return arguments
}

func (that GenericReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if nil != arguments {
		for index, argument := range arguments {
			if 0 == index {
				that["content"] = argument
			} else {
				that[fmt.Sprintf("content%d", index)] = argument
			}
		}
	}
}

func (that GenericReturns) GetCause(ctx context.Context) *macro.Cause {
	if nil == that["cause"] {
		return nil
	}
	if c, ok := that["cause"].(*macro.Cause); ok {
		return c
	}
	cdc, ok := macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	if !ok {
		return macro.Errors(cause.Errorf("No codec named %s exist ", codec.JSON))
	}
	buff, err := cdc.Encode(that["cause"])
	if nil != err {
		return macro.Errors(err)
	}
	var c macro.Cause
	if _, err := cdc.Decode(buff, &c); nil != err {
		return macro.Errors(err)
	}
	return &c
}

func (that GenericReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that["cause"] = cause
}
