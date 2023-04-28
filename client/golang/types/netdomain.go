/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

// Domain
//
// create table mesh_domain (
//  `instance_id` varchar(255) not null default '' comment '服务实例ID',
//  `node_id` varchar(255) not null default '' comment '节点编号',
//  `inst_id` varchar(255) not null default '' comment '机构编号',
//  `zone` varchar(255) not null default '' comment '区，保留',
//  `cluster` varchar(255) not null default '' comment '集群，保留',
//  `unit` varchar(255) not null default '' comment '单元，保留',
//  `group` varchar(255) not null default '' comment '组，保留',
//  `address` varchar(255) not null default '' comment '服务地址',
//  `urn` varchar(255) not null default '' comment '资源统一域名',
//  `namespace` varchar(255) not null default '' comment '资源命名空间',
//  `name` varchar(255) not null default '' comment '资源名',
//  `version` varchar(255) not null default '' comment '版本',
//  `alias` varchar(255) not null default '' comment '别名',
//  `kind` varchar(255) not null default '' comment '类型',
//  `proto` varchar(255) not null default '' comment '网络协议',
//  `codec` varchar(255) not null default '' comment '序列化协议',
//  `timeout` int not null default 3000 comment '默认超时时间',
//  `asyncable` int not null default 1 comment '是否支持异步',
//  `args` varchar(255) not null default '' comment '扩展参数，业务路由',
//  primary key `instance_id`
//)
type Domain struct {

	// 实例ID
	InstanceID string `index:"0" json:"instance_id" xml:"instance_id" yaml:"instance_id" comment:"实例ID"`
	// 节点编号
	NodeID string `index:"5" json:"node_id" xml:"node_id" yaml:"node_id" comment:"节点编号"`
	// 机构编号
	InstID string `index:"10" json:"inst_id" xml:"inst_id" yaml:"inst_id" comment:"机构编号"`
	// 区，保留
	Zone string `index:"15" json:"zone" xml:"zone" yaml:"zone" comment:"区，保留"`
	// 集群，保留
	Cluster string `index:"20" json:"cluster" xml:"cluster" yaml:"cluster" comment:"集群，保留"`
	// 单元，保留
	Unit string `index:"25" json:"unit" xml:"unit" yaml:"unit" comment:"单元，保留"`
	// 组，保留
	Group string `index:"30" json:"group" xml:"group" yaml:"group" comment:"组，保留"`
	// 服务地址
	Address string `index:"35" json:"address" xml:"address" yaml:"address" comment:"服务地址"`
	// 资源统一域名
	URN string `index:"40" json:"urn" xml:"urn" yaml:"urn" comment:"资源统一域名"`
	// 资源命名空间
	Namespace string `index:"45" json:"namespace" xml:"namespace" yaml:"namespace" comment:"资源命名空间"`
	// 资源名
	Name string `index:"50" json:"name" xml:"name" yaml:"name" comment:"资源名"`
	// 版本
	Version string `index:"55" json:"version" xml:"version" yaml:"version" comment:"版本"`
	// 别名
	Alias string `index:"60" json:"alias" xml:"alias" yaml:"alias" comment:"别名"`
	// 类型
	Kind string `index:"65" json:"kind" xml:"kind" yaml:"kind" comment:"类型"`
	// 网络协议
	Proto string `index:"70" json:"proto" xml:"proto" yaml:"proto" comment:"网络协议"`
	// 网络协议
	Codec string `index:"75" json:"codec" xml:"codec" yaml:"codec" comment:"序列化协议"`
	// 默认超时时间
	Timeout string `index:"80" json:"timeout" xml:"timeout" yaml:"timeout" comment:"默认超时时间"`
	// 是否支持异步
	Asyncable bool `index:"85" json:"asyncable" xml:"asyncable" yaml:"asyncable" comment:"是否支持异步"`
	// 扩展参数，业务路
	Args bool `index:"90" json:"args" xml:"args" yaml:"args" comment:"扩展参数，业务路"`
	// 创建时间
	CreateAt string `index:"95" json:"create_at" xml:"create_at" yaml:"create_at" comment:"创建时间"`
	// 更新时间
	UpdateAt string `index:"100" json:"update_at" xml:"update_at" yaml:"update_at" comment:"更新时间"`
	// 创建人
	CreateBy string `index:"105" json:"create_by" xml:"create_by" yaml:"create_by" comment:"创建人"`
	// 更新人
	UpdateBy string `index:"110" json:"update_by" xml:"update_by" yaml:"update_by" comment:"更新人"`
}
