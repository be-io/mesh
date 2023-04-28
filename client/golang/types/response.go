/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import "github.com/be-io/mesh/client/golang/cause"

type Response struct {
	Code    string      `index:"0" json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty"`
	Content interface{} `index:"1" json:"content,omitempty" yaml:"content,omitempty" xml:"content,omitempty"`
	Message string      `index:"2" json:"message,omitempty" yaml:"message,omitempty" xml:"message,omitempty"`
	Version string      `index:"3" json:"version,omitempty" yaml:"version,omitempty" xml:"version,omitempty"`
}

func (that Response) IsOK() bool {
	return that.Code == cause.Success.Code
}

func OK() Response {
	return Response{Code: cause.Success.Code}
}

func OK0(data interface{}) Response {
	return Response{Code: cause.Success.Code, Content: data}
}

func NOTOK(message string) Response {
	return Response{Code: cause.SystemError.Code, Message: message}
}
