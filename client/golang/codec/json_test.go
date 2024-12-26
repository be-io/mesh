/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package codec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opendatav/mesh/client/golang/macro"
	"testing"
	"time"
)

func TestMarshall(t *testing.T) {
	x := map[string]json.RawMessage{}
	x["1"] = json.RawMessage(fmt.Sprintf("\"%s\"", "JG0100000100000000"))
	if buff, err := Jsonizer.Marshal(x); nil != err {
		t.Error(err)
		return
	} else {
		t.Log(string(buff))
	}
}

func TestBinaryCodec(t *testing.T) {
	type TestStruct struct {
		Buff []byte `json:"datetime"`
	}
	cdc := macro.Load(ICodec).Get(JSON).(Codec)
	buff, err := cdc.Encode(&TestStruct{Buff: []byte("mesh</>^;&#$!@*(&)[]{}=''``?/\\")})
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(buff.String())

	var x TestStruct
	if _, err = cdc.Decode(buff, &x); nil != err {
		t.Error(err)
		return
	}
	t.Log(string(x.Buff))
}

func TestTimeCodec(t *testing.T) {
	type TestStruct struct {
		Datetime time.Time `json:"datetime"`
	}
	cdc := macro.Load(ICodec).Get(JSON).(Codec)
	buff, err := cdc.Encode(&TestStruct{Datetime: time.Now()})
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(buff.String())

	var ts TestStruct
	if _, err = cdc.Decode(buff, &ts); nil != err {
		t.Error(err)
		return
	}
	t.Log(ts.Datetime)

	var tt TestStruct
	if _, err = cdc.Decode(bytes.NewBufferString("{\"datetime\":\"2022-02-26 14:16:42\"}"), &tt); nil != err {
		t.Error(err)
		//return
	}
	t.Log(tt.Datetime)

	var tz TestStruct
	if _, err = cdc.Decode(bytes.NewBufferString("{\"datetime\":\"1645856291727\"}"), &tz); nil != err {
		t.Error(err)
		return
	}
	t.Log(tz.Datetime)

}
