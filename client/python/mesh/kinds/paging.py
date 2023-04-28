#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any, Dict

from mesh.macro import index as idx, serializable


@serializable
class Paging:

    @idx(0)
    def sid(self) -> str:
        return ''

    @idx(5)
    def index(self) -> int:
        return 0

    @idx(10)
    def limit(self) -> int:
        return 0

    @idx(15)
    def factor(self) -> Dict[str, Any]:
        return {}
