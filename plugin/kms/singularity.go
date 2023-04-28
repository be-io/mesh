/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package kms

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.Singularity = new(kmsSingularity)
	macro.Provide(prsim.ISingularity, new(kmsSingularity))
}

type kmsSingularity struct {
}

func (that *kmsSingularity) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSPI}
}

func (that *kmsSingularity) Reveal(ctx context.Context) (*prsim.RootKey, error) {
	return nil, cause.UrnNotPermit.Error()
}
