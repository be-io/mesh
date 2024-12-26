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
)

func init() {
	var _ prsim.Cluster = new(MeshCluster)
	macro.Provide(prsim.ICluster, new(MeshCluster))
}

type MeshCluster struct {
}

func (that *MeshCluster) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}
func (that *MeshCluster) Election(ctx context.Context, buff []byte) ([]byte, error) {
	return aware.Cluster.Election(ctx, buff)
}

func (that *MeshCluster) IsLeader(ctx context.Context) (bool, error) {
	return aware.Cluster.IsLeader(ctx)
}
