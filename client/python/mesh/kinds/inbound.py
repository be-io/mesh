#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any, Sequence

from mesh.macro import index, serializable


@serializable
class Inbound:

    @index(0)
    def arguments(self) -> Sequence[Any]:
        return []

    @index(1)
    def attachments(self) -> {}:
        return {}
