/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package raft

import (
	"context"
)

//go:generate go run ../../client/golang/dyn/mpc.go -d ./ -m github.com/be-io/mesh/plugin/raft -p github.com/be-io/mesh/plugin/raft

var iCluster = (*channel)(nil)

// channel spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type channel interface {

	// SendMessage send message with rpc.
	// @MPI("mesh.raft.message")
	SendMessage(ctx context.Context, buff []byte) ([]byte, error)
	// SendChunk send chunk with rpc.
	// @MPI("mesh.raft.chunk")
	SendChunk(ctx context.Context, buff []byte) ([]byte, error)
}

var _ channel = new(raftChannel)

// raftChannel
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type raftChannel struct {
}

func (that *raftChannel) SendMessage(ctx context.Context, buff []byte) ([]byte, error) {
	return nil, nil
}

func (that *raftChannel) SendChunk(ctx context.Context, buff []byte) ([]byte, error) {
	return nil, nil
}
