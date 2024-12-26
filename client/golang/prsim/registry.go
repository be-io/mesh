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

var IRegistry = (*Registry)(nil)

var (
	MetadataRegisterEvent = &macro.Btt{Topic: "mesh.registry.event.metadata", Code: "*"}
	ProxyRegisterEvent    = &macro.Btt{Topic: "mesh.registry.event.proxy", Code: "*"}
	RegistryEventRefresh  = &macro.Btt{Topic: "mesh.registry.event.refresh", Code: "*"}
)

// Registry spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Registry interface {

	// Register
	// @MPI("mesh.registry.put")
	Register(ctx context.Context, registration *types.Registration[any]) error

	// Registers
	// @MPI("mesh.registry.puts")
	Registers(ctx context.Context, registrations []*types.Registration[any]) error

	// Unregister
	// @MPI("mesh.registry.remove")
	Unregister(ctx context.Context, registration *types.Registration[any]) error

	// Export
	// @MPI("mesh.registry.export")
	Export(ctx context.Context, kind string) ([]*types.Registration[any], error)
}
