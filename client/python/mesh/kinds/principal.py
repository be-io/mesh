#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.macro import index, serializable


@serializable
class Principal:
    """Any fixed information of principal."""

    @index(0)
    def node_id(self) -> str:
        return ""

    @index(5)
    def inst_id(self) -> str:
        return ""
