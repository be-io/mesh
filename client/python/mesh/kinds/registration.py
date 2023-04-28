#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any, List, Dict

from mesh.kinds.reference import Reference
from mesh.kinds.service import Service
from mesh.macro import index, serializable

Binding = "Binding"
MPS = "MPS"
Forward = "Forward"
Proxy = "Proxy"


@serializable
class Registration:
    METADATA = "metadata"
    PROXY = "proxy"
    SERVER = "server"
    COMPLEX = "complex"

    @index(0)
    def instance_id(self) -> str:
        return ''

    @index(5)
    def name(self) -> str:
        return ''

    @index(10)
    def kind(self) -> str:
        return ''

    @index(15)
    def address(self) -> str:
        return ''

    @index(20)
    def content(self) -> Any:
        return None

    @index(25)
    def timestamp(self) -> int:
        return 0

    @index(30)
    def attachments(self) -> Dict[str, str]:
        return {}


class Resource:

    @index(0)
    def references(self) -> List[Reference]:
        return []

    @index(5)
    def services(self) -> List[Service]:
        return []
