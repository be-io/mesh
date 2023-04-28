/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"bytes"
	"context"
	"io"
)

var IConsumer = (*Consumer)(nil)

type Consumer interface {

	// Start the mesh broker.
	Start() error

	// Consume the input payload.
	// urn       Actual uniform resource Domain Name.
	// execution Service reference.
	// inbound   Input arguments.
	// Output    payload
	Consume(ctx context.Context, urn string, execution Execution, inbound *bytes.Buffer) (*bytes.Buffer, error)

	// Closer is release hook
	io.Closer
}
