/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"bytes"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"testing"
)

func TestGrpcs(t *testing.T) {
	consumer, ok := macro.Load(mpc.IConsumer).Get(Name).(mpc.Consumer)
	if !ok {
		t.Error(cause.Errorf("No consumer named %s exist. ", Name))
		return
	}
	mtx := mpc.Context()
	mtx.SetAttribute(mpc.AddressKey, "127.0.0.1:8864")
	mtx.SetAttribute(mpc.InsecureKey, true)
	buff, err := consumer.Consume(mtx, "version.net.mesh.0001000000000000000000000000000000000.JG0100000100000000.trustbe.cn", nil, bytes.NewBufferString("{\"node_id\":\"JG0100000200000000\"}"))
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(buff.String())

	buff, err = consumer.Consume(mtx, "transparent.net.mesh.0001000000000000000000000000000000000.JG0100000200000000.trustbe.cn", nil, bytes.NewBufferString("{\"node_id\":\"JG0100000200000000\"}"))
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(buff.String())

}
