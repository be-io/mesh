/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package test

import (
	_ "github.com/opendatav/mesh/client/golang/grpc"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	_ "github.com/opendatav/mesh/client/golang/proxy"
	"github.com/opendatav/mesh/client/golang/prsim"
	_ "github.com/opendatav/mesh/client/golang/system"
	"testing"
)

func TestNameReplace(t *testing.T) {
	ctx := mpc.Context()
	ctx.SetAttribute(mpc.AddressKey, "10.99.31.33:570")
	runtimes := macro.Load(prsim.IRuntimeAware).List()
	for _, runner := range runtimes {
		ra, ok := runner.(prsim.RuntimeAware)
		if !ok {
			t.Error("F")
			return
		}
		if err := ra.Init(); nil != err {
			t.Error(err)
			return
		}
	}
	dispatcher, ok := macro.Load(prsim.IBuiltin).Get(macro.MeshMPI).(prsim.Builtin)
	if !ok {
		t.Error("No dispatcher")
		return
	}
	ctx.SetAttribute(mpc.RemoteName, "theta")
	x, err := dispatcher.Doc(ctx, "", "")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(x)
}

func TestKVGet(t *testing.T) {
	ctx := mpc.Context()
	ctx.SetAttribute(mpc.AddressKey, "10.99.31.33:570")
	runtimes := macro.Load(prsim.IRuntimeAware).List()
	for _, runner := range runtimes {
		ra, ok := runner.(prsim.RuntimeAware)
		if !ok {
			t.Error("F")
			return
		}
		if err := ra.Init(); nil != err {
			t.Error(err)
			return
		}
	}
	kv, ok := macro.Load(prsim.IKV).Get(macro.MeshMPI).(prsim.KV)
	if !ok {
		t.Error("No KV")
		return
	}
	entity, err := kv.Get(ctx, "mesh.license")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(entity)
}
