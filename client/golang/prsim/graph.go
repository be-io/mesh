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
	// @MPI("mesh.graph.graphql")
	GraphQL(ctx context.Context, mql string, args map[string]any) ([]map[string]any, error)

	// Sort apply graph query language.
	// @MPI("mesh.graph.sort")
	Sort(ctx context.Context, index *types.Paging) (*types.Page, error)

	// Link point with edge.
	// @MPI("mesh.graph.link")
	Link(ctx context.Context, quads []*types.Quad) error

	// Unlink point with edge.
	// @MPI("mesh.graph.unlink")
	Unlink(ctx context.Context, quads []*types.Quad) error

	// Paths paths.
	// @MPI("mesh.graph.paths")
	Paths(ctx context.Context, cursor *types.Cursor) ([][]*types.Quad, error)

	// Dijkstra
	// @MPI("mesh.graph.dijkstra")
	Dijkstra(ctx context.Context, triple *types.Triple) ([]*types.Quad, error)

	// Drop point.
	// @MPI("mesh.graph.drop")
	Drop(ctx context.Context, quad *types.Quad) error

	// Vertex points.
	// @MPI("mesh.graph.vertexes")
	Vertex(ctx context.Context) ([]*types.Quad, error)

	// Dump
	// @MPI("mesh.graph.dump")
	Dump(ctx context.Context) ([]*types.Quad, error)
}
