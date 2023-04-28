#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import List

from mesh.kinds import Registration
from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import Registry


@spi("mesh")
class MeshRegistry(Registry):

    def __init__(self):
        self.proxy = ServiceProxy.default_proxy(Registry)

    def register(self, registration: Registration):
        return self.proxy.register(registration)

    def registers(self, registrations: List[Registration]):
        return self.proxy.registers(registrations)

    def unregister(self, registration: Registration):
        return self.proxy.unregister(registration)

    def export(self, kind: str) -> List[Registration]:
        return self.proxy.export()
