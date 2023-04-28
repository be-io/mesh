#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod
from typing import Any, List

from mesh.kinds import Document, Paging, Page
from mesh.macro import mpi, spi


@spi("mesh")
class DataHouse(ABC):

    @abstractmethod
    @mpi("mesh.dh.writes")
    def writes(self, docs: List[Document]):
        """
        batch write log message
        :param docs:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.dh.write")
    def write(self, doc: Document):
        """
        write log message
        :param doc:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.dh.read")
    def read(self, index: Paging) -> Page[Any]:
        """
        write log message
        :param index:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.dh.indies")
    def indies(self, index: Paging) -> Page[Any]:
        """
        Export index list.
        :param index:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.dh.tables")
    def tables(self, index: Paging) -> Page[Any]:
        """
        Export table list.
        :param index:
        :return:
        """
        pass
