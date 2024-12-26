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
	"github.com/opendatav/mesh/client/golang/types"
	"time"
)

func init() {
	var _ macro.Inspector = new(InspectExecution)
	var _ Execution = new(InspectExecution)
}

type InspectExecution struct {
	URN *types.URN
}

func (that *InspectExecution) Call(ctx context.Context, proxy interface{}, method macro.Inspector, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (that *InspectExecution) String() string {
	return that.URN.Name
}

func (that *InspectExecution) Signature() string {
	return that.URN.Name
}

func (that *InspectExecution) GetDeclaredKind() interface{} {
	return (*interface{})(nil)
}

func (that *InspectExecution) GetName() string {
	return that.URN.Name
}

func (that *InspectExecution) GetTName() string {
	return that.URN.Name
}

func (that *InspectExecution) GetMPS() macro.MPS {
	return &macro.MPSAnnotation{Meta: &macro.Stt{Name: that.URN.Name}}
}

func (that *InspectExecution) GetMPI() macro.MPI {
	return &macro.MPIAnnotation{Meta: &macro.Rtt{Name: that.URN.Name}}
}

func (that *InspectExecution) GetSPI() macro.SPI {
	return &macro.SPIAnnotation{Attribute: &macro.Att{Name: that.URN.Name}}
}

func (that *InspectExecution) GetBinding() []macro.Binding {
	return []macro.Binding{&macro.BindingAnnotation{Binding: &macro.Btt{Topic: that.URN.Name}}}
}

func (that *InspectExecution) Schema() Generic {
	return &types.Reference{
		URN:       that.URN.String(),
		Namespace: "",
		Name:      that.URN.Name,
		Version:   that.URN.Flag.Version,
		Proto:     MeshFlag.OfProto(that.URN.Flag.Proto).Name(),
		Codec:     MeshFlag.OfCodec(that.URN.Flag.Codec).Name(),
		Flags:     0,
		Timeout:   time.Second.Milliseconds() * 10,
		Retries:   5,
		Node:      that.URN.NodeId,
		Inst:      "",
		Zone:      that.URN.Flag.Zone,
		Cluster:   that.URN.Flag.Cluster,
		Group:     that.URN.Flag.Group,
		Address:   that.URN.Flag.Address,
	}
}

func (that *InspectExecution) Inspect() macro.Inspector {
	return that
}

func (that *InspectExecution) GetIntype() macro.Parameters {
	parameters := &GenericParameters{}
	return parameters
}

func (that *InspectExecution) GetRetype() macro.Returns {
	var returns GenericReturns
	return &returns
}

func (that *InspectExecution) NewInbound() macro.Parameters {
	return GenericParameters{}
}

func (that *InspectExecution) NewOutbound() macro.Returns {
	return GenericReturns{}
}

func (that *InspectExecution) Invoke(ctx context.Context, invocation Invocation) (interface{}, error) {
	invoker := &ServiceHandler{inspector: that}
	return composite(invoker, PROVIDER).Invoke(ctx, invocation)
}
