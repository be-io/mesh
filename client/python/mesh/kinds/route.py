#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.macro import index, serializable


@serializable
class Route:

    @index(0)
    def node_id(self) -> str:
        """ 节点编号 """
        return ""

    @index(5)
    def inst_id(self) -> str:
        """ 机构编号 """
        return ""

    @index(10)
    def name(self) -> str:
        """ 节点名称 """
        return ""

    @index(15)
    def inst_name(self) -> str:
        """ 机构编号 """
        return ""

    @index(20)
    def address(self) -> str:
        """ 节点地址 """
        return ""

    @index(25)
    def describe(self) -> str:
        """ 节点描述 """
        return ""

    @index(30)
    def host_root(self) -> str:
        """ Host root certifications """
        return ""

    @index(35)
    def host_key(self) -> str:
        """ Host private certifications key """
        return ""

    @index(40)
    def host_crt(self) -> str:
        """ Host certification """
        return ""

    @index(45)
    def guest_root(self) -> str:
        """ Guest root certifications """
        return ""

    @index(50)
    def guest_key(self) -> str:
        """ Guest private certifications key """
        return ""

    @index(55)
    def guest_crt(self) -> str:
        """ Guest certifications key """
        return ""

    @index(60)
    def status(self) -> int:
        """ 状态 """
        return 0

    @index(65)
    def version(self) -> int:
        """ 状态 """
        return 0

    @index(70)
    def auth_code(self) -> str:
        """ 授权码 """
        return ""

    @index(75)
    def expire_at(self) -> int:
        """ 失效时间 """
        return 0

    @index(80)
    def extra(self) -> str:
        """ 扩展信息 """
        return ""

    @index(85)
    def create_at(self) -> str:
        """ 创建时间 """
        return ""

    @index(90)
    def create_by(self) -> str:
        """ 创建人 """
        return ""

    @index(95)
    def update_at(self) -> str:
        """ 更新时间 """
        return ""

    @index(100)
    def update_by(self) -> str:
        """ 跟新人 """
        return ""

    @index(105)
    def group(self) -> str:
        """ Network group. """
        return ""

    @index(110)
    def upstream(self) -> int:
        """ Upstream bandwidth. """
        return 0

    @index(115)
    def downstream(self) -> int:
        """ Downstream bandwidth. """
        return 0

    @index(120)
    def static_ip(self) -> str:
        """ Static public ip address """
        return ""

    @index(125)
    def proxy(self) -> str:
        """ Proxy endpoint in transport """
        return ""

    @index(130)
    def concurrency(self) -> int:
        """ MPC concurrency """
        return 0
