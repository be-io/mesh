#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Dict

from mesh.macro import index, serializable


@serializable
class Document:

    def __init__(self, metadata: Dict[str, str], content: str, timestamp: int):
        self.metadata = metadata
        self.content = content
        self.timestamp = timestamp

    @index(0)
    def metadata(self) -> Dict[str, str]:
        return self.metadata

    @index(5)
    def content(self) -> str:
        return self.content

    @index(10)
    def timestamp(self) -> int:
        return self.timestamp
