#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.kinds.environ import Environ
from mesh.kinds.license import License
from mesh.macro import index, serializable


@serializable
class CommerceLicense:

    @index(0)
    def cipher(self) -> str:
        return ""

    @index(5)
    def explain(self) -> License:
        return License()


@serializable
class CommerceEnviron:

    @index(0)
    def cipher(self) -> str:
        return ""

    @index(5)
    def explain(self) -> Environ:
        return Environ()

    @index(10)
    def node_key(self) -> str:
        return ""
