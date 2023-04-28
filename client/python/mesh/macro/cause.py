#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#


from mesh.macro.index import index


class Cause:

    @index(0)
    def name(self) -> str:
        return ''

    @index(5)
    def pos(self) -> str:
        return ''

    @index(10)
    def text(self) -> str:
        return ''

    @index(15)
    def buff(self) -> bytes:
        return b''

    @staticmethod
    def of(e: BaseException) -> "Cause":
        return Cause()

    @staticmethod
    def of_cause(code: str, message: str, cause: "Cause") -> BaseException:
        raise BaseException("")
