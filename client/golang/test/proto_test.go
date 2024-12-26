/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package test

import (
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/macro"
	"testing"
)

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
//go:generate go run ../proto/generate.go -m github.com/opendatav/mesh/client/golang/test -j ../../java/src/test/java

type PureBody struct {
	A float64              `index:"1" json:"a,omitempty"`
	B float32              `index:"2" json:"b,omitempty"`
	C int32                `index:"3" json:"c,omitempty"`
	D int64                `index:"4" json:"d,omitempty"`
	E uint32               `index:"5" json:"e,omitempty"`
	F uint64               `index:"6" json:"f,omitempty"`
	G int32                `index:"7" json:"g,omitempty"`
	H int64                `index:"8" json:"h,omitempty"`
	I uint32               `index:"9" json:"i,omitempty"`
	J uint64               `index:"10" json:"j,omitempty"`
	K int32                `index:"11" json:"k,omitempty"`
	L int64                `index:"12" json:"l,omitempty"`
	M bool                 `index:"13" json:"m,omitempty"`
	N string               `index:"14" json:"n,omitempty"`
	O []byte               `index:"15" json:"o,omitempty"`
	P map[string]string    `index:"16" json:"p,omitempty"`
	Q []int32              `index:"17" json:"q,omitempty"`
	R *PureBody            `index:"18" json:"r,omitempty"`
	S map[string]*PureBody `index:"19" json:"s,omitempty"`
	T []*PureBody          `index:"20" json:"t,omitempty"`
}

type PureMessage struct {
	A float64              `index:"1" json:"a,omitempty"`
	B float32              `index:"2" json:"b,omitempty"`
	C int32                `index:"3" json:"c,omitempty"`
	D int64                `index:"4" json:"d,omitempty"`
	E uint32               `index:"5" json:"e,omitempty"`
	F uint64               `index:"6" json:"f,omitempty"`
	G int32                `index:"7" json:"g,omitempty"`
	H int64                `index:"8" json:"h,omitempty"`
	I uint32               `index:"9" json:"i,omitempty"`
	J uint64               `index:"10" json:"j,omitempty"`
	K int32                `index:"11" json:"k,omitempty"`
	L int64                `index:"12" json:"l,omitempty"`
	M bool                 `index:"13" json:"m,omitempty"`
	N string               `index:"14" json:"n,omitempty"`
	O []byte               `index:"15" json:"o,omitempty"`
	P map[string]string    `index:"16" json:"p,omitempty"`
	Q []int32              `index:"17" json:"q,omitempty"`
	R *PureBody            `index:"18" json:"r,omitempty"`
	S map[string]*PureBody `index:"19" json:"s,omitempty"`
	T []*PureBody          `index:"20" json:"t,omitempty"`
}

func TestProtoCodec(t *testing.T) {
	message := &PureMessage{
		A: 1.1,
		B: 2.1,
		C: 3,
		D: 4,
		E: 5,
		F: 6,
		G: 7,
		H: 8,
		I: 9,
		J: 10,
		K: 11,
		L: 12,
		M: true,
		N: "xxx",
		O: []byte("xxx"),
		P: map[string]string{"xxx": "yyy"},
		Q: []int32{1, 2, 3, 4, 5},
		R: &PureBody{
			A: 1.1,
			B: 2.1,
			C: 3,
			D: 4,
			E: 5,
			F: 6,
			G: 7,
			H: 8,
			I: 9,
			J: 10,
			K: 11,
			L: 12,
			M: true,
			N: "xxx",
			O: []byte("xxx"),
			P: map[string]string{"xxx": "yyy"},
			Q: []int32{1, 2, 3, 4, 5},
			R: &PureBody{},
			S: map[string]*PureBody{"xxx": {}},
			T: []*PureBody{{}},
		},
		S: map[string]*PureBody{"xxx": {
			A: 1.1,
			B: 2.1,
			C: 3,
			D: 4,
			E: 5,
			F: 6,
			G: 7,
			H: 8,
			I: 9,
			J: 10,
			K: 11,
			L: 12,
			M: true,
			N: "xxx",
			O: []byte("xxx"),
			P: map[string]string{"xxx": "yyy"},
			Q: []int32{1, 2, 3, 4, 5},
			R: &PureBody{},
			S: map[string]*PureBody{"xxx": {}},
			T: []*PureBody{{}},
		}},
		T: []*PureBody{{
			A: 1.1,
			B: 2.1,
			C: 3,
			D: 4,
			E: 5,
			F: 6,
			G: 7,
			H: 8,
			I: 9,
			J: 10,
			K: 11,
			L: 12,
			M: true,
			N: "xxx",
			O: []byte("xxx"),
			P: map[string]string{"xxx": "yyy"},
			Q: []int32{1, 2, 3, 4, 5},
			R: &PureBody{},
			S: map[string]*PureBody{"xxx": {}},
			T: []*PureBody{{}},
		}},
	}
	cdc := macro.Load(codec.ICodec).Get(codec.PROTOBUF).(codec.Codec)
	buff, err := cdc.Encode(message)
	if nil != err {
		t.Error(err)
		return
	}
	var dm PureMessage
	if _, err = cdc.Decode(buff, &dm); nil != err {
		t.Error(err)
		return
	}
	t.Log(dm)
}
