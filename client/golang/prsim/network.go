/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/types"
)

const (
	AutoDomain = "auto"
	ManuDomain = "manu"
)

var INetwork = (*Network)(nil)

var (
	NetworkRouteRefresh = &macro.Btt{Topic: "mesh.net.route.refresh", Code: "*"}
	RoutePeriodRefresh  = &macro.Btt{Topic: "mesh.plugin.proxy.dynamic.refresh", Code: "*"}
)

// Network spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Network interface {

	// GetEnviron Get the meth network environment fixed information.
	// @MPI("mesh.net.environ")
	GetEnviron(ctx context.Context) (*types.Environ, error)

	// Accessible Check the mesh network is accessible.
	// @MPI("mesh.net.accessible")
	Accessible(ctx context.Context, route *types.Route) (bool, error)

	// Refresh the routes to mesh network.
	// @MPI("mesh.net.refresh")
	Refresh(ctx context.Context, routes []*types.Route) error

	// GetRoute the network edge route.
	// @MPI("mesh.net.edge")
	GetRoute(ctx context.Context, nodeId string) (*types.Route, error)

	// GetRoutes the network edge routes.
	// @MPI("mesh.net.edges")
	GetRoutes(ctx context.Context) ([]*types.Route, error)

	// GetDomains the network domains.
	// @MPI("mesh.net.domains")
	GetDomains(ctx context.Context, kind string) ([]*types.Domain, error)

	// PutDomains the network domains.
	// @MPI("mesh.net.resolve")
	PutDomains(ctx context.Context, kind string, domains []*types.Domain) error

	// Weave the network.
	// @MPI("mesh.net.weave")
	Weave(ctx context.Context, route *types.Route) error

	// Ack the network.
	// @MPI("mesh.net.ack")
	Ack(ctx context.Context, route *types.Route) error

	// Disable the network
	// @MPI("mesh.net.disable")
	Disable(ctx context.Context, nodeId string) error

	// Enable the network
	// @MPI("mesh.net.enable")
	Enable(ctx context.Context, nodeId string) error

	// Index the network edges
	// @MPI("mesh.net.index")
	Index(ctx context.Context, index *types.Paging) (*types.Page[*types.Route], error)

	// Version
	// Network environment version.
	// @MPI("mesh.net.version")
	Version(ctx context.Context, nodeId string) (*types.Versions, error)

	// Instx
	// Network institutions.
	// @MPI("mesh.net.instx")
	Instx(ctx context.Context, index *types.Paging) (*types.Page[*types.Institution], error)

	// Instr
	// Network institutions.
	// @MPI("mesh.net.instr")
	Instr(ctx context.Context, institutions []*types.Institution) error

	// Ally
	// Network form alliance.
	// @MPI("mesh.net.ally")
	Ally(ctx context.Context, nodeIds []string) error

	// Disband
	// Network quit alliance.
	// @MPI("mesh.net.disband")
	Disband(ctx context.Context, nodeIds []string) error

	// Assert
	// Network feature assert.
	// @MPI("mesh.net.assert")
	Assert(ctx context.Context, feature string, nodeIds []string) (bool, error)
}
