/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"github.com/opendatav/mesh/client/golang/log"
	"golang.org/x/net/context"
)

const EnvironVersion = "1.0.0"

// Environ is mesh fixed environ information.
type Environ struct {
	Version  string `index:"0" json:"version" xml:"version"  yaml:"version" comment:"节点证书版本."`
	NodeId   string `index:"5" json:"node_id" xml:"node_id"  yaml:"node_id" comment:"节点ID，所有节点按照标准固定分配一个全网唯一nodeId."`
	InstId   string `index:"10" json:"inst_id" xml:"inst_id" yaml:"inst_id" comment:"每一个节点有一个初始的机构ID作为该节点的拥有者."`
	InstName string `index:"15" json:"inst_name" xml:"inst_name" yaml:"inst_name" comment:"一级机构名称."`
	RootCrt  string `index:"20" json:"root_crt" xml:"root_crt"  yaml:"root_crt" comment:"节点根证书，每一个节点有一个统一的入网证书私钥，私钥用来作为节点内服务入网凭证."`
	RootKey  string `index:"25" json:"root_key" xml:"root_key" yaml:"root_key" comment:"节点根证书私钥."`
	NodeCrt  string `index:"30" json:"node_crt" xml:"node_crt"  yaml:"node_crt" comment:"节点许可证书"`
}

type Lattice struct {
	Zone    string `index:"0" json:"zone" xml:"zone" yaml:"zone" comment:"Zone"`
	Cluster string `index:"5" json:"cluster" xml:"cluster" yaml:"cluster" comment:"Cluster"`
	Cell    string `index:"10" json:"cell" xml:"cell" yaml:"cell" comment:"Cell"`
	Group   string `index:"15" json:"group" xml:"group" yaml:"group" comment:"Group"`
	Address string `index:"20" json:"address" xml:"address" yaml:"address" comment:"Address"`
}

func (that *Environ) ID(ctx context.Context) *NodeID {
	nodeId, err := FromNodeID(that.NodeId)
	if nil == err {
		return nodeId
	}
	log.Debug(ctx, err.Error())
	return &NodeID{SEQ: that.NodeId}
}
