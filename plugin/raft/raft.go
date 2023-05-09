/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package raft

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/lni/dragonboat/v4/plugin/tan"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lni/dragonboat/v4/config"
	"github.com/lni/dragonboat/v4/logger"
	"github.com/lni/goutils/syncutil"
)

func init() {
	plugin.Provide(new(raftPlugin))
}

const (
	// we use two raft groups in this example, they are identified by the cluster
	// ID values below
	shardID1 uint64 = 100
	shardID2 uint64 = 101
)

type raftPlugin struct {
	Members string               `json:"members" dft:"" usage:"Raft members like 10.99.1.33:570,10.99.2.33:570."`
	Home    string               `json:"plugin.raft.home" dft:"${MESH_HOME}/mesh/raft/" usage:"Path to store raft snapshots and wal"`
	Host    *dragonboat.NodeHost `json:"-"`
}

func (that *raftPlugin) Ptt() *plugin.Ptt {
	return &plugin.Ptt{
		Name:  plugin.Raft,
		Flags: raftPlugin{},
		Create: func() plugin.Plugin {
			return that
		},
	}
}

func (that *raftPlugin) LoggerFactory(name string) logger.ILogger {
	return &rlog{name: name, level: log.INFO, ctx: mpc.Context()}
}

func (that *raftPlugin) Start(ctx context.Context, rt plugin.Runtime) {
	err := rt.Parse(that)
	if nil != err {
		log.Error(ctx, err.Error())
		return
	}
	if "" == that.Members {

	}
	members := map[uint64]string{}
	for idx, v := range strings.Split(that.Members, ",") {
		// key is the ReplicaID, ReplicaID is not allowed to be 0
		// value is the raft address
		members[uint64(idx+1)] = v
	}

	logger.SetLoggerFactory(that.LoggerFactory)
	// change the log verbosity
	logger.GetLogger("raft").SetLevel(logger.ERROR)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)

	replicaID := flag.Int("nodeid", 1, "ReplicaID to use")

	// config for raft
	// note the ShardID value is not specified here
	rc := config.Config{
		ReplicaID:          uint64(*replicaID),
		ElectionRTT:        5,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	home := filepath.Join(that.Home, fmt.Sprintf("%d", replicaID))
	hostConf := config.NodeHostConfig{
		DeploymentID:        0,
		NodeHostID:          "",
		WALDir:              home,
		NodeHostDir:         home,
		RTTMillisecond:      200,
		RaftAddress:         "",
		AddressByNodeHostID: false,
		ListenAddress:       "",
		EnableMetrics:       true,
		RaftEventListener:   nil,
		SystemEventListener: nil,
		MaxSendQueueSize:    0,
		MaxReceiveQueueSize: 0,
		NotifyCommit:        false,
		Gossip:              config.GossipConfig{},
		Expert: config.ExpertConfig{
			LogDBFactory:            tan.Factory,
			TransportFactory:        new(mt),
			Engine:                  config.GetDefaultEngineConfig(),
			LogDB:                   config.GetDefaultLogDBConfig(),
			TestGossipProbeInterval: time.Second * 10,
		},
	}
	that.Host, err = dragonboat.NewNodeHost(hostConf)
	if nil != err {
		log.Error(ctx, err.Error())
		return
	}
	if err = that.Host.StartOnDiskReplica(members, true, RStore, rc); nil != err {
		log.Error(ctx, "Failed to set cluster node storage, %s", err.Error())
		return
	}
	// start the first cluster
	// we use ExampleStateMachine as the IStateMachine for this cluster, its
	// behaviour is identical to the one used in the Hello World example.
	rc.ShardID = shardID1
	if err = that.Host.StartReplica(members, false, NewExampleStateMachine, rc); nil != err {
		log.Error(ctx, "Failed to add cluster node, %s", err.Error())
		return
	}
	// start the second cluster
	// we use SecondStateMachine as the IStateMachine for the second cluster
	rc.ShardID = shardID2
	if err := that.Host.StartReplica(members, false, NewSecondStateMachine, rc); nil != err {
		log.Error(ctx, "Failed to add cluster node, %s", err.Error())
		return
	}
	raftStopper := syncutil.NewStopper()
	consoleStopper := syncutil.NewStopper()
	ch := make(chan string, 16)
	consoleStopper.RunWorker(func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if nil != err {
				close(ch)
				return
			}
			if s == "exit\n" {
				raftStopper.Stop()
				// no data will be lost/corrupted if nodehost.Stop() is not called
				that.Stop(ctx, rt)
				return
			}
			ch <- s
		}
	})
	raftStopper.RunWorker(func() {
		// use NO-OP client session here
		// check the example in godoc to see how to use a regular client session
		cs1 := that.Host.GetNoOPSession(shardID1)
		cs2 := that.Host.GetNoOPSession(shardID2)
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				// remove the \n char
				msg := strings.Replace(strings.TrimSpace(v), "\n", "", 1)
				var err error
				// In this example, the strategy on how data is sharded across different
				// Raft groups is based on whether the input message ends with a "?".
				// In your application, you are free to choose strategies suitable for
				// your application.
				if strings.HasSuffix(msg, "?") {
					// user message ends with "?", make a proposal to update the second
					// raft group
					_, err = that.Host.SyncPropose(ctx, cs2, []byte(msg))
				} else {
					// message not ends with "?", make a proposal to update the first
					// raft group
					_, err = that.Host.SyncPropose(ctx, cs1, []byte(msg))
				}
				cancel()
				if nil != err {
					fmt.Fprintf(os.Stderr, "SyncPropose returned error %v\n", err)
				}
			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.Wait()
}

func (that *raftPlugin) Stop(ctx context.Context, rt plugin.Runtime) {
	if nil != that.Host {
		that.Host.Close()
	}
}
