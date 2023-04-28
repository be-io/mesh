#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod
from typing import List

from mesh.kinds import Registration
from mesh.macro import mpi, spi


@spi("mesh")
class Registry(ABC):

    @abstractmethod
    @mpi("mesh.registry.put")
    def register(self, registration: Registration):
        """
        Register metadata to mesh graph database.
        :param registration:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.registry.puts")
    def registers(self, registrations: List[Registration]):
        """
        Register metadata to mesh graph database.
        :param registrations:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.registry.remove")
    def unregister(self, registration: Registration):
        """
        Unregister metadata from mesh graph database.
        :param registration:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.registry.export")
    def export(self, kind: str) -> List[Registration]:
        """
        Export register metadata of mesh graph database.
        :param kind:
        :return:
        """
        pass
