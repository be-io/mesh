#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.kinds.principal import Principal
from mesh.macro import index, serializable


@serializable
class Location(Principal):

    @index(10)
    def ip(self) -> str:
        return ""

    @index(15)
    def port(self) -> str:
        return ""

    @index(20)
    def host(self) -> str:
        return ""

    @index(25)
    def name(self) -> str:
        return ""
