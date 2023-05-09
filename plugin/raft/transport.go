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
	"github.com/lni/dragonboat/v4/config"
	"github.com/lni/dragonboat/v4/raftio"
	"github.com/lni/dragonboat/v4/raftpb"
)

type mt struct {
}

func (that *mt) Create(config config.NodeHostConfig, mh raftio.MessageHandler, ch raftio.ChunkHandler) raftio.ITransport {
	return &transport{
		config: config,
		mh:     mh,
		ch:     ch,
	}
}

func (that *mt) Validate(addr string) bool {
	return true
}

type transport struct {
	config config.NodeHostConfig
	mh     raftio.MessageHandler
	ch     raftio.ChunkHandler
}

func (that *transport) Name() string {
	return that.Name()
}

func (that *transport) Start() error {
	return nil
}

func (that *transport) Close() error {
	return nil
}

func (that *transport) GetConnection(ctx context.Context, target string) (raftio.IConnection, error) {
	return &connection{ctx: ctx, target: target}, nil
}

func (that *transport) GetSnapshotConnection(ctx context.Context, target string) (raftio.ISnapshotConnection, error) {
	return &snapshotConnection{ctx: ctx, target: target}, nil
}

type connection struct {
	ctx    context.Context
	target string
}

func (that *connection) Close() {

}

func (that *connection) SendMessageBatch(batch raftpb.MessageBatch) error {
	buff, err := batch.Marshal()
	if nil != err {
		return cause.Error(err)
	}
	_, err = aware.Cluster.SendMessage(that.ctx, buff)
	return cause.Error(err)
}

type snapshotConnection struct {
	ctx    context.Context
	target string
}

func (that *snapshotConnection) Close() {

}

func (that *snapshotConnection) SendChunk(chunk raftpb.Chunk) error {
	buff, err := chunk.Marshal()
	if nil != err {
		return cause.Error(err)
	}
	_, err = aware.Cluster.SendChunk(that.ctx, buff)
	return cause.Error(err)
}
