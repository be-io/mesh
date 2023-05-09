/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package raft

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.Cluster = new(meshCluster)
	macro.Provide(prsim.ICluster, new(meshCluster))
}

const Name = "raft"

type meshCluster struct {
}

func (that *meshCluster) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *meshCluster) Election(ctx context.Context, buff []byte) ([]byte, error) {
	return nil, cause.NotImplementError()
}

func (that *meshCluster) IsLeader(ctx context.Context) (bool, error) {
	return false, cause.NotImplementError()
}
