#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any

from mesh.macro import index, serializable


@serializable
class Outbound:

    @index(0)
    def code(self) -> str:
        return ''

    @index(1)
    def message(self) -> str:
        return ''

    @index(2)
    def cause(self) -> str:
        return ''

    @index(3)
    def content(self) -> Any:
        return None
