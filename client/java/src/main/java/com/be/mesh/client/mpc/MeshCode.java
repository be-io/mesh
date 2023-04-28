/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.prsim.Codeable;
import lombok.AllArgsConstructor;
import lombok.Getter;

/**
 * @author coyzeng@gmail.com
 */
@Getter
@AllArgsConstructor
public enum MeshCode implements Codeable {

    SUCCESS("E0000000000", "请求成功"),
    NOT_FOUND("E0000000404", "请求资源不存在"),
    SYSTEM_ERROR("E0000000500", "系统异常"),
    SERVICE_UNAVAILABLE("E0000000503", "循环请求服务不可达"),
    VALIDATE("E0000000400", "请求非法"),
    UNAUTHORIZED("E0000000403", "请求资源未被授权"),
    UNKNOWN("E0000000520", "未知异常"),
    COMPATIBLE("E0000000600", "系统不兼容"),
    TIMEOUT("E0000000601", "请求超时"),
    NO_PROVIDER("E0000000602", "无服务实例"),
    CRYPT_ERROR("E0000000603", "数字证书非法"),
    TOKEN_EXPIRE("E0000000604", "节点授权码已过期"),
    NET_EXPIRE("E0000000605", "节点组网时间已过期"),
    NET_DISABLE("E0000000606", "对方节点已禁用网络"),
    NET_UNAVAILABLE("E0000000607", "网络不通"),
    LICENSE_FORMAT_ERROR("E0000000608", "License非法"),
    LICENSE_EXPIRED("E0000000609", "License已过期"),
    MAX_REPLICAS_LIMIT("E0000000610", "集群副本已达许可上限"),
    MAX_COOP_LIMIT("E0000000611", "合作方链接数已达许可上限"),
    MAX_TENANT_LIMIT("E0000000612", "租户数已达许可上限"),
    MAX_USER_LIMIT("E0000000613", "用户数已达许可上限"),
    URN_NOT_PERMIT("E0000000614", "接口未被许可调用"),
    SIGNATURE_ERROR("E0000000615", "证书签名非法"),
    CRYPT_CODEC_ERROR("E0000000616", "报文编解码异常"),
    NO_SERVICE("E0000000617", "下游版本不匹配服务不存在"),
    NET_NOT_WEAVE("E0000000618", "节点或机构未组网"),
    ADDRESS_ERROR("E0000000619", "地址非法或无法访问"),
    UNEXPECTED_SYNTAX("E0000000620", "规则语法错误"),
    UNKNOWN_SCRIPT("E0000000621", "规则不存在"),
    ;

    private final String code;
    private final String message;
}
