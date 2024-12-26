/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import (
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
)

func init() {
	var _ Inspector = new(Method)
}

var ICaller = (*Caller)(nil)

type Interface interface {
	GetMethods() map[string]*Method
}

// Inspector will inspec the method signature.
type Inspector interface {

	// Stringer is inspector name.
	fmt.Stringer

	// Caller is inspector invocation.
	Caller

	// Signature return inspector signature
	Signature() string

	// GetDeclaredKind return the declared kind
	GetDeclaredKind() interface{}

	// GetName return method name
	GetName() string

	// GetTName return type name
	GetTName() string

	// GetIntype is input arguments
	GetIntype() Parameters

	// GetRetype is output arguments
	GetRetype() Returns

	// NewInbound return new inbound.
	NewInbound() Parameters

	// NewOutbound return new outbound.
	NewOutbound() Returns

	GetMPS() MPS

	GetMPI() MPI

	GetSPI() SPI

	GetBinding() []Binding
}

// Parameters is the method argument array type
type Parameters interface {
	GetKind() interface{}
	GetArguments(ctx context.Context) []interface{}
	SetArguments(ctx context.Context, arguments ...interface{})
	GetAttachments(ctx context.Context) map[string]string
	SetAttachments(ctx context.Context, attachments map[string]string)
}

type Returns interface {
	GetCode() string
	SetCode(code string)
	GetMessage() string
	SetMessage(message string)
	GetContent(ctx context.Context) []interface{}
	SetContent(ctx context.Context, arguments ...interface{})
	GetCause(ctx context.Context) *Cause
	SetCause(ctx context.Context, cause *Cause)
}

type Caller interface {

	// Call will execute the method
	Call(ctx context.Context, proxy interface{}, method Inspector, args ...interface{}) (interface{}, error)
}

type Method struct {
	DeclaredKind interface{}
	TName        string
	Name         string
	MPI          MPI
	MPS          MPS
	SPI          SPI
	Bindings     []Binding
	Intype       func() Parameters
	Retype       func() Returns
	Inbound      func() Parameters
	Outbound     func() Returns
}

func (that *Method) String() string {
	return fmt.Sprintf("%p", that)
}

func (that *Method) GetMPS() MPS {
	return that.MPS
}

func (that *Method) GetMPI() MPI {
	return that.MPI
}

func (that *Method) GetSPI() SPI {
	return that.SPI
}

func (that *Method) GetBinding() []Binding {
	return that.Bindings
}

func (that *Method) Signature() string {
	return that.Name
}

func (that *Method) GetDeclaredKind() interface{} {
	return that.DeclaredKind
}

func (that *Method) GetName() string {
	return that.Name
}

func (that *Method) GetTName() string {
	return that.TName
}

func (that *Method) GetIntype() Parameters {
	return that.Intype()
}

func (that *Method) GetRetype() Returns {
	return that.Retype()
}

func (that *Method) NewInbound() Parameters {
	return that.Inbound()
}

func (that *Method) NewOutbound() Returns {
	return that.Outbound()
}

func (that *Method) Call(ctx context.Context, proxy interface{}, method Inspector, args ...interface{}) (interface{}, error) {
	caller, ok := Load(ICaller).Get(method.String()).(Caller)
	if !ok {
		return nil, cause.NotFoundErrorf("No service named %s exist ", method.String())
	}
	return caller.Call(ctx, proxy, method, args...)
}
