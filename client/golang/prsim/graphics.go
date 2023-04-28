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

var IGraphics = (*Graphics)(nil)

// Graphics spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Graphics interface {

	// Captcha apply a graphics captcha.
	// @MPI("mesh.graphics.captcha.apply")
	Captcha(ctx context.Context, kind string, features map[string]string) (*types.Captcha, error)

	// Verify a graphics captcha value.
	// @MPI("mesh.graphics.captcha.verify")
	Verify(ctx context.Context, mno string, value string) (bool, error)
}
