/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package raft

import (
	"bytes"
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/types"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"math"
	"time"
)

type EtcdRaftNode struct {
	id          uint64
	ctx         *types.ClusterContext
	startTime   time.Time
	store       raft.Storage
	pstore      map[string]string
	cfg         *raft.Config
	raft        raft.Node
	ticker      <-chan time.Time
	done        <-chan struct{}
	peers       map[uint64]string
	confChanges map[uint64]chan error
}

func StartNode(rc *types.ClusterContext, store raft.Storage) (*EtcdRaftNode, error) {
	heartbeat := 1
	snap, err := store.Snapshot()
	if nil != err {
		return nil, cause.Error(err)
	}
	n := &EtcdRaftNode{
		id:        rc.Id,
		ctx:       rc,
		startTime: time.Now(),
		store:     store,
		cfg: &raft.Config{
			ID:                       rc.Id,
			ElectionTick:             20 * heartbeat, // 2s if we call Tick() every 100 ms.
			HeartbeatTick:            heartbeat,      // 100ms if we call Tick() every 100 ms.
			Storage:                  store,
			MaxSizePerMsg:            math.MaxUint16, // 256 KB should allow more batching.
			MaxInflightMsgs:          256,
			MaxCommittedSizePerReady: 64 << 20, // Avoid loading entire Raft log into memory.
			// We don't need lease based reads. They cause issues because they
			// require CheckQuorum to be true, and that causes a lot of issues
			// for us during cluster bootstrapping and later. A seemingly
			// healthy cluster would just cause leader to step down due to
			// "inactive" quorum, and then disallow anyone from becoming leader.
			// So, let's stick to default options.  Let's achieve correctness,
			// then we achieve performance. Plus, for the Dgraph alphas, we'll
			// be soon relying only on Timestamps for blocking reads and
			// achieving linearizability, than checking quorums (Zero would
			// still check quorums).
			ReadOnlyOption: raft.ReadOnlySafe,
			// When a disconnected EtcdRaftNode joins back, it forces a leader change,
			// as it starts with a higher term, as described in Raft thesis (not
			// the paper) in section 9.6. This setting can avoid that by only
			// increasing the term, if the EtcdRaftNode has a good chance of becoming
			// the leader.
			PreVote: true,
			// We can explicitly set Applied to the first index in the Raft log,
			// so it does not derive it separately, thus avoiding a crash when
			// the Applied is set to below snapshot index by Raft.
			// In case this is a new Raft log, first would be 1, and therefore
			// Applied would be zero, hence meeting the condition by the library
			// that Applied should only be set during a restart.
			//
			// Update: Set the Applied to the latest snapshot, because it seems
			// like somehow the first index can be out of sync with the latest
			// snapshot.
			Applied: snap.Metadata.Index,
		},
		pstore:      make(map[string]string),
		ticker:      time.Tick(time.Second),
		done:        make(chan struct{}),
		confChanges: make(map[uint64]chan error),
		peers:       make(map[uint64]string),
	}

	n.raft = raft.StartNode(n.cfg, nil)
	return n, nil
}

func RestartNode(snap raftpb.Snapshot, st raftpb.HardState, entries []raftpb.Entry) (raft.Node, error) {
	storage := raft.NewMemoryStorage()
	// Recover the in-memory storage from persistent snapshot, state and entries.
	if err := storage.ApplySnapshot(snap); nil != err {
		return nil, cause.Error(err)
	}
	if err := storage.SetHardState(st); nil != err {
		return nil, cause.Error(err)
	}
	if err := storage.Append(entries); nil != err {
		return nil, cause.Error(err)
	}
	config := &raft.Config{
		ID:              0x01,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}

	// Restart raft without peer information.
	// Peer information is already included in the storage.
	return raft.RestartNode(config), nil
}

func (that *EtcdRaftNode) run() {
	for {
		select {
		case <-that.ticker:
			that.raft.Tick()
		case ready := <-that.raft.Ready():
			ctx := mpc.Context()
			that.saveToStorage(&ready.HardState, ready.Entries, &ready.Snapshot)
			that.send(ctx, ready.Messages)
			if !raft.IsEmptySnap(ready.Snapshot) {
				that.processSnapshot(ready.Snapshot)
			}
			for _, entry := range ready.CommittedEntries {
				that.process(ctx, entry)
				if entry.Type == raftpb.EntryConfChange {
					var cc raftpb.ConfChange
					cc.Unmarshal(entry.Data)
					that.raft.ApplyConfChange(cc)
				}
			}
			that.raft.Advance()
		case <-that.done:
			return
		}
	}
}

func (that *EtcdRaftNode) saveToStorage(hardState *raftpb.HardState, entries []raftpb.Entry, snapshot *raftpb.Snapshot) {

}

func (that *EtcdRaftNode) send(ctx context.Context, messages []raftpb.Message) {
	for _, message := range messages {
		log.Info(ctx, raft.DescribeMessage(message, nil))
		// send message to other EtcdRaftNode
		that.receive(ctx, message)
	}
}

func (that *EtcdRaftNode) processSnapshot(snapshot raftpb.Snapshot) {
	panic(fmt.Sprintf("Applying snapshot on EtcdRaftNode %v is not implemented", that.id))
}

func (that *EtcdRaftNode) process(ctx context.Context, entry raftpb.Entry) {
	log.Info(ctx, "EtcdRaftNode %v: processing entry: %v\n", that.id, entry)
	if entry.Type == raftpb.EntryNormal && entry.Data != nil {
		parts := bytes.SplitN(entry.Data, []byte(":"), 2)
		that.pstore[string(parts[0])] = string(parts[1])
	}
}

func (that *EtcdRaftNode) receive(ctx context.Context, message raftpb.Message) {
	that.raft.Step(ctx, message)
}

func main() error {
	ctx := mpc.Context()
	rc := &types.ClusterContext{
		Id:        0,
		Addr:      "",
		Group:     0,
		IsLearner: false,
	}
	// start a small cluster
	node, err := StartNode(rc, nil)
	if nil != err {
		return cause.Error(err)
	}
	if err = node.raft.Campaign(ctx); nil != err {
		return cause.Error(err)
	}
	node.run()
	if err = node.raft.ProposeConfChange(ctx, raftpb.ConfChange{
		ID:      3,
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  3,
		Context: []byte(""),
	}); nil != err {
		return cause.Error(err)
	}
	// Wait for leader, is there a better way to do this
	for node.raft.Status().Lead != 1 {
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
