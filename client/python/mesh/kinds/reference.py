#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.macro import index, serializable


@serializable
class Reference:

    @index(0)
    def urn(self) -> str:
        return ''

    @index(5)
    def namespace(self) -> str:
        return ''

    @index(10)
    def name(self) -> str:
        return ''

    @index(15)
    def version(self) -> str:
        return ''

    @index(20)
    def proto(self) -> str:
        return ''

    @index(25)
    def codec(self) -> str:
        return ''

    @index(30)
    def flags(self) -> int:
        return 0

    @index(35)
    def timeout(self) -> int:
        return 5000

    @index(40)
    def retries(self) -> int:
        return 5

    @index(45)
    def node(self) -> str:
        return ''

    @index(50)
    def inst(self) -> str:
        return ''

    @index(55)
    def zone(self) -> str:
        return ''

    @index(60)
    def cluster(self) -> str:
        return ''

    @index(65)
    def cell(self) -> str:
        return ''

    @index(70)
    def group(self) -> str:
        return ''

    @index(75)
    def address(self) -> str:
        return ''
