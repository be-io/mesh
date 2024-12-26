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

var IKMS = (*KMS)(nil)

// KMS spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type KMS interface {

	// Reset will override the keystore environ.
	// @MPI("kms.store.reset")
	Reset(ctx context.Context, env *types.Environ) error

	// Environ will return the keystore environ.
	// @MPI("kms.store.environ")
	Environ(ctx context.Context) (*types.Environ, error)

	// List will return the keystore environ.
	// @MPI("kms.crt.store.list")
	List(ctx context.Context, cno string) ([]*types.Keys, error)

	// ApplyRoot will apply the root certification.
	// @MPI("kms.crt.apply.root")
	ApplyRoot(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error)

	// ApplyIssue will apply the common certification.
	// @MPI("kms.crt.apply.issue")
	ApplyIssue(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error)
}
