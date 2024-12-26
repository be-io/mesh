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
	var _ macro.Interface = new(meshLicenserMPI)
	macro.Provide((*prsim.Licenser)(nil), &meshLicenserMPI{
		invoker: mpc.ServiceProxy.Reference(&macro.Rtt{}),
		methods: map[string]*macro.Method{
			"Imports": {
				DeclaredKind: (*prsim.Licenser)(nil),
				TName:        "prsim.Licenser",
				Name:         "Imports",
				Intype:       func() macro.Parameters { var parameters MeshLicenserImportsParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshLicenserImportsReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshLicenserImportsParameters) },
				Outbound:     func() macro.Returns { return new(MeshLicenserImportsReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Flags: 2, Name: "mesh.license.imports"},
				},
			},
			"Exports": {
				DeclaredKind: (*prsim.Licenser)(nil),
				TName:        "prsim.Licenser",
				Name:         "Exports",
				Intype:       func() macro.Parameters { var parameters MeshLicenserExportsParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshLicenserExportsReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshLicenserExportsParameters) },
				Outbound:     func() macro.Returns { return new(MeshLicenserExportsReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Flags: 2, Name: "mesh.license.exports"},
				},
			},
			"Explain": {
				DeclaredKind: (*prsim.Licenser)(nil),
				TName:        "prsim.Licenser",
				Name:         "Explain",
				Intype:       func() macro.Parameters { var parameters MeshLicenserExplainParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshLicenserExplainReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshLicenserExplainParameters) },
				Outbound:     func() macro.Returns { return new(MeshLicenserExplainReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Flags: 2, Name: "mesh.license.explain"},
				},
			},
			"Verify": {
				DeclaredKind: (*prsim.Licenser)(nil),
				TName:        "prsim.Licenser",
				Name:         "Verify",
				Intype:       func() macro.Parameters { var parameters MeshLicenserVerifyParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshLicenserVerifyReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshLicenserVerifyParameters) },
				Outbound:     func() macro.Returns { return new(MeshLicenserVerifyReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Flags: 2, Name: "mesh.license.verify"},
				},
			},
			"Features": {
				DeclaredKind: (*prsim.Licenser)(nil),
				TName:        "prsim.Licenser",
				Name:         "Features",
				Intype:       func() macro.Parameters { var parameters MeshLicenserFeaturesParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshLicenserFeaturesReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshLicenserFeaturesParameters) },
				Outbound:     func() macro.Returns { return new(MeshLicenserFeaturesReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Flags: 2, Name: "mesh.license.features"},
				},
			},
		},
	})
}

// meshLicenserMPI is an implementation of Licenser
type meshLicenserMPI struct {
	invoker macro.Caller
	methods map[string]*macro.Method
}

func (that *meshLicenserMPI) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshMPI}
}

func (that *meshLicenserMPI) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: macro.MeshMPI}
}

func (that *meshLicenserMPI) GetMethods() map[string]*macro.Method {
	return that.methods
}

// Imports the licenses.
// @MPI(name = "mesh.license.imports", flags = 2)
func (that *meshLicenserMPI) Imports(ctx context.Context, license string) error {
	_, err := that.invoker.Call(ctx, that.invoker, that.methods["Imports"], license)
	return err
}

// Exports the licenses.
// @MPI(name = "mesh.license.exports", flags = 2)
func (that *meshLicenserMPI) Exports(ctx context.Context) (string, error) {
	ret, err := that.invoker.Call(ctx, that.invoker, that.methods["Exports"])
	if nil != err {
		x := new(MeshLicenserExportsReturns)
		return x.Content, err
	}
	x, ok := ret.(*MeshLicenserExportsReturns)
	if ok {
		return x.Content, err
	}
	x = new(MeshLicenserExportsReturns)
	return x.Content, cause.Errorf("Cant resolve response ")
}

// Explain the license.
// @MPI(name = "mesh.license.explain", flags = 2)
func (that *meshLicenserMPI) Explain(ctx context.Context) (*types.License, error) {
	ret, err := that.invoker.Call(ctx, that.invoker, that.methods["Explain"])
	if nil != err {
		x := new(MeshLicenserExplainReturns)
		return x.Content, err
	}
	x, ok := ret.(*MeshLicenserExplainReturns)
	if ok {
		return x.Content, err
	}
	x = new(MeshLicenserExplainReturns)
	return x.Content, cause.Errorf("Cant resolve response ")
}

// Verify the license.
// @MPI(name = "mesh.license.verify", flags = 2)
func (that *meshLicenserMPI) Verify(ctx context.Context) (int64, error) {
	ret, err := that.invoker.Call(ctx, that.invoker, that.methods["Verify"])
	if nil != err {
		x := new(MeshLicenserVerifyReturns)
		return x.Content, err
	}
	x, ok := ret.(*MeshLicenserVerifyReturns)
	if ok {
		return x.Content, err
	}
	x = new(MeshLicenserVerifyReturns)
	return x.Content, cause.Errorf("Cant resolve response ")
}

// Features is license features.
// @MPI(name = "mesh.license.features", flags = 2)
func (that *meshLicenserMPI) Features(ctx context.Context) (map[string]string, error) {
	ret, err := that.invoker.Call(ctx, that.invoker, that.methods["Features"])
	if nil != err {
		x := new(MeshLicenserFeaturesReturns)
		return x.Content, err
	}
	x, ok := ret.(*MeshLicenserFeaturesReturns)
	if ok {
		return x.Content, err
	}
	x = new(MeshLicenserFeaturesReturns)
	return x.Content, cause.Errorf("Cant resolve response ")
}

type MeshLicenserImportsParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	License     string            `index:"0" json:"license" xml:"license" yaml:"license"`
}

func (that *MeshLicenserImportsParameters) GetKind() interface{} {
	return new(MeshLicenserImportsParameters)
}

func (that *MeshLicenserImportsParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.License)
	return arguments
}

func (that *MeshLicenserImportsParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.License = arguments[0].(string)
		}
	}
}

func (that *MeshLicenserImportsParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshLicenserImportsParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshLicenserImportsReturns struct {
	Code    string       `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
}

func (that *MeshLicenserImportsReturns) GetCode() string {
	return that.Code
}

func (that *MeshLicenserImportsReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshLicenserImportsReturns) GetMessage() string {
	return that.Message
}

func (that *MeshLicenserImportsReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshLicenserImportsReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshLicenserImportsReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshLicenserImportsReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshLicenserImportsReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

type MeshLicenserExportsParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
}

func (that *MeshLicenserExportsParameters) GetKind() interface{} {
	return new(MeshLicenserExportsParameters)
}

func (that *MeshLicenserExportsParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshLicenserExportsParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

func (that *MeshLicenserExportsParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshLicenserExportsParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshLicenserExportsReturns struct {
	Code    string       `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content string       `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *MeshLicenserExportsReturns) GetCode() string {
	return that.Code
}

func (that *MeshLicenserExportsReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshLicenserExportsReturns) GetMessage() string {
	return that.Message
}

func (that *MeshLicenserExportsReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshLicenserExportsReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshLicenserExportsReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshLicenserExportsReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *MeshLicenserExportsReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].(string)
		}
	}
}

type MeshLicenserExplainParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
}

func (that *MeshLicenserExplainParameters) GetKind() interface{} {
	return new(MeshLicenserExplainParameters)
}

func (that *MeshLicenserExplainParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshLicenserExplainParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

func (that *MeshLicenserExplainParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshLicenserExplainParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshLicenserExplainReturns struct {
	Code    string         `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string         `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause   `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content *types.License `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *MeshLicenserExplainReturns) GetCode() string {
	return that.Code
}

func (that *MeshLicenserExplainReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshLicenserExplainReturns) GetMessage() string {
	return that.Message
}

func (that *MeshLicenserExplainReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshLicenserExplainReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshLicenserExplainReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshLicenserExplainReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *MeshLicenserExplainReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].(*types.License)
		}
	}
}

type MeshLicenserVerifyParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
}

func (that *MeshLicenserVerifyParameters) GetKind() interface{} {
	return new(MeshLicenserVerifyParameters)
}

func (that *MeshLicenserVerifyParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshLicenserVerifyParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

func (that *MeshLicenserVerifyParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshLicenserVerifyParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshLicenserVerifyReturns struct {
	Code    string       `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content int64        `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *MeshLicenserVerifyReturns) GetCode() string {
	return that.Code
}

func (that *MeshLicenserVerifyReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshLicenserVerifyReturns) GetMessage() string {
	return that.Message
}

func (that *MeshLicenserVerifyReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshLicenserVerifyReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshLicenserVerifyReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshLicenserVerifyReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *MeshLicenserVerifyReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].(int64)
		}
	}
}

type MeshLicenserFeaturesParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
}

func (that *MeshLicenserFeaturesParameters) GetKind() interface{} {
	return new(MeshLicenserFeaturesParameters)
}

func (that *MeshLicenserFeaturesParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshLicenserFeaturesParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

func (that *MeshLicenserFeaturesParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshLicenserFeaturesParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshLicenserFeaturesReturns struct {
	Code    string            `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string            `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause      `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content map[string]string `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *MeshLicenserFeaturesReturns) GetCode() string {
	return that.Code
}

func (that *MeshLicenserFeaturesReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshLicenserFeaturesReturns) GetMessage() string {
	return that.Message
}

func (that *MeshLicenserFeaturesReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshLicenserFeaturesReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshLicenserFeaturesReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshLicenserFeaturesReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *MeshLicenserFeaturesReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].(map[string]string)
		}
	}
}
