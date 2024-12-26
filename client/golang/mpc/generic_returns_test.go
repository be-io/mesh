/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"encoding/json"
	"github.com/opendatav/mesh/client/golang/codec"
	"testing"
)

func TestGenericReturns(t *testing.T) {
	x := GenericReturns{}
	x.SetCode("1")
	x.SetMessage("xx")
	x.SetContent(Context(), json.RawMessage("{\"x\":\"x\"}"))

	buff, err := codec.Jsonizer.Marshal(x)
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(string(buff))
}
