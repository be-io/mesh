/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
)

var _ prsim.Graph = new(PRSIGraph)

// PRSIGraph
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSIGraph struct {
}

func (that *PRSIGraph) GraphQL(ctx context.Context, mql string, args map[string]any) ([]map[string]any, error) {
	return aware.Graph.GraphQL(ctx, mql, args)
}

func (that *PRSIGraph) Sort(ctx context.Context, index *types.Paging) (*types.Page, error) {
	return aware.Graph.Sort(ctx, index)
}

func (that *PRSIGraph) Link(ctx context.Context, quads []*types.Quad) error {
	return aware.Graph.Link(ctx, quads)
}

func (that *PRSIGraph) Unlink(ctx context.Context, quads []*types.Quad) error {
	return aware.Graph.Unlink(ctx, quads)
}

func (that *PRSIGraph) Paths(ctx context.Context, cursor *types.Cursor) ([][]*types.Quad, error) {
	return aware.Graph.Paths(ctx, cursor)
}

func (that *PRSIGraph) Dijkstra(ctx context.Context, triple *types.Triple) ([]*types.Quad, error) {
	return aware.Graph.Dijkstra(ctx, triple)
}

func (that *PRSIGraph) Drop(ctx context.Context, quad *types.Quad) error {
	return aware.Graph.Drop(ctx, quad)
}

func (that *PRSIGraph) Vertex(ctx context.Context) ([]*types.Quad, error) {
	return aware.Graph.Vertex(ctx)
}

func (that *PRSIGraph) Dump(ctx context.Context) ([]*types.Quad, error) {
	return aware.Graph.Dump(ctx)
}
