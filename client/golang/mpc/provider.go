/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"crypto/tls"
	"io"
)

var IProvider = (*Provider)(nil)

type Provider interface {

	// Start the mesh broker.
	Start(ctx context.Context, address string, tc *tls.Config) error

	io.Closer
}
