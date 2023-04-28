#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod
from typing import Type, Any

from mesh.kinds import Entity
from mesh.macro import mpi, spi, T


@spi("mesh")
class KV(ABC):

    @abstractmethod
    @mpi("mesh.kv.get")
    def get(self, key: str) -> Entity:
        """
        Get the value from kv store.
        :param key:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.kv.put")
    def put(self, key: str, value: Entity):
        """
        Put the value to kv store.
        :param key:
        :param value:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.kv.remove")
    def remove(self):
        """
        Remove the kv store.
        :return:
        """
        pass

    def get_with_type(self, key: str, kind: Type[T]) -> T:
        """
        Get by codec.
        :param key:
        :param kind:
        :return:
        """
        entity = self.get(key)
        return entity.try_read_object(kind) if entity else None

    def put_object(self, key: str, value: Any):
        """
        Put by codec.
        :param key:
        :param value:
        :return:
        """
        self.put(key, Entity.wrap(value))
