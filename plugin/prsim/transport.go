/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

var _ prsim.Transport = new(PRSITransport)

// PRSITransport
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSITransport struct {
}

func (that *PRSITransport) Open(ctx context.Context, sessionId string, metadata map[string]string) (prsim.Session, error) {
	return nil, nil
}

func (that *PRSITransport) Close(ctx context.Context, timeout types.Duration) error {
	return nil
}

func (that *PRSITransport) Roundtrip(ctx context.Context, payload []byte, metadata map[string]string) ([]byte, error) {
	return nil, nil
}
