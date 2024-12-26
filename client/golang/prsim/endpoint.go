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
)

var IEndpoint = (*Endpoint)(nil)

// Endpoint
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Endpoint interface {

	// Fuzzy call with generic param
	// In multi returns, it's an array.
	// @MPI("${mesh.uname}")
	Fuzzy(ctx context.Context, buff []byte) ([]byte, error)
}

var IEndpointSticker = (*EndpointSticker[any, any])(nil)

type EndpointSticker[I any, O any] interface {

	// MPI attributes
	macro.MPI

	// I is the input
	I() I

	// O is the output
	O() O

	// Stick with generic param
	Stick(ctx context.Context, varg I) (O, error)
}
