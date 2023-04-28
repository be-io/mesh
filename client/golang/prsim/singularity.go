/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
)

func init() {
	var _ Singularity = new(singularity)
	macro.Provide(ISingularity, new(singularity))
}

var ISingularity = (*Singularity)(nil)

// Singularity is a root extension endpoint.
type Singularity interface {

	// Reveal is start up point.
	Reveal(ctx context.Context) (*RootKey, error)
}

type RootKey struct {
	PublicKey  string `json:"public_key" xml:"public_key" yaml:"public_key"`
	PrivateKey string `json:"private_key" xml:"private_key" yaml:"private_key"`
}

type singularity struct {
}

func (that *singularity) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshNop}
}

func (that *singularity) Reveal(ctx context.Context) (*RootKey, error) {
	return nil, cause.NotImplementError()
}
