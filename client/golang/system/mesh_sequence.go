/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.Sequence = new(systemSequence)
	macro.Provide(prsim.ISequence, new(systemSequence))
}

type systemSequence struct {
}

func (that *systemSequence) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *systemSequence) Next(ctx context.Context, kind string, length int) (string, error) {
	return aware.Sequence.Next(ctx, kind, length)
}

func (that *systemSequence) Section(ctx context.Context, kind string, size int, length int) ([]string, error) {
	return aware.Sequence.Section(ctx, kind, size, length)
}
