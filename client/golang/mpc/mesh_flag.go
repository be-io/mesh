/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

func init() {
	MeshFlag.OfProto = func(code string) codename {
		for _, value := range protoFlags {
			if value.Code() == code {
				return value
			}
		}
		return protoFlags[1]
	}
	MeshFlag.OfCodec = func(code string) codename {
		for _, value := range codecFlags {
			if value.Code() == code {
				return value
			}
		}
		return codecFlags[0]
	}
	MeshFlag.OfName = func(name string) codename {
		for _, value := range meshFlags {
			if value.Name() == name {
				return value
			}
		}
		return meshFlags[0]
	}
}

type codename interface {
	Code() string
	Name() string
}
type meshFlag struct {
	code string
	name string
}

var (
	meshFlags = []codename{
		MeshFlag.HTTP, MeshFlag.GRPC, MeshFlag.MQTT, MeshFlag.TCP,
		MeshFlag.JSON, MeshFlag.PROTOBUF, MeshFlag.XML, MeshFlag.YAML, MeshFlag.THRIFT,
	}
	protoFlags = []codename{
		MeshFlag.HTTP, MeshFlag.GRPC, MeshFlag.MQTT, MeshFlag.TCP,
	}
	codecFlags = []codename{
		MeshFlag.JSON, MeshFlag.PROTOBUF, MeshFlag.XML, MeshFlag.YAML, MeshFlag.THRIFT,
	}
)

var MeshFlag = struct {
	HTTP     codename
	GRPC     codename
	MQTT     codename
	TCP      codename
	JSON     codename
	PROTOBUF codename
	XML      codename
	THRIFT   codename
	YAML     codename
	OfProto  func(code string) codename
	OfCodec  func(code string) codename
	OfName   func(name string) codename
}{
	HTTP:     &meshFlag{code: "00", name: "http"},
	GRPC:     &meshFlag{code: "01", name: "grpc"},
	MQTT:     &meshFlag{code: "02", name: "mqtt"},
	TCP:      &meshFlag{code: "03", name: "tcp"},
	JSON:     &meshFlag{code: "00", name: "json"},
	PROTOBUF: &meshFlag{code: "01", name: "protobuf"},
	XML:      &meshFlag{code: "02", name: "xml"},
	THRIFT:   &meshFlag{code: "03", name: "thrift"},
	YAML:     &meshFlag{code: "04", name: "yaml"},
}

func (that *meshFlag) Code() string {
	return that.code
}

func (that *meshFlag) Name() string {
	return that.name
}
