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

var ILicenser = (*Licenser)(nil)

var LicenseImports = &macro.Btt{Topic: "mesh.license.imports", Code: "*"}

// Licenser
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Licenser interface {

	// Imports the licenses.
	// @MPI(name = "mesh.license.imports", flags = 2)
	Imports(ctx context.Context, license string) error

	// Exports the licenses.
	// @MPI(name = "mesh.license.exports", flags = 2)
	Exports(ctx context.Context) (string, error)

	// Explain the license.
	// @MPI(name = "mesh.license.explain", flags = 2)
	Explain(ctx context.Context) (*types.License, error)

	// Verify the license.
	// @MPI(name = "mesh.license.verify", flags = 2)
	Verify(ctx context.Context) (int64, error)

	// Features is license features.
	// @MPI(name = "mesh.license.features", flags = 2)
	Features(ctx context.Context) (map[string]string, error)
}
