/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"bytes"
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/dsa"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"reflect"
)

var _ macro.ServiceInspector = new(MeshEndpoint)
var _ prsim.EndpointSticker[*bytes.Buffer, *bytes.Buffer] = new(MeshEndpoint)
var _ prsim.Endpoint = new(MeshEndpoint)
var _ prsim.EndpointSticker[reflect.Value, reflect.Value] = new(reflectEndpointSticker)
var zero = reflect.Value{}

// MeshEndpoint
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
// @ServiceInspector
type MeshEndpoint struct {
	stickers dsa.Map[string, prsim.EndpointSticker[reflect.Value, reflect.Value]]
}

func (that *MeshEndpoint) Fuzzy(ctx context.Context, buff []byte) ([]byte, error) {
	r, err := that.Stick(ctx, bytes.NewBuffer(buff))
	if nil != err {
		return nil, cause.Error(err)
	}
	if nil == r {
		return nil, nil
	}
	return r.Bytes(), nil
}

func (that *MeshEndpoint) Inspect() []macro.MPI {
	var services []macro.MPI
	for _, s := range macro.Load(prsim.IEndpointSticker).List() {
		if sticker, ok := s.(macro.MPI); ok {
			services = append(services, sticker)
		}
	}
	return services
}

func (that *MeshEndpoint) Rtt() *macro.Rtt {
	return &macro.Rtt{}
}

func (that *MeshEndpoint) I() *bytes.Buffer {
	return new(bytes.Buffer)
}

func (that *MeshEndpoint) O() *bytes.Buffer {
	return new(bytes.Buffer)
}

func (that *MeshEndpoint) Stickers() dsa.Map[string, prsim.EndpointSticker[reflect.Value, reflect.Value]] {
	if nil == that.stickers {
		that.stickers = dsa.NewStringMap[prsim.EndpointSticker[reflect.Value, reflect.Value]]()
	}
	return that.stickers
}

func (that *MeshEndpoint) Stick(ctx context.Context, varg *bytes.Buffer) (*bytes.Buffer, error) {
	mtx := mpc.ContextWith(ctx)
	urn := types.FromURN(ctx, mtx.GetUrn())
	spi := macro.Load(prsim.IEndpointSticker).Get(urn.Name)
	if nil == spi {
		return nil, mpc.NoService(ctx, urn)
	}
	sticker := that.Stickers().PutIfy(urn.Name, func(k string) prsim.EndpointSticker[reflect.Value, reflect.Value] {
		s := reflect.ValueOf(spi)
		return &reflectEndpointSticker{
			sticker: s,
			stt:     s.MethodByName("Rtt"),
			i:       s.MethodByName("I"),
			o:       s.MethodByName("O"),
			stick:   s.MethodByName("Stick"),
		}
	})
	cn := mpc.MeshFlag.OfCodec(urn.Flag.Codec).Name()
	cs := macro.Load(codec.ICodec).Get(cn)
	if nil == cs {
		return nil, cause.NoImplement(cn)
	}
	cdc, ok := cs.(codec.Codec)
	if !ok {
		return nil, cause.NoImplement(cn)
	}
	input, err := func() (reflect.Value, error) {
		input := sticker.I()
		if zero == input {
			return zero, cause.CompatibleError("Unexpected definition of %s", urn.Name)
		}
		if input.Kind() == reflect.Pointer {
			_, err := cdc.Decode(varg, input.Interface())
			return input, cause.Error(err)
		}
		pointer := reflect.New(input.Type())
		pointer.Elem().Set(input)
		_, err := cdc.Decode(varg, pointer.Interface())
		return pointer.Elem(), cause.Error(err)
	}()
	output, err := sticker.Stick(ctx, input)
	if nil != err {
		return nil, cause.Error(err)
	}
	return cdc.Encode(output.Interface())
}

type reflectEndpointSticker struct {
	sticker reflect.Value
	stt     reflect.Value
	i       reflect.Value
	o       reflect.Value
	stick   reflect.Value
}

func (that *reflectEndpointSticker) Rtt() *macro.Rtt {
	vs := that.stt.Call(nil)
	if len(vs) < 1 {
		return nil
	}
	if v, ok := vs[0].Interface().(*macro.Rtt); ok {
		return v
	}
	return nil
}

func (that *reflectEndpointSticker) I() reflect.Value {
	vs := that.i.Call(nil)
	if len(vs) < 1 {
		return zero
	}
	return vs[0]
}

func (that *reflectEndpointSticker) O() reflect.Value {
	vs := that.o.Call(nil)
	if len(vs) < 1 {
		return zero
	}
	return vs[0]
}

func (that *reflectEndpointSticker) Stick(ctx context.Context, varg reflect.Value) (reflect.Value, error) {
	urn := types.FromURN(ctx, mpc.ContextWith(ctx).GetUrn())
	vs := that.stick.Call([]reflect.Value{reflect.ValueOf(ctx), varg})
	if len(vs) < 2 {
		return zero, cause.CompatibleError("Unexpected definition of %s", urn.Name)
	}
	if err, ok := vs[1].Interface().(error); ok {
		return zero, err
	}
	return vs[0], nil
}
