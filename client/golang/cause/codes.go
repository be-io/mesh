/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cause

import (
	"fmt"
)

func init() {
	var _ Codeable = new(MeshCode)
}

type MeshCode struct {
	Code    string
	Message string
}

func (that *MeshCode) GetCode() string {
	return that.Code
}

func (that *MeshCode) GetMessage() string {
	return that.Message
}

func (that *MeshCode) Error() error {
	return &Cause{
		code: that.Code,
		at:   Caller(2),
		err:  fmt.Errorf(that.Message),
	}
}

func Format(code Codeable) string {
	for _, coder := range codes {
		if coder.GetCode() == code.GetCode() {
			return coder.GetMessage()
		}
	}
	return code.GetMessage()
}

func From(code string) *MeshCode {
	for _, coder := range codes {
		if coder.GetCode() == code {
			return coder
		}
	}
	return Unknown
}

func Form(code int) string {
	return fmt.Sprintf("E0000000%d", code)
}

var codes = []*MeshCode{
	Success,
	NotFound,
	SystemError,
	Validate,
	Unauthorized,
	Unknown,
	Compatible,
	Timeout,
	NoProvider,
	CryptError,
	TokenExpire,
	NetExpire,
	NetDisable,
	NetUnavailable,
	LicenseFormatError,
	LicenseExpired,
	MaxReplicasLimit,
	MaxCoopLimit,
	MaxTenantLimit,
	MaxUserLimit,
	UrnNotPermit,
	SignatureError,
	CryptCodecError,
	NoService,
	NetNotWeave,
	AddressError,
	UnexpectedSyntax,
	UnknownScript,
	DataIDError,
}

var (
	Success            = &MeshCode{"E0000000000", "请求成功"}
	NotFound           = &MeshCode{"E0000000404", "请求资源不存在"}
	SystemError        = &MeshCode{"E0000000500", "系统异常"}
	ServiceUnavailable = &MeshCode{"E0000000503", "循环请求服务不可达"}
	Validate           = &MeshCode{"E0000000400", "请求非法"}
	Unauthorized       = &MeshCode{"E0000000403", "请求资源未被授权"}
	Unknown            = &MeshCode{"E0000000520", "未知异常"}
	Compatible         = &MeshCode{"E0000000600", "系统不兼容"}
	Timeout            = &MeshCode{"E0000000601", "请求超时"}
	NoProvider         = &MeshCode{"E0000000602", "无服务实例"}
	CryptError         = &MeshCode{"E0000000603", "数字证书校验异常"}
	TokenExpire        = &MeshCode{"E0000000604", "节点授权码已过期"}
	NetExpire          = &MeshCode{"E0000000605", "节点组网时间已过期"}
	NetDisable         = &MeshCode{"E0000000606", "对方节点已禁用网络"}
	NetUnavailable     = &MeshCode{"E0000000607", "网络不通"}
	LicenseFormatError = &MeshCode{"E0000000608", "License非法"}
	LicenseExpired     = &MeshCode{"E0000000609", "License已过期"}
	MaxReplicasLimit   = &MeshCode{"E0000000610", "集群副本已达许可上限"}
	MaxCoopLimit       = &MeshCode{"E0000000611", "合作方链接数已达许可上限"}
	MaxTenantLimit     = &MeshCode{"E0000000612", "租户数已达许可上限"}
	MaxUserLimit       = &MeshCode{"E0000000613", "用户数已达许可上限"}
	UrnNotPermit       = &MeshCode{"E0000000614", "接口未被许可调用"}
	SignatureError     = &MeshCode{"E0000000615", "证书签名非法"}
	CryptCodecError    = &MeshCode{"E0000000616", "报文编解码异常"}
	NoService          = &MeshCode{"E0000000617", "下游版本不匹配服务不存在"}
	NetNotWeave        = &MeshCode{"E0000000618", "节点或机构未组网"}
	AddressError       = &MeshCode{"E0000000619", "地址非法或无法访问"}
	UnexpectedSyntax   = &MeshCode{"E0000000620", "规则语法错误"}
	UnknownScript      = &MeshCode{"E0000000621", "规则不存在"}
	DataIDError        = &MeshCode{"E0000000622", "数据身份不存在或非法"}
	StartPending       = &MeshCode{"E0000000623", "系统启动中"}
)
