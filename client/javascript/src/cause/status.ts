/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Codeable} from "@/cause/codeable";

export class MeshCode extends Error implements Codeable {

    constructor(code: string, message: string) {
        super(message)
        this.code = code;
        this.message = message;
    }

    readonly code: string;
    readonly message: string

    getCode(): string {
        return this.code;
    }

    getMessage(): string {
        return this.message;
    }


}

export class Status {

    public static SUCCESS = new MeshCode("E0000000000", "请求成功")
    public static NOT_FOUND = new MeshCode("E0000000404", "请求资源不存在")
    public static SYSTEM_ERROR = new MeshCode("E0000000500", "系统异常")
    public static SERVICE_UNAVAILABLE = new MeshCode("E0000000503", "循环请求服务不可达")
    public static VALIDATE = new MeshCode("E0000000400", "请求非法")
    public static UNAUTHORIZED = new MeshCode("E0000000403", "请求资源未被授权")
    public static UNKNOWN = new MeshCode("E0000000520", "未知异常")
    public static COMPATIBLE = new MeshCode("E0000000600", "系统不兼容")
    public static TIMEOUT = new MeshCode("E0000000601", "请求超时")
    public static NO_PROVIDER = new MeshCode("E0000000602", "无服务实例")
    public static CRYPT_ERROR = new MeshCode("E0000000603", "数字证书校验异常")
    public static TOKEN_EXPIRE = new MeshCode("E0000000604", "节点授权码已过期")
    public static NET_EXPIRE = new MeshCode("E0000000605", "节点组网时间已过期")
    public static NET_DISABLE = new MeshCode("E0000000606", "对方节点已禁用网络")
    public static NET_UNAVAILABLE = new MeshCode("E0000000607", "网络不通")
    public static LICENSE_FORMAT_ERROR = new MeshCode("E0000000608", "License非法")
    public static LICENSE_EXPIRED = new MeshCode("E0000000609", "License已过期")
    public static MAX_REPLICAS_LIMIT = new MeshCode("E0000000610", "集群副本已达许可上限")
    public static MAX_COOP_LIMIT = new MeshCode("E0000000611", "合作方链接数已达许可上限")
    public static MAX_TENANT_LIMIT = new MeshCode("E0000000612", "租户数已达许可上限")
    public static MAX_USER_LIMIT = new MeshCode("E0000000613", "用户数已达许可上限")
    public static URN_NOT_PERMIT = new MeshCode("E0000000614", "接口未被许可调用")
    public static SIGNATURE_ERROR = new MeshCode("E0000000615", "证书签名非法")
    public static CRYPT_CODEC_ERROR = new MeshCode("E0000000616", "报文编解码异常")
    public static NO_SERVICE = new MeshCode("E0000000617", "下游版本不匹配服务不存在")
    public static NET_NOT_WEAVE = new MeshCode("E0000000618", "节点或机构未组网")
    public static ADDRESS_ERROR = new MeshCode("E0000000619", "地址非法或无法访问")
    public static UNEXPECTED_SYNTAX = new MeshCode("E0000000620", "规则语法错误")
    public static UNKNOWN_SCRIPT = new MeshCode("E0000000621", "规则不存在")
}