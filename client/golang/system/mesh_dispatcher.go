/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.Dispatcher = new(MeshDispatcher)
	macro.Provide(prsim.IDispatcher, &MeshDispatcher{
		invoker: &mpc.GenericHandler{},
	})
}

type MeshDispatcher struct {
	invoker *mpc.GenericHandler
}

func (that *MeshDispatcher) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSPI}
}

func (that *MeshDispatcher) Invoke(ctx context.Context, urn string, param map[string]interface{}) ([]interface{}, error) {
	return that.invoker.Invoke00(ctx, urn, param)
}

func (that *MeshDispatcher) Invoke0(ctx context.Context, urn string, param interface{}) ([]interface{}, error) {
	cdc, ok := macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	if !ok {
		return nil, cause.Errorf("No Codec named %s exist ", codec.JSON)
	}
	buff, err := cdc.Encode(param)
	if nil != err {
		return nil, cause.Error(err)
	}
	var kvs map[string]interface{}
	if _, err = cdc.Decode(buff, &kvs); nil != err {
		return nil, cause.Error(err)
	}
	return that.invoker.Invoke00(ctx, urn, kvs)
}

func (that *MeshDispatcher) InvokeLR(ctx context.Context, urn string, param map[string]interface{}) (interface{}, error) {
	rets, err := that.invoker.Invoke00(ctx, urn, param)
	if len(rets) > 0 {
		return rets[0], err
	}
	return nil, err
}

func (that *MeshDispatcher) InvokeLRG(ctx context.Context, urn string, param interface{}) (interface{}, error) {
	rets, err := that.Invoke0(ctx, urn, param)
	if len(rets) > 0 {
		return rets[0], err
	}
	return nil, err
}
