#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.macro import index, serializable


@serializable
class Lattice:

    @index(0)
    def zone(self) -> str:
        """zone"""
        return ''

    @index(5)
    def cluster(self) -> str:
        """cluster"""
        return ''

    @index(10)
    def cell(self) -> str:
        """cell"""
        return ''

    @index(15)
    def group(self) -> str:
        """group"""
        return ''

    @index(20)
    def address(self) -> str:
        """address"""
        return ''


class Environ:

    @index(0)
    def version(self) -> str:
        """Node certification version"""
        return ''

    @index(5)
    def node_id(self) -> str:
        """节点ID，所有节点按照标准固定分配一个全网唯一nodeId."""
        return ''

    @index(10)
    def inst_id(self) -> str:
        """每一个节点有一个初始的机构ID作为该节点的拥有者."""
        return ''

    @index(15)
    def inst_name(self) -> str:
        """一级机构名称."""
        return ''

    @index(20)
    def root_crt(self) -> str:
        """每一个节点内拥有一副证书用于通信，该证书可以被动态替换."""
        return ''

    @index(25)
    def root_key(self) -> str:
        """每一个节点内拥有一副证书用于通信，该证书可以被动态替换."""
        return ''

    @index(30)
    def node_crt(self) -> str:
        """节点许可证书"""
        return ''

    @index(35)
    def lattice(self) -> Lattice:
        """Mesh data center."""
        return Lattice()
