#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Dict

from mesh.macro import index, serializable


@serializable
class Script:
    EXPRESSION = "EXPRESSION"
    VALUE = "VALUE"
    SCRIPT = "SCRIPT"

    @index(0)
    def code(self) -> str:
        return ""

    @index(5)
    def name(self) -> str:
        return ""

    @index(10)
    def desc(self) -> str:
        return ""

    @index(15)
    def kind(self) -> str:
        return ""

    @index(20)
    def expr(self) -> str:
        return ""

    @index(25)
    def attachment(self) -> Dict[str, str]:
        return {}
