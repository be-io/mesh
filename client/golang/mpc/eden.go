/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
)

var IEden = (*Eden)(nil)

// Eden
// @SPI("mesh")
type Eden interface {

	// Define the reference object.
	Define(ctx context.Context, mpi macro.MPI, reference interface{}) error

	// Refer the service reference by method.
	Refer(ctx context.Context, mpi macro.MPI, reference interface{}, method macro.Inspector) (Execution, error)

	// Store the service object.
	Store(ctx context.Context, kind interface{}, service interface{}) error

	// Infer the reference service by Domain.
	Infer(ctx context.Context, urn string) (Execution, error)

	// ReferTypes Get all reference types.
	ReferTypes(ctx context.Context) ([]macro.MPI, error)

	// InferTypes Get all service types.
	InferTypes(ctx context.Context) ([]macro.MPS, error)
}
