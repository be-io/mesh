#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import List

from mesh.kinds import Event
from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import Publisher


@spi("mesh")
class MeshPublisher(Publisher):

    def __init__(self):
        self.proxy = ServiceProxy.default_proxy(Publisher)

    def publish(self, events: List[Event]) -> List[str]:
        return self.proxy.publish(events)

    def broadcast(self, events: List[Event]) -> List[str]:
        return self.proxy.broadcast(events)
