/*
* Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
* TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
*
*
 */

// Code generated by mesh; DO NOT EDIT.

package proxy

import (
	"context"

	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ macro.Interface = new(meshGraphicsMPI)
	macro.Provide((*prsim.Graphics)(nil), &meshGraphicsMPI{
		invoker: mpc.ServiceProxy.Reference(&macro.Rtt{}),
		methods: map[string]*macro.Method{
			"Captcha": {
				DeclaredKind: (*prsim.Graphics)(nil),
				TName:        "prsim.Graphics",
				Name:         "Captcha",
				Intype:       func() macro.Parameters { var parameters MeshGraphicsCaptchaParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshGraphicsCaptchaReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshGraphicsCaptchaParameters) },
				Outbound:     func() macro.Returns { return new(MeshGraphicsCaptchaReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Name: "mesh.graphics.captcha.apply"},
				},
			},
			"Verify": {
				DeclaredKind: (*prsim.Graphics)(nil),
				TName:        "prsim.Graphics",
				Name:         "Verify",
				Intype:       func() macro.Parameters { var parameters MeshGraphicsVerifyParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshGraphicsVerifyReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshGraphicsVerifyParameters) },
				Outbound:     func() macro.Returns { return new(MeshGraphicsVerifyReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Name: "mesh.graphics.captcha.verify"},
				},
			},
		},
	})
}

// meshGraphicsMPI is an implementation of Graphics
type meshGraphicsMPI struct {
	invoker macro.Caller
	methods map[string]*macro.Method
}

func (that *meshGraphicsMPI) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshMPI}
}

func (that *meshGraphicsMPI) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: macro.MeshMPI}
}

func (that *meshGraphicsMPI) GetMethods() map[string]*macro.Method {
	return that.methods
}

// Captcha apply a graphics captcha.
// @MPI("mesh.graphics.captcha.apply")
func (that *meshGraphicsMPI) Captcha(ctx context.Context, kind string, features map[string]string) (*types.Captcha, error) {
	ret, err := that.invoker.Call(ctx, that.invoker, that.methods["Captcha"], kind, features)
	if nil != err {
		x := new(MeshGraphicsCaptchaReturns)
		return x.Content, err
	}
	x, ok := ret.(*MeshGraphicsCaptchaReturns)
	if ok {
		return x.Content, err
	}
	x = new(MeshGraphicsCaptchaReturns)
	return x.Content, cause.Errorf("Cant resolve response ")
}

// Verify a graphics captcha value.
// @MPI("mesh.graphics.captcha.verify")
func (that *meshGraphicsMPI) Verify(ctx context.Context, mno string, value string) (bool, error) {
	ret, err := that.invoker.Call(ctx, that.invoker, that.methods["Verify"], mno, value)
	if nil != err {
		x := new(MeshGraphicsVerifyReturns)
		return x.Content, err
	}
	x, ok := ret.(*MeshGraphicsVerifyReturns)
	if ok {
		return x.Content, err
	}
	x = new(MeshGraphicsVerifyReturns)
	return x.Content, cause.Errorf("Cant resolve response ")
}

type MeshGraphicsCaptchaParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	Kind        string            `index:"0" json:"kind" xml:"kind" yaml:"kind"`
	Features    map[string]string `index:"1" json:"features" xml:"features" yaml:"features"`
}

func (that *MeshGraphicsCaptchaParameters) GetKind() interface{} {
	return new(MeshGraphicsCaptchaParameters)
}

func (that *MeshGraphicsCaptchaParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Kind)
	arguments = append(arguments, that.Features)
	return arguments
}

func (that *MeshGraphicsCaptchaParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Kind = arguments[0].(string)
		}
		if len(arguments) > 1 && nil != arguments[1] {
			that.Features = arguments[1].(map[string]string)
		}
	}
}

func (that *MeshGraphicsCaptchaParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshGraphicsCaptchaParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshGraphicsCaptchaReturns struct {
	Code    string         `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string         `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause   `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content *types.Captcha `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *MeshGraphicsCaptchaReturns) GetCode() string {
	return that.Code
}

func (that *MeshGraphicsCaptchaReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshGraphicsCaptchaReturns) GetMessage() string {
	return that.Message
}

func (that *MeshGraphicsCaptchaReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshGraphicsCaptchaReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshGraphicsCaptchaReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshGraphicsCaptchaReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *MeshGraphicsCaptchaReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].(*types.Captcha)
		}
	}
}

type MeshGraphicsVerifyParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	Mno         string            `index:"0" json:"mno" xml:"mno" yaml:"mno"`
	Value       string            `index:"1" json:"value" xml:"value" yaml:"value"`
}

func (that *MeshGraphicsVerifyParameters) GetKind() interface{} {
	return new(MeshGraphicsVerifyParameters)
}

func (that *MeshGraphicsVerifyParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Mno)
	arguments = append(arguments, that.Value)
	return arguments
}

func (that *MeshGraphicsVerifyParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Mno = arguments[0].(string)
		}
		if len(arguments) > 1 && nil != arguments[1] {
			that.Value = arguments[1].(string)
		}
	}
}

func (that *MeshGraphicsVerifyParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshGraphicsVerifyParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshGraphicsVerifyReturns struct {
	Code    string       `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content bool         `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *MeshGraphicsVerifyReturns) GetCode() string {
	return that.Code
}

func (that *MeshGraphicsVerifyReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshGraphicsVerifyReturns) GetMessage() string {
	return that.Message
}

func (that *MeshGraphicsVerifyReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshGraphicsVerifyReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshGraphicsVerifyReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshGraphicsVerifyReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *MeshGraphicsVerifyReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].(bool)
		}
	}
}
