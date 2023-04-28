#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any, List

from mesh.kinds import Document, Paging, Page
from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import DataHouse


@spi("mesh")
class MeshDataHouse(DataHouse):

    def __init__(self):
        self.proxy = ServiceProxy.default_proxy(DataHouse)

    def writes(self, docs: List[Document]):
        return self.proxy.writes(docs)

    def write(self, doc: Document):
        return self.proxy.write(doc)

    def read(self, index: Paging) -> Page[Any]:
        return self.proxy.read(index)

    def indies(self, index: Paging) -> Page[Any]:
        return self.proxy.indies(index)

    def tables(self, index: Paging) -> Page[Any]:
        return self.proxy.tables(index)
