/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/types"
)

var IGraph = (*Graph)(nil)

// Graph spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Graph interface {

	// GraphQL apply graph query language.
	// @MPI("mesh.graph.expr")
	GraphQL(ctx context.Context, script *types.MeshQL) ([]map[string]any, error)

	// Sort apply graph query language.
	// @MPI("mesh.graph.sort")
	Sort(ctx context.Context, index *types.Paging) (*types.Page, error)

	// Link point with edge.
	// @MPI("mesh.graph.link")
	Link(ctx context.Context, quads []*types.Quad) error

	// Unlink point with edge.
	// @MPI("mesh.graph.unlink")
	Unlink(ctx context.Context, quads []*types.Quad) error

	// Sides edge.
	// @MPI("mesh.graph.sides")
	Sides(ctx context.Context, cursor *types.Cursor) ([]*types.Side, error)

	// Attrib point.
	// @MPI("mesh.graph.attrib")
	Attrib(ctx context.Context, vertex *types.Vertex) error

	// Dijkstra
	// @MPI("mesh.graph.dijkstra")
	Dijkstra(ctx context.Context, vector *types.Vec2) ([]*types.Quad, error)

	// Drop point.
	// @MPI("mesh.graph.drop")
	Drop(ctx context.Context, vector *types.Vec1) error

	// Vertex point.
	// @MPI("mesh.graph.vertex")
	Vertex(ctx context.Context, name string) ([]string, error)

	// Dump
	// @MPI("mesh.graph.dump")
	Dump(ctx context.Context, name string) ([]*types.Quad, error)
}
