#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import time
from abc import abstractmethod, ABC
from typing import Any, Type

from mesh.kinds import Entity
from mesh.kinds.entity import CacheEntity
from mesh.macro import spi, mpi, T


@spi(name="mesh")
class Cache(ABC):

    @abstractmethod
    @mpi("mesh.cache.get")
    def get(self, key: str) -> CacheEntity:
        """
        :param key:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.cache.put")
    def put(self, cell: CacheEntity) -> None:
        """
        :param cell:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.cache.remove")
    def remove(self, key: str):
        """
        :param key:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.cache.incr")
    def incr(self, key: str, value: int, duration: int) -> int:
        """
        :param key:
        :param value:
        :param duration:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.cache.decr")
    def decr(self, key: str, value: int, duration: int) -> int:
        """
        :param key:
        :param value:
        :param duration:
        :return:
        """
        pass

    def get_with_type(self, key: str, kind: Type[T]) -> T:
        """
        Get cache value with generic type.
        :param key:
        :param kind:
        :return:
        """
        ce = self.get(key)
        if ce is None:
            return None
        return ce.entity().try_read_object(kind)

    def put_with_duration(self, key: str, value: Any, duration: int):
        """
        Put value to cache with expire duration.
        :param key:
        :param value:
        :param duration:
        :return:
        """
        cell = CacheEntity()
        cell.version = "1.0.0"
        cell.entity = Entity.wrap(value)
        cell.timestamp = int(time.time() * 1000)
        cell.duration = duration
        cell.key = key
        self.put(cell)

    def compute_if_absent(self, key: str, kind: Type[T], duration: int, fn):
        """
        Default compute and put if absent.
        :param key:
        :param kind:
        :param duration:
        :param fn:
        :return:
        """
        value = self.get_with_type(key, kind)
        if value is not None:
            return value
        value = fn(key)
        if value is not None:
            self.put_with_duration(key, value, duration)
        return value
