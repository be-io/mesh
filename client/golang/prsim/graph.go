/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/types"
)

var IGraph = (*Graph)(nil)

// Graph spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Graph interface {

	// GraphQL apply graph query language.
	// @MPI("mesh.graph.graphql")
	GraphQL(ctx context.Context, mql string, args map[string]any) ([]map[string]any, error)

	// Vertex
	// @MPI("mesh.graph.vertex")
	Vertex(ctx context.Context, id string) (*types.Vertex, error)

	// Vertices
	// @MPI("mesh.graph.vertices")
	Vertices(ctx context.Context, pattern string) ([]string, error)

	// Vertexes points.
	// @MPI("mesh.graph.vertexes")
	Vertexes(ctx context.Context, index *types.Paging) (*types.Page[*types.Vertex], error)

	// Side
	// @MPI("mesh.graph.side")
	Side(ctx context.Context, triple *types.Triple) (*types.Side, error)

	// Sides
	// @MPI("mesh.graph.sides")
	Sides(ctx context.Context, index *types.Paging) (*types.Page[*types.Side], error)

	// Link vertex with side.
	// @MPI("mesh.graph.link")
	Link(ctx context.Context, quads []*types.Quad) error

	// Unlink vertex with side.
	// @MPI("mesh.graph.unlink")
	Unlink(ctx context.Context, triples []*types.Triple) error

	// Paths paths.
	// @MPI("mesh.graph.paths")
	Paths(ctx context.Context, tuple *types.Tuple) ([][]*types.Quad, error)

	// Dijkstra
	// @MPI("mesh.graph.dijkstra")
	Dijkstra(ctx context.Context, triple *types.Triple) ([][]*types.Quad, error)

	// Drop point.
	// @MPI("mesh.graph.drop")
	Drop(ctx context.Context, triples []*types.Triple) error

	// Dump
	// @MPI("mesh.graph.dump")
	Dump(ctx context.Context) ([]*types.Quad, error)
}
