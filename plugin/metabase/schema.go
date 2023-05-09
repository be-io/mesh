/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import "time"

type KVEnt struct {
	// 配置KEY
	Key string `json:"key"`
	// 配置内容
	Value string `json:"value"`
	// 创建时间
	CreateAt time.Time `json:"create_at"`
	// 更新时间
	UpdateAt time.Time `json:"update_at"`
	// 创建人
	CreateBy string `json:"create_by"`
	// 更新人
	UpdateBy string `json:"update_by"`
}

type EdgeEnt struct {
	// 节点编号
	NodeID string `json:"node_id"`
	// 机构编号
	InstID string `json:"inst_id" badgerholdIndex:"InstId"`
	// 机构名称
	InstName string `json:"inst_name" badgerholdIndex:"InstName"`
	// 节点地址
	Address string `json:"address"`
	// 节点说明
	Describe string `json:"describe"`
	// 节点证书
	Certificate string `json:"certificate"`
	// 状态
	Status int32 `json:"status"`
	// 乐观锁版本
	Version int32 `json:"version"`
	// 授权码
	AuthCode string `json:"auth_code"`
	// 补充信息
	Extra string `json:"extra"`
	// 过期时间
	ExpireAt time.Time `json:"expire_at"`
	// 创建时间
	CreateAt time.Time `json:"create_at"`
	// 更新时间
	UpdateAt time.Time `json:"update_at"`
	// 创建人
	CreateBy string `json:"create_by"`
	// 更新人
	UpdateBy string `json:"update_by"`
	// 联盟中心节点机构id-多个用逗号分割
	Group string `json:"group"`
}

type Oauth2ClientEnt struct {
	// 客户端ID
	ID string `json:"id"`
	// 客户端名称
	Name string `json:"name"`
	// 客户端密钥
	Secret string `json:"secret"`
	// 客户端域名
	Domain string `json:"domain"`
	// 补充数据
	Data string `json:"data"`
}

type Oauth2TokenEnt struct {
	// 授权码
	Code string `json:"code"`
	// 准入TOKEN
	Access string `json:"access" badgerholdIndex:"Access"`
	// 刷新TOKEN
	Refresh string `json:"refresh" badgerholdIndex:"Refresh"`
	// 补充数据
	Data string `json:"data"`
	// 创建时间
	CreateAt time.Time `json:"create_at"`
	// 过期时间
	ExpireAt time.Time `json:"expire_at" badgerholdIndex:"ExpireAt"`
}

type SequenceEnt struct {
	// 序列号类型
	Kind string `json:"kind"`
	// 当前范围最小值
	Min int64 `json:"min"`
	// 当前范围最大值
	Max int64 `json:"max"`
	// 每次取号段大小
	Size int32 `json:"size"`
	// 序列号长度不足补零
	Length int32 `json:"length"`
	// 状态
	Status int32 `json:"status"`
	// 乐观锁版本
	Version int32 `json:"version"`
	// 创建时间
	CreateAt time.Time `json:"create_at"`
	// 更新时间
	UpdateAt time.Time `json:"update_at"`
}
