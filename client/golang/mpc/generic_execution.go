/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/types"
	"time"
)

func init() {
	var _ Execution = new(GenericExecution)
}

type GenericExecution struct {
	args      map[string]interface{}
	inspector macro.Inspector
	reference Generic
}

func (that *GenericExecution) Init(urn *types.URN, arguments map[string]interface{}) {
	that.reference = &types.Reference{
		URN:       urn.String(),
		Namespace: "",
		Name:      urn.Name,
		Version:   urn.Flag.Version,
		Proto:     MeshFlag.OfProto(urn.Flag.Proto).Name(),
		Codec:     MeshFlag.OfCodec(urn.Flag.Codec).Name(),
		Flags:     0,
		Timeout:   time.Second.Milliseconds() * 10,
		Retries:   5,
		Node:      urn.NodeId,
		Inst:      "",
		Zone:      urn.Flag.Zone,
		Cluster:   urn.Flag.Cluster,
		Group:     urn.Flag.Group,
		Address:   urn.Flag.Address,
	}
	that.args = arguments
	that.inspector = &GenericInspector{Name: that.reference.GetURN(), Args: arguments}
}

func (that *GenericExecution) Schema() Generic {
	return that.reference
}

func (that *GenericExecution) Inspect() macro.Inspector {
	return that.inspector
}

func (that *GenericExecution) Invoke(ctx context.Context, invocation Invocation) (interface{}, error) {
	invoker := &ServiceHandler{inspector: that.inspector}
	return composite(invoker, PROVIDER).Invoke(ctx, invocation)
}
