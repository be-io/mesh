/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	_ "github.com/be-io/mesh/plugin/metabase"
	"time"
)

var _ prsim.Sequence = new(PRSISequence)

// PRSISequence
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSISequence struct {
}

func (that *PRSISequence) Next(ctx context.Context, kind string, length int) (string, error) {
	environ, err := aware.LocalNet.GetEnviron(ctx)
	if nil != err {
		return "", cause.Error(err)
	}
	node, err := types.FromNodeID(environ.NodeId)
	if nil != err {
		return "", cause.Error(err)
	}
	switch kind {
	case types.MeshNode:
		return aware.Sequence.Next(ctx, kind, types.MeshIDLength)
	case types.MeshInstitution:
		id, err := aware.Sequence.Next(ctx, kind, types.MeshIDLength)
		if nil != err {
			return "", cause.Error(err)
		}
		return types.SaaSInstID(node.SEQ, id).String(), nil
	default:
		seq, err := aware.Sequence.Next(ctx, kind, length)
		if nil != err {
			return "", cause.Error(err)
		}
		return fmt.Sprintf("%s%s%s", time.Now().Format("20060102"), node.SEQ, seq), nil
	}
}

func (that *PRSISequence) Section(ctx context.Context, kind string, size int, length int) ([]string, error) {
	return aware.Sequence.Section(ctx, kind, size, length)
}
