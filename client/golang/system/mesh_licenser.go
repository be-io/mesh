/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ prsim.Licenser = new(systemLicenser)
	macro.Provide(prsim.ILicenser, new(systemLicenser))
}

type systemLicenser struct {
}

func (that *systemLicenser) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *systemLicenser) Imports(ctx context.Context, license string) error {
	return aware.Licenser.Imports(ctx, license)
}

func (that *systemLicenser) Exports(ctx context.Context) (string, error) {
	return aware.Licenser.Exports(ctx)
}

func (that *systemLicenser) Explain(ctx context.Context) (*types.License, error) {
	return aware.Licenser.Explain(ctx)
}

func (that *systemLicenser) Verify(ctx context.Context) (int64, error) {
	return aware.Licenser.Verify(ctx)
}

func (that *systemLicenser) Features(ctx context.Context) (map[string]string, error) {
	return aware.Licenser.Features(ctx)
}
