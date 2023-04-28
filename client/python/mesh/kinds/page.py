#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Generic

from mesh.macro import T, index as idx, serializable


@serializable
class Page(Generic[T]):

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
    def total(self) -> int:
        return 0

    @idx(20)
    def next(self) -> bool:
        return False

    @idx(25)
    def data(self) -> T:
        return None
