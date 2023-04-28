#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import List

from mesh.macro import index, serializable


@serializable
class License:

    @index(0)
    def version(self) -> str:
        return ""

    @index(5)
    def level(self) -> int:
        return 0

    @index(10)
    def name(self) -> str:
        return ""

    @index(15)
    def create_by(self) -> str:
        return ""

    @index(20)
    def create_at(self) -> int:
        return 0

    @index(25)
    def active_at(self) -> int:
        return 0

    @index(30)
    def factors(self) -> List[str]:
        return []

    @index(35)
    def signature(self) -> str:
        return ""

    @index(40)
    def node_id(self) -> str:
        return ""

    @index(45)
    def inst_id(self) -> str:
        return ""

    @index(50)
    def server(self) -> str:
        return ""

    @index(55)
    def crt(self) -> str:
        return ""

    @index(60)
    def group(self) -> List[str]:
        return []

    @index(65)
    def replicas(self) -> int:
        return 0

    @index(70)
    def max_cooperators(self) -> int:
        return 0

    @index(75)
    def max_tenants(self) -> int:
        return 0

    @index(80)
    def max_users(self) -> int:
        return 0

    @index(85)
    def max_mills(self) -> int:
        return 0

    @index(90)
    def white_urns(self) -> List[str]:
        return []

    @index(95)
    def black_urns(self) -> List[str]:
        return []
