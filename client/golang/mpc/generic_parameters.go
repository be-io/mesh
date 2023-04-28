/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
)

func init() {
	var _ macro.Parameters = new(GenericParameters)
}

type GenericParameters map[string]interface{}

func (that GenericParameters) PutAll(dict map[string]interface{}) {
	for key, value := range dict {
		that[key] = value
	}
}

func (that GenericParameters) GetKind() interface{} {
	return (*map[string]interface{})(nil)
}

func (that GenericParameters) GetArguments(ctx context.Context) []interface{} {
	var arguments []interface{}
	for key, value := range that {
		if key == "attachments" {
			continue
		}
		arguments = append(arguments, value)
	}
	return arguments
}

func (that GenericParameters) SetArguments(ctx context.Context, arguments ...interface{}) {

}

func (that GenericParameters) GetAttachments(ctx context.Context) map[string]string {
	attachments := that["attachments"]
	if nil == attachments {
		return nil
	}
	if attach, ok := attachments.(map[string]string); ok {
		return attach
	}
	cdc, ok := macro.Load(codec.ICodec).Get("json").(codec.Codec)
	if !ok {
		log.Error(ctx, "Parse attachments without codec named json exist")
		return nil
	}
	buff, err := cdc.Encode(attachments)
	if nil != err {
		log.Error(ctx, err.Error())
		return nil
	}
	var attach map[string]string
	if _, err = cdc.Decode(buff, &attach); nil != err {
		log.Error(ctx, err.Error())
		return nil
	}
	return attach
}

func (that GenericParameters) SetAttachments(ctx context.Context, attachments map[string]string) {
	that["attachments"] = attachments
}

func (that GenericParameters) Get(key string) interface{} {
	return that[key]
}
