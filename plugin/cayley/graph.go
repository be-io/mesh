/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cayley

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
)

func init() {
	var _ prsim.Graph = new(cayleyGraph)
	macro.Provide(prsim.IGraph, new(cayleyGraph))
}

type cayleyGraph struct {
}

func (that *cayleyGraph) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *cayleyGraph) GraphQL(ctx context.Context, mql string, args map[string]any) ([]map[string]any, error) {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Sort(ctx context.Context, index *types.Paging) (*types.Page, error) {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Link(ctx context.Context, quads []*types.Quad) error {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Unlink(ctx context.Context, quads []*types.Quad) error {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Paths(ctx context.Context, cursor *types.Cursor) ([][]*types.Quad, error) {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Dijkstra(ctx context.Context, triple *types.Triple) ([]*types.Quad, error) {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Drop(ctx context.Context, quad *types.Quad) error {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Vertex(ctx context.Context) ([]*types.Quad, error) {
	//TODO implement me
	panic("implement me")
}

func (that *cayleyGraph) Dump(ctx context.Context) ([]*types.Quad, error) {
	//TODO implement me
	panic("implement me")
}
