/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
	"strings"
)

const (
	MeshIDLength    = 6
	MeshNode        = "NODE"
	MeshInstitution = "INSTITUTION"
	NodeIdLength    = 15
	InstIdLength    = 18
)

// NodeID
// 节点ID由15位编码构成，构建规则：
// 厂商编码2位  + 5位扩展位 + 1位节点类型 + 6位节点序号 + 1位版本号
//  1. 前两位是厂商编码，LX（蓝象）
//  2. 5位扩展位：默认为00000
//  3. 1位节点类型：
//     a. GLab： 0
//     b. GAIA： 1
//     c. GLite：2
//     d. TEE：  3
type NodeID struct {
	LX            string // 2  Prefix
	Supplementary string // 5  Union number.
	Kind          string // 1  Kind
	SEQ           string // 6  Sequence number.
	Version       string // 1  Version.
}

func (that *NodeID) String() string {
	return fmt.Sprintf("%s%s%s%s%s", that.LX, that.Supplementary, that.Kind, that.SEQ, that.Version)
}

func NodeSEQ(nodeId string) string {
	if len(nodeId) == NodeIdLength {
		n, _ := FromNodeID(nodeId)
		return n.SEQ
	}
	if len(nodeId) == InstIdLength {
		n, _ := FromInstID(nodeId)
		return n.NodeSEQ
	}
	return nodeId
}

func FromNodeID(nodeId string) (*NodeID, error) {
	if len(nodeId) < NodeIdLength {
		return nil, cause.Errorf("NodeID %s is illegal, length should gather than 15.", nodeId)
	}
	return &NodeID{
		LX:            nodeId[0:2],
		Supplementary: nodeId[2:7],
		Kind:          nodeId[7:8],
		SEQ:           nodeId[8:14],
		Version:       nodeId[14:15],
	}, nil
}

func ApplyNodeID(kind string, seq string) *NodeID {
	return ApplyCNodeID(kind, seq, "LX")
}

func ApplyCNodeID(kind string, seq string, cname string) *NodeID {
	return &NodeID{LX: cname, Supplementary: "00000", Kind: kind, SEQ: seq, Version: "0"}
}

// InstitutionID
// 机构ID由18位编码构成，构建规则：
// JG + 机构类型(2位) + 一级机构编码(6位) +  二级机构编码(6位)  +  预留(1位) +  1位版本号
//
// 机构类型：
// 1. 00-蓝象
// 2. 01-普通机构，SaaS类、联盟类平台中成员机构
// 3. 02-SaaS类平台的管理机构，例如：常州、深圳等数据要素流通平台
// 4. 03-联盟类平台的盟主机构，例如：人行
type InstitutionID struct {
	JG            string // 2
	Kind          string // 2
	NodeSEQ       string // 6 same as node seq
	SEQ           string // 6
	Supplementary string // 1
	Version       string // 1
}

func (that *InstitutionID) String() string {
	return strings.ToUpper(fmt.Sprintf("%s%s%s%s%s%s",
		that.JG, that.Kind, that.NodeSEQ, that.SEQ, that.Supplementary, that.Version))
}

func (that *InstitutionID) MatchNode(nodeId string) bool {
	if nid, err := FromNodeID(nodeId); nil != err {
		return false
	} else {
		return nid.SEQ == that.NodeSEQ
	}
}

func (that *InstitutionID) Match(instId string) bool {
	if nid, err := FromInstID(instId); nil != err {
		return false
	} else {
		return nid.NodeSEQ == that.NodeSEQ
	}
}

func FromInstID(instId string) (*InstitutionID, error) {
	if len(instId) < InstIdLength {
		return nil, cause.Errorf("InstID %s is illegal.", instId)
	}
	return &InstitutionID{
		JG:            instId[0:2],
		Kind:          instId[2:4],
		NodeSEQ:       instId[4:10],
		SEQ:           instId[10:16],
		Supplementary: instId[16:17],
		Version:       instId[17:18],
	}, nil
}

func ApplyInstID(kind string, nodeSeq string, seq string) *InstitutionID {
	return &InstitutionID{JG: "JG", Kind: kind, NodeSEQ: nodeSeq, SEQ: seq, Supplementary: "0", Version: "0"}
}

func LXInstID(nodeSeq string, seq string) *InstitutionID {
	return ApplyInstID("00", nodeSeq, seq)
}

func AloneInstID(nodeSeq string, seq string) *InstitutionID {
	return ApplyInstID("01", nodeSeq, seq)
}

func SaaSInstID(nodeSeq string, seq string) *InstitutionID {
	return ApplyInstID("02", nodeSeq, seq)
}

func GroupInstID(nodeSeq string, seq string) *InstitutionID {
	return ApplyInstID("03", nodeSeq, seq)
}
