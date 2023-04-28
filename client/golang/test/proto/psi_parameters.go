/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proto

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
)

func init() {
	var _ macro.Parameters = new(PsiParameters)
}

type PsiParameters struct {
	Attachments map[string]string  `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	Request     *PsiExecuteRequest `index:"0" json:"request" xml:"request" yaml:"request"`
}

func (that *PsiParameters) GetKind() interface{} {
	return (*map[string]interface{})(nil)
}

func (that *PsiParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Request)
	return arguments
}

func (that *PsiParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Request = arguments[0].(*PsiExecuteRequest)
		}
	}
}

func (that *PsiParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *PsiParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}
