/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/types"
)

var IInterconnection = (*Interconnection)(nil)

// Interconnection spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Interconnection interface {

	// Instx apply institutions by contract id.
	// @MPI("mesh.inc.instx")
	Instx(ctx context.Context, contractId string) ([]*types.Institution, error)

	// ContractIds infer the contract ids by institution ids.
	// @MPI("mesh.inc.contract.ids")
	ContractIds(ctx context.Context, instIds []string) ([]string, error)

	// Contract /v1/interconn/node/contract/query
	// APPLIED已申请
	// APPROVED已授权
	// REJECTED已拒绝
	// TERMINATED已解除
	// @MPI("mesh.inc.contract")
	Contract(ctx context.Context, req *IncContractID) (*IncState, error)

	// Describe /v1/interconn/node/query  /v1/platform/node/query
	// @MPI(name = "mesh.inc.describe", flags = 8)
	Describe(ctx context.Context, req *IncNodeID) (*IncNode, error)

	// Weave /v1/interconn/node/contract/apply
	// @MPI(name = "mesh.inc.weave", flags = 8)
	Weave(ctx context.Context, req *IncNode) (*IncContractID, error)

	// Ack /v1/interconn/node/contract/confirm
	// @MPI(name = "mesh.inc.ack", flags = 8)
	Ack(ctx context.Context, req *IncAck) (*IncOption, error)

	// Abort /v1/interconn/node/contract/terminate
	// @MPI(name = "mesh.inc.abort", flags = 8)
	Abort(ctx context.Context, req *IncContractID) (*IncOption, error)

	// Refresh /v1/interconn/node/update /v1/platform/node/update
	// @MPI(name = "mesh.inc.refresh", flags = 8)
	Refresh(ctx context.Context, req *IncNode) (*IncOption, error)

	// Probe /v1/interconn/node/health
	// 节点健康状态。直接返回ok
	// @MPI(name = "mesh.inc.probe", flags = 8)
	Probe(ctx context.Context, req *IncOption) (*IncState, error)
}

type IncContractID struct {
	ContractID string `json:"contract_id"`
}

type IncState struct {
	Status string `json:"status"` // APPLIED已申请 APPROVED已授权 REJECTED已拒绝 TERMINATED已 解除
}
type IncNodeID struct {
	NodeId string `json:"node_id"`
}

type IncNode struct {
	NodeId         string `json:"node_id"`         // 合作方的节点ID
	Name           string `json:"name"`            // 节点名称
	Institution    string `json:"institution"`     // 节点所属机构
	InstId         string `json:"inst_id"`         // 机构Id
	System         string `json:"system"`          // 技术服务提供系统
	SystemVersion  string `json:"system_version"`  // 系统版本
	Address        string `json:"address"`         // 节点服务地址
	Description    string `json:"description"`     // 节点说明 optional
	AuthType       string `json:"auth_type"`       // 认证方式，枚举值：SHA256_RSA、 SHA256_ECDSA、CERT等
	AuthCredential string `json:"auth_credential"` // 凭证内容：公钥值、证书内容等
	ExpiredTime    int64  `json:"expired_time"`    // 合约过期时间 optional
	Status         string `json:"status"`          // 授权状态 枚举值：修改状态类型，和流程审批一致 APPLIED已申请 APPROVED已授权 REJECTED已拒绝 TERMINATED已解除
	ContractID     string `json:"contract_id"`
}

type IncAck struct {
	ContractID     string `json:"contract_id"`     // 合约id，全局唯一
	Status         string `json:"status"`          // 授权状态 枚举值：修改状态类型，和流程审批一致 APPLIED已申请 APPROVED已授权 REJECTED已拒绝 TERMINATED已解除
	AuthType       string `json:"auth_type"`       // 认证方式，枚举值：SHA256_RSA、 SHA256_ECDSA、CERT等
	AuthCredential string `json:"auth_credential"` // 凭证内容：公钥值、证书内容等
	ExpiredTime    int64  `json:"expired_time"`    // 合约过期时间 optional
}

type IncOption struct {
}
