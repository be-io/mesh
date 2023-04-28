#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from mesh.macro import index, serializable


@serializable
class Profile:

    @index(0)
    def data_id(self) -> str:
        return ""

    @index(5)
    def content(self) -> str:
        return ""
