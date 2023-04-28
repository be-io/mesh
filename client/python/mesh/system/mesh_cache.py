#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.kinds import CacheEntity
from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import Cache


@spi("mesh")
class MeshCache(Cache):

    def __init__(self):
        self.proxy = ServiceProxy.default_proxy(Cache)

    def get(self, key: str) -> CacheEntity:
        return self.proxy.get(key)

    def put(self, cell: CacheEntity) -> None:
        return self.proxy.put(cell)

    def remove(self, key: str):
        return self.proxy.remove(key)

    def incr(self, key: str, value: int, duration: int) -> int:
        return self.proxy.incr(key, value, duration)

    def decr(self, key: str, value: int, duration: int) -> int:
        return self.proxy.decr(key, value, duration)
