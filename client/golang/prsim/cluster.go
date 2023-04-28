/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
)

var ICluster = (*Cluster)(nil)

// Cluster spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Cluster interface {
	// Election will election leader of instances.
	// @MPI("mesh.cluster.election")
	Election(ctx context.Context, buff []byte) ([]byte, error)
	// IsLeader if same level.
	// @MPI("mesh.cluster.leader")
	IsLeader(ctx context.Context) (bool, error)
}
