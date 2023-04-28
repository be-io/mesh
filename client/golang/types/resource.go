/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

const (
	TCP   = "TCP"
	UDP   = "UDP"
	HTTP1 = "HTTP1"
	HTTP2 = "HTTP2"
	HTTP3 = "HTTP3"

	JSON = "JSON"
)

type Resource struct {
	ID           string `json:"id" yaml:"id"`
	Module       string `json:"module" yaml:"module"`
	Kind         string `json:"kind" yaml:"kind"`
	Name         string `json:"name" yaml:"name"`
	Version      string `json:"version" yaml:"version"`
	Lang         string `json:"lang" yaml:"lang"`
	NetworkProto string `json:"network_proto" yaml:"network_proto"`
	PayloadProto string `json:"payload_proto" yaml:"payload_proto"`
	IP           string `json:"ip" yaml:"ip"`
	Port         string `json:"port" yaml:"port"`
}
