/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import "encoding/json"

type Tablet struct {
	// Served by which group.
	GroupId     uint32 `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3" json:"groupId,omitempty"`
	Predicate   string `protobuf:"bytes,2,opt,name=predicate,proto3" json:"predicate,omitempty"`
	Force       bool   `protobuf:"varint,3,opt,name=force,proto3" json:"force,omitempty"`
	OnDiskBytes int64  `protobuf:"varint,7,opt,name=on_disk_bytes,json=onDiskBytes,proto3" json:"on_disk_bytes,omitempty"`
	Remove      bool   `protobuf:"varint,8,opt,name=remove,proto3" json:"remove,omitempty"`
	// If true, do not ask zero to serve any tablets.
	ReadOnly bool   `protobuf:"varint,9,opt,name=read_only,json=readOnly,proto3" json:"readOnly,omitempty"`
	MoveTs   uint64 `protobuf:"varint,10,opt,name=move_ts,json=moveTs,proto3" json:"moveTs,omitempty"`
	// Estimated uncompressed size of tablet in bytes
	UncompressedBytes int64 `protobuf:"varint,11,opt,name=uncompressed_bytes,json=uncompressedBytes,proto3" json:"uncompressed_bytes,omitempty"`
}

type Group struct {
	Members      map[uint64]*Member `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Tablets      map[string]*Tablet `protobuf:"bytes,2,rep,name=tablets,proto3" json:"tablets,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SnapshotTs   uint64             `protobuf:"varint,3,opt,name=snapshot_ts,json=snapshotTs,proto3" json:"snapshot_ts,omitempty"`
	Checksum     uint64             `protobuf:"varint,4,opt,name=checksum,proto3" json:"checksum,omitempty"`
	CheckpointTs uint64             `protobuf:"varint,5,opt,name=checkpoint_ts,json=checkpointTs,proto3" json:"checkpoint_ts,omitempty"`
}

// Member stores information about RAFT group member for a single RAFT node.
// Note that each server can be serving multiple RAFT groups. Each group would
// have one RAFT node per server serving that group.
type Member struct {
	Id              uint64 `protobuf:"fixed64,1,opt,name=id,proto3" json:"id,omitempty"`
	GroupId         uint32 `protobuf:"varint,2,opt,name=group_id,json=groupId,proto3" json:"groupId,omitempty"`
	Addr            string `protobuf:"bytes,3,opt,name=addr,proto3" json:"addr,omitempty"`
	Leader          bool   `protobuf:"varint,4,opt,name=leader,proto3" json:"leader,omitempty"`
	AmDead          bool   `protobuf:"varint,5,opt,name=am_dead,json=amDead,proto3" json:"amDead,omitempty"`
	LastUpdate      uint64 `protobuf:"varint,6,opt,name=last_update,json=lastUpdate,proto3" json:"lastUpdate,omitempty"`
	Learner         bool   `protobuf:"varint,7,opt,name=learner,proto3" json:"learner,omitempty"`
	ClusterInfoOnly bool   `protobuf:"varint,13,opt,name=cluster_info_only,json=clusterInfoOnly,proto3" json:"clusterInfoOnly,omitempty"`
	ForceGroupId    bool   `protobuf:"varint,14,opt,name=force_group_id,json=forceGroupId,proto3" json:"forceGroupId,omitempty"`
}

type ClusterContext struct {
	Id         uint64 `protobuf:"fixed64,1,opt,name=id,proto3" json:"id,omitempty"`
	Group      uint32 `protobuf:"varint,2,opt,name=group,proto3" json:"group,omitempty"`
	Addr       string `protobuf:"bytes,3,opt,name=addr,proto3" json:"addr,omitempty"`
	SnapshotTs uint64 `protobuf:"varint,4,opt,name=snapshot_ts,json=snapshotTs,proto3" json:"snapshot_ts,omitempty"`
	IsLearner  bool   `protobuf:"varint,5,opt,name=is_learner,json=isLearner,proto3" json:"is_learner,omitempty"`
}

type ClusterBatch struct {
	Context *ClusterContext `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	Payload json.RawMessage `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

type HealthInfo struct {
	Instance    string   `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	Address     string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Status      string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Group       string   `protobuf:"bytes,4,opt,name=group,proto3" json:"group,omitempty"`
	Version     string   `protobuf:"bytes,5,opt,name=version,proto3" json:"version,omitempty"`
	Uptime      int64    `protobuf:"varint,6,opt,name=uptime,proto3" json:"uptime,omitempty"`
	LastEcho    int64    `protobuf:"varint,7,opt,name=lastEcho,proto3" json:"lastEcho,omitempty"`
	Ongoing     []string `protobuf:"bytes,8,rep,name=ongoing,proto3" json:"ongoing,omitempty"`
	Indexing    []string `protobuf:"bytes,9,rep,name=indexing,proto3" json:"indexing,omitempty"`
	EeFeatures  []string `protobuf:"bytes,10,rep,name=ee_features,json=eeFeatures,proto3" json:"ee_features,omitempty"`
	MaxAssigned uint64   `protobuf:"varint,11,opt,name=max_assigned,json=maxAssigned,proto3" json:"max_assigned,omitempty"`
}
