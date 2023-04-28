#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from mesh.macro import index, serializable


@serializable
class Captcha:

    @index(0)
    def mno(self) -> str:
        return ""

    @index(5)
    def kind(self) -> str:
        return ""

    @index(10)
    def mime(self) -> bytes:
        return b''

    @index(15)
    def text(self) -> str:
        return ""
