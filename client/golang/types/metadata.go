/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

const (
	Binding  = "Binding"
	MPS      = "MPS"
	Forward  = "Forward"
	Restful  = "Restful"
	Protobuf = "Protobuf"
)

type Metadata struct {
	References []*Reference `index:"0" json:"references" xml:"references" yaml:"references"`

	Services []*Service `index:"5" json:"services" xml:"services" yaml:"services"`
}
