#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import abstractmethod
from enum import Enum


class Codeable:
    """
    /**
     * See https://blue-elephant.yuque.com/blue-elephant/dg3hry/srnhhr.
     *
     * <pre>
     *
     * 1.5版本后，算子、算法、业务应用均往该规范靠.
     *
     * 错误码格式：
     *     E      0       00       000         0000
     *    前缀  异常等级  异常类型  模块或应用    业务码含义
     *
     * 总共10位错误码:
     * 1. 1位是等级，为0的等级直接对客
     * 2. 23位是异常类型
     * 3. 456位是模块或应用标志
     * 4. 78910位是自定义的业务码含义.
     *
     * 错误message可以是开发便于排查的，在网关层转化为产品定义的message，可配置.
     *
     * 错误码	模块
     * 000	   mesh
     * 001	   omega
     * 002	   asset
     * 003	   theta
     * 004	   cube
     * 005	   edge
     * 006     tensor
     * 007     base
     * </pre>
     *
     * @author coyzeng@gmail.com
     */
    """

    @abstractmethod
    def get_code(self) -> str:
        pass

    @abstractmethod
    def get_message(self) -> str:
        pass

    def match(self, codeable: "Codeable") -> bool:
        """
        Error code is matchable.
        """
        return self.matches(codeable.get_code())

    def matches(self, code: str) -> bool:
        if code is None:
            return False
        return str.__eq__(self.get_code(), code)


class MeshCode(Codeable, Enum):
    SUCCESS = "E0000000000", "请求成功"
    NOT_FOUND = "E0000000404", "请求资源不存在"
    SYSTEM_ERROR = "E0000000500", "系统异常"
    SERVICE_UNAVAILABLE = "E0000000503", "循环请求服务不可达"
    VALIDATE_ERROR = "E0000000400", "请求非法"
    UNAUTHORIZED = "E0000000403", "请求资源未被授权"
    UNKNOWN = "E0000000520", "未知异常"
    COMPATIBLE_ERROR = "E0000000600", "系统不兼容"
    TIMEOUT_ERROR = "E0000000601", "请求超时"
    NO_PROVIDER_ERROR = "E0000000602", "无服务实例"
    CRYPT_ERROR = "E0000000603", "数字证书校验异常"
    TOKEN_EXPIRE = "E0000000604", "节点授权码已过期"
    NET_EXPIRE = "E0000000605", "节点组网时间已过期"
    NET_DISABLE = "E0000000606", "对方节点已禁用网络"
    NET_UNAVAILABLE = "E0000000607", "网络不通"
    LICENSE_FORMAT_ERROR = "E0000000608", "License非法"
    LICENSE_EXPIRED = "E0000000609", "License已过期"
    MAX_REPLICAS_LIMIT = "E0000000610", "集群副本已达许可上限"
    MAX_COOP_LIMIT = "E0000000611", "合作方链接数已达许可上限"
    MAX_TENANT_LIMIT = "E0000000612", "租户数已达许可上限"
    MAX_USER_LIMIT = "E0000000613", "用户数已达许可上限"
    URN_NOT_PERMIT = "E0000000614", "接口未被许可调用"
    SIGNATURE_ERROR = "E0000000615", "证书签名非法"
    CRYPT_CODEC_ERROR = "E0000000616", "报文编解码异常"
    NO_SERVICE = "E0000000617", "下游版本不匹配服务不存在"
    NET_NOT_WEAVE = "E0000000618", "节点或机构未组网"
    ADDRESS_ERROR = "E0000000619", "地址非法或无法访问"
    UNEXPECTED_SYNTAX = "E0000000620", "规则语法错误"
    UNKNOWN_SCRIPT = "E0000000621", "规则不存在"

    def __init__(self, code: str, message: str) -> None:
        self.code = code
        self.message = message

    def __str__(self):
        return f'{self.code}({self.message})'

    def get_code(self) -> str:
        return self.value[0]

    def get_message(self) -> str:
        return self.value[1]
