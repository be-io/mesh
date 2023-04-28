package test

import (
	"errors"
	"fmt"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/test/proto"
	"github.com/be-io/mesh/client/golang/types"
	"testing"
)

var cdc = macro.Load(codec.ICodec).Get(codec.PROTOBUF4).(codec.Codec)

func TestParametersParser(t *testing.T) {
	// Encode
	parameters := &proto.PsiParameters{
		Attachments: map[string]string{"11": "22"},
		Request: &proto.PsiExecuteRequest{
			SourcePartnerCode: "333",
			TaskId:            "444",
			AuthorityCode:     "555",
			Index:             "666",
			EncodeData:        [][]byte{[]byte("333")},
		},
	}
	encode, err := cdc.Encode(parameters)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Parameters Encode: ", encode)

	// Decode
	var decodeParameters proto.PsiParameters
	decode, err := cdc.Decode(encode, &decodeParameters)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Parameters Decode: ", decode)
}

func TestReturnsParser(t *testing.T) {
	/******************Returns********************/
	// Encode
	returns := &proto.PsiReturns{
		Code:    "000",
		Message: "111",
		Content: &proto.PsiExecuteResponse{
			Stage:    "222",
			TaskId:   "333",
			MaskData: [][]byte{[]byte("333")},
			Encoders: [][]byte{[]byte("333")},
		},
	}
	encode, err := cdc.Encode(returns)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Returns Encode: ", encode)

	// Decode
	var decodeReturns proto.PsiReturns
	decode, err := cdc.Decode(encode, &decodeReturns)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Returns Decode: ", decode, decodeReturns)

	/******************Outbound********************/
	// Encode
	outbound := &types.Outbound{
		Code:    "000",
		Message: "111",
		Cause:   macro.Errors(errors.New("mock error")),
	}
	outboundEncode, err := cdc.Encode(outbound)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Outbound Encode: ", outboundEncode)

	// Decode
	var decodeOutbound types.Outbound
	decodeOut, err := cdc.Decode(outboundEncode, &decodeOutbound)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Outbound Decode: ", decodeOut, decodeOutbound)
}
