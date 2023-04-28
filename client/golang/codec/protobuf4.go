package codec

import (
	"bytes"
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec/proto4"
	"github.com/be-io/mesh/client/golang/macro"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func init() {
	macro.Provide(ICodec, new(Protobuf4))
}

const PROTOBUF4 = "protobuf4"

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
//go:generate go run ../proto/generate.go -m github.com/be-io/mesh/client/golang/codec/proto4
type Protobuf4 struct {
}

func (that *Protobuf4) Att() *macro.Att {
	return &macro.Att{Name: PROTOBUF4}
}

func (that *Protobuf4) Encode(value interface{}) (*bytes.Buffer, error) {
	if nil == value {
		return nil, nil
	}
	return Encode(value, func(input any) (*bytes.Buffer, error) {
		if message, ok := value.(proto.Message); ok {
			if buf, err := proto.Marshal(message); nil != err {
				return nil, cause.Error(err)
			} else {
				return bytes.NewBuffer(buf), nil
			}
		}
		if parameters, ok := value.(macro.Parameters); ok {
			return that.EncodeParameters(macro.Context(), parameters)
		}
		if returns, ok := value.(macro.Returns); ok {
			return that.EncodeReturns(macro.Context(), returns)
		}
		return nil, cause.Errorable(cause.CryptCodecError)
	})
}

func (that *Protobuf4) EncodeParameters(ctx context.Context, parameters macro.Parameters) (*bytes.Buffer, error) {
	arguments := make(map[int32][]byte)
	for idx, argument := range parameters.GetArguments(ctx) {
		if msg, ok := argument.(proto.Message); ok {
			if buf, err := proto.Marshal(msg); nil != err {
				return nil, err
			} else {
				arguments[int32(idx)] = buf
			}
		}
	}
	message := &proto4.InBound{
		Attachments: parameters.GetAttachments(ctx),
		Arguments:   arguments,
	}
	return that.doEncode(message)
}

func (that *Protobuf4) EncodeReturns(ctx context.Context, returns macro.Returns) (*bytes.Buffer, error) {
	message := &proto4.OutBound{
		Code:    returns.GetCode(),
		Message: returns.GetMessage(),
	}
	content := returns.GetContent(ctx)[0]
	if msg, ok := content.(proto.Message); ok {
		if buf, err := proto.Marshal(msg); nil != err {
			return nil, err
		} else {
			message.Content = buf
		}
	}
	ca := returns.GetCause(ctx)
	if ca != nil {
		message.Cause = &proto4.Cause{
			Name: ca.Name,
			Pos:  ca.Pos,
			Text: ca.Text,
			Buff: ca.Buff,
		}
	}
	return that.doEncode(message)
}

func (that *Protobuf4) Decode(value *bytes.Buffer, kind interface{}) (interface{}, error) {
	return Decode(value, kind, func(input *bytes.Buffer, kind any) (any, error) {
		if m, ok := kind.(proto.Message); ok {
			err := proto.Unmarshal(value.Bytes(), m)
			return m, cause.Error(err)
		}
		if parameters, ok := kind.(macro.Parameters); ok {
			return that.DecodeParameters(macro.Context(), value, parameters)
		}
		if returns, ok := kind.(macro.Returns); ok {
			return that.DecodeReturns(macro.Context(), value, returns)
		}
		return nil, cause.Errorable(cause.CryptCodecError)
	})
}

func (that *Protobuf4) DecodeParameters(ctx context.Context, value *bytes.Buffer, parameters macro.Parameters) (interface{}, error) {
	var message proto4.InBound
	err := proto.Unmarshal(value.Bytes(), &message)
	if nil != err {
		return nil, cause.Error(err)
	}
	parameters.SetAttachments(ctx, message.Attachments)
	var arguments []interface{}
	for idx, argument := range parameters.GetArguments(ctx) {
		if msg, ok := reflect.New(reflect.TypeOf(argument).Elem()).Interface().(proto.Message); ok {
			err = proto.Unmarshal(message.GetArguments()[int32(idx)], msg)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, msg)
		}
	}
	parameters.SetArguments(ctx, arguments...)
	return parameters, nil
}

func (that *Protobuf4) DecodeReturns(ctx context.Context, value *bytes.Buffer, returns macro.Returns) (interface{}, error) {
	var message proto4.OutBound
	err := proto.Unmarshal(value.Bytes(), &message)
	if err != nil {
		return nil, err
	}
	if len(returns.GetContent(ctx)) > 0 {
		var contents []interface{}
		for _, content := range returns.GetContent(ctx) {
			if content == nil {
				continue
			}
			if msg, ok := reflect.New(reflect.TypeOf(content).Elem()).Interface().(proto.Message); ok {
				err = proto.Unmarshal(message.GetContent(), msg)
				if err != nil {
					return nil, err
				}
				content = msg
				contents = append(contents, msg)
			}
		}
		returns.SetContent(ctx, contents...)
	}
	returns.SetCode(message.GetCode())
	returns.SetMessage(message.GetMessage())
	if message.GetCause() != nil {
		returns.SetCause(ctx, &macro.Cause{
			Name: message.GetCause().Name,
			Pos:  message.GetCause().Pos,
			Text: message.GetCause().Text,
			Buff: message.GetCause().Buff,
		})
	}
	return returns, nil
}

func (that *Protobuf4) EncodeString(value interface{}) (string, error) {
	return EncodeString(value, that)
}

func (that *Protobuf4) DecodeString(value string, kind interface{}) (interface{}, error) {
	return DecodeString(value, kind, that)
}

func (that *Protobuf4) doEncode(message proto.Message) (*bytes.Buffer, error) {
	buf, err := proto.Marshal(message)
	if err != nil {
		return nil, cause.Error(err)
	}
	return bytes.NewBuffer(buf), nil
}
