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

var ICommercialize = (*Commercialize)(nil)

// Commercialize
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Commercialize interface {

	// Sign license, only glab has permission to sign.
	// @MPI("mesh.license.sign")
	Sign(ctx context.Context, lsr *types.License) (*types.CommerceLicense, error)

	// History list the sign license in history, the latest is the first index.
	// @MPI("mesh.license.history")
	History(ctx context.Context, instId string) ([]*types.CommerceLicense, error)

	// Issued mesh node identity.
	// @MPI("mesh.net.issued")
	Issued(ctx context.Context, name string, kind string, cname string) (*types.CommerceEnviron, error)

	// Dump the node identity.
	// @MPI("mesh.net.dump")
	Dump(ctx context.Context, nodeId string) ([]*types.CommerceEnviron, error)
}
