/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
)

func init() {
	var _ macro.Inspector = new(GenericInspector)
}

type GenericInspector struct {
	Name string
	Args map[string]interface{}
}

func (that *GenericInspector) Call(ctx context.Context, proxy interface{}, method macro.Inspector, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (that *GenericInspector) String() string {
	return that.Name
}

func (that *GenericInspector) Signature() string {
	return that.Name
}

func (that *GenericInspector) GetDeclaredKind() interface{} {
	return (*interface{})(nil)
}

func (that *GenericInspector) GetName() string {
	return that.Name
}

func (that *GenericInspector) GetTName() string {
	return that.Name
}

func (that *GenericInspector) GetIntype() macro.Parameters {
	parameters := &GenericParameters{}
	return parameters
}

func (that *GenericInspector) GetRetype() macro.Returns {
	var returns GenericReturns
	return &returns
}

func (that *GenericInspector) NewInbound() macro.Parameters {
	return GenericParameters{}
}

func (that *GenericInspector) NewOutbound() macro.Returns {
	return GenericReturns{}
}

func (that *GenericInspector) GetMPS() macro.MPS {
	return &macro.MPSAnnotation{Meta: &macro.Stt{Name: that.Name}}
}

func (that *GenericInspector) GetMPI() macro.MPI {
	return &macro.MPIAnnotation{Meta: &macro.Rtt{Name: that.Name}}
}

func (that *GenericInspector) GetSPI() macro.SPI {
	return &macro.SPIAnnotation{Attribute: &macro.Att{Name: that.Name}}
}

func (that *GenericInspector) GetBinding() []macro.Binding {
	return []macro.Binding{&macro.BindingAnnotation{Binding: &macro.Btt{Topic: that.Name}}}
}
