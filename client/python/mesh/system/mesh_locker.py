#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import Locker


@spi("mesh")
class MeshLocker(Locker):

    def __init__(self):
        self.locker = ServiceProxy.default_proxy(Locker)

    def lock(self, rid: str, timeout: int) -> bool:
        return self.locker.lock(rid, timeout)

    def unlock(self, rid: str):
        return self.locker.unlock(rid)

    def read_lock(self, rid: str, timeout: int) -> bool:
        return self.locker.read_lock(rid, timeout)

    def read_unlock(self, rid: str):
        return self.locker.read_unlock(rid)
