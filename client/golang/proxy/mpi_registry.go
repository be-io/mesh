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
	var _ macro.Interface = new(meshRegistryMPI)
	macro.Provide((*prsim.Registry)(nil), &meshRegistryMPI{
		invoker: mpc.ServiceProxy.Reference(&macro.Rtt{}),
		methods: map[string]*macro.Method{
			"Register": {
				DeclaredKind: (*prsim.Registry)(nil),
				TName:        "prsim.Registry",
				Name:         "Register",
				Intype:       func() macro.Parameters { var parameters MeshRegistryRegisterParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshRegistryRegisterReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshRegistryRegisterParameters) },
				Outbound:     func() macro.Returns { return new(MeshRegistryRegisterReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Name: "mesh.registry.put"},
				},
			},
			"Registers": {
				DeclaredKind: (*prsim.Registry)(nil),
				TName:        "prsim.Registry",
				Name:         "Registers",
				Intype:       func() macro.Parameters { var parameters MeshRegistryRegistersParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshRegistryRegistersReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshRegistryRegistersParameters) },
				Outbound:     func() macro.Returns { return new(MeshRegistryRegistersReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Name: "mesh.registry.puts"},
				},
			},
			"Unregister": {
				DeclaredKind: (*prsim.Registry)(nil),
				TName:        "prsim.Registry",
				Name:         "Unregister",
				Intype:       func() macro.Parameters { var parameters MeshRegistryUnregisterParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshRegistryUnregisterReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshRegistryUnregisterParameters) },
				Outbound:     func() macro.Returns { return new(MeshRegistryUnregisterReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Name: "mesh.registry.remove"},
				},
			},
			"Export": {
				DeclaredKind: (*prsim.Registry)(nil),
				TName:        "prsim.Registry",
				Name:         "Export",
				Intype:       func() macro.Parameters { var parameters MeshRegistryExportParameters; return &parameters },
				Retype:       func() macro.Returns { var returns MeshRegistryExportReturns; return &returns },
				Inbound:      func() macro.Parameters { return new(MeshRegistryExportParameters) },
				Outbound:     func() macro.Returns { return new(MeshRegistryExportReturns) },
				MPI: &macro.MPIAnnotation{
					Meta: &macro.Rtt{Name: "mesh.registry.export"},
				},
			},
		},
	})
}

// meshRegistryMPI is an implementation of Registry
type meshRegistryMPI struct {
	invoker macro.Caller
	methods map[string]*macro.Method
}

func (that *meshRegistryMPI) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshMPI}
}

func (that *meshRegistryMPI) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: macro.MeshMPI}
}

func (that *meshRegistryMPI) GetMethods() map[string]*macro.Method {
	return that.methods
}

// Register
// @MPI("mesh.registry.put")
func (that *meshRegistryMPI) Register(ctx context.Context, registration *types.Registration[any]) error {
	_, err := that.invoker.Call(ctx, that.invoker, that.methods["Register"], registration)
	return err
}

// Registers
// @MPI("mesh.registry.puts")
func (that *meshRegistryMPI) Registers(ctx context.Context, registrations []*types.Registration[any]) error {
	_, err := that.invoker.Call(ctx, that.invoker, that.methods["Registers"], registrations)
	return err
}

// Unregister
// @MPI("mesh.registry.remove")
func (that *meshRegistryMPI) Unregister(ctx context.Context, registration *types.Registration[any]) error {
	_, err := that.invoker.Call(ctx, that.invoker, that.methods["Unregister"], registration)
	return err
}

// Export
// @MPI("mesh.registry.export")
func (that *meshRegistryMPI) Export(ctx context.Context, kind string) ([]*types.Registration[any], error) {
	ret, err := that.invoker.Call(ctx, that.invoker, that.methods["Export"], kind)
	if nil != err {
		x := new(MeshRegistryExportReturns)
		return x.Content, err
	}
	x, ok := ret.(*MeshRegistryExportReturns)
	if ok {
		return x.Content, err
	}
	x = new(MeshRegistryExportReturns)
	return x.Content, cause.Errorf("Cant resolve response ")
}

type MeshRegistryRegisterParameters struct {
	Attachments  map[string]string        `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	Registration *types.Registration[any] `index:"0" json:"registration" xml:"registration" yaml:"registration"`
}

func (that *MeshRegistryRegisterParameters) GetKind() interface{} {
	return new(MeshRegistryRegisterParameters)
}

func (that *MeshRegistryRegisterParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Registration)
	return arguments
}

func (that *MeshRegistryRegisterParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Registration = arguments[0].(*types.Registration[any])
		}
	}
}

func (that *MeshRegistryRegisterParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshRegistryRegisterParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshRegistryRegisterReturns struct {
	Code    string       `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
}

func (that *MeshRegistryRegisterReturns) GetCode() string {
	return that.Code
}

func (that *MeshRegistryRegisterReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshRegistryRegisterReturns) GetMessage() string {
	return that.Message
}

func (that *MeshRegistryRegisterReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshRegistryRegisterReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshRegistryRegisterReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshRegistryRegisterReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshRegistryRegisterReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

type MeshRegistryRegistersParameters struct {
	Attachments   map[string]string          `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	Registrations []*types.Registration[any] `index:"0" json:"registrations" xml:"registrations" yaml:"registrations"`
}

func (that *MeshRegistryRegistersParameters) GetKind() interface{} {
	return new(MeshRegistryRegistersParameters)
}

func (that *MeshRegistryRegistersParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Registrations)
	return arguments
}

func (that *MeshRegistryRegistersParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Registrations = arguments[0].([]*types.Registration[any])
		}
	}
}

func (that *MeshRegistryRegistersParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshRegistryRegistersParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshRegistryRegistersReturns struct {
	Code    string       `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
}

func (that *MeshRegistryRegistersReturns) GetCode() string {
	return that.Code
}

func (that *MeshRegistryRegistersReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshRegistryRegistersReturns) GetMessage() string {
	return that.Message
}

func (that *MeshRegistryRegistersReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshRegistryRegistersReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshRegistryRegistersReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshRegistryRegistersReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshRegistryRegistersReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

type MeshRegistryUnregisterParameters struct {
	Attachments  map[string]string        `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	Registration *types.Registration[any] `index:"0" json:"registration" xml:"registration" yaml:"registration"`
}

func (that *MeshRegistryUnregisterParameters) GetKind() interface{} {
	return new(MeshRegistryUnregisterParameters)
}

func (that *MeshRegistryUnregisterParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Registration)
	return arguments
}

func (that *MeshRegistryUnregisterParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Registration = arguments[0].(*types.Registration[any])
		}
	}
}

func (that *MeshRegistryUnregisterParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshRegistryUnregisterParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshRegistryUnregisterReturns struct {
	Code    string       `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string       `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
}

func (that *MeshRegistryUnregisterReturns) GetCode() string {
	return that.Code
}

func (that *MeshRegistryUnregisterReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshRegistryUnregisterReturns) GetMessage() string {
	return that.Message
}

func (that *MeshRegistryUnregisterReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshRegistryUnregisterReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshRegistryUnregisterReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshRegistryUnregisterReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	return arguments
}

func (that *MeshRegistryUnregisterReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
	}
}

type MeshRegistryExportParameters struct {
	Attachments map[string]string `index:"-1" json:"attachments" xml:"attachments" yaml:"attachments"`
	Kind        string            `index:"0" json:"kind" xml:"kind" yaml:"kind"`
}

func (that *MeshRegistryExportParameters) GetKind() interface{} {
	return new(MeshRegistryExportParameters)
}

func (that *MeshRegistryExportParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Kind)
	return arguments
}

func (that *MeshRegistryExportParameters) SetArguments(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Kind = arguments[0].(string)
		}
	}
}

func (that *MeshRegistryExportParameters) GetAttachments(ctx context.Context) map[string]string {
	return that.Attachments
}

func (that *MeshRegistryExportParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that.Attachments = attachments
}

type MeshRegistryExportReturns struct {
	Code    string                     `index:"0" json:"code" xml:"code" yaml:"code" comment:"Result code"`
	Message string                     `index:"5" json:"message" xml:"message" yaml:"message" comment:"Result message"`
	Cause   *macro.Cause               `index:"10" json:"cause" xml:"cause" yaml:"cause" comment:"Service cause stacktrace"`
	Content []*types.Registration[any] `index:"15" json:"content" xml:"content" yaml:"content"`
}

func (that *MeshRegistryExportReturns) GetCode() string {
	return that.Code
}

func (that *MeshRegistryExportReturns) SetCode(code string) {
	that.Code = code
}

func (that *MeshRegistryExportReturns) GetMessage() string {
	return that.Message
}

func (that *MeshRegistryExportReturns) SetMessage(message string) {
	that.Message = message
}

func (that *MeshRegistryExportReturns) GetCause(ctx context.Context) *macro.Cause {
	return that.Cause
}

func (that *MeshRegistryExportReturns) SetCause(ctx context.Context, cause *macro.Cause) {
	that.Cause = cause
}

func (that *MeshRegistryExportReturns) GetContent(ctx context.Context) []interface{} {
	var arguments []interface{}
	arguments = append(arguments, that.Content)
	return arguments
}

func (that *MeshRegistryExportReturns) SetContent(ctx context.Context, arguments ...interface{}) {
	if len(arguments) > 0 {
		if len(arguments) > 0 && nil != arguments[0] {
			that.Content = arguments[0].([]*types.Registration[any])
		}
	}
}
