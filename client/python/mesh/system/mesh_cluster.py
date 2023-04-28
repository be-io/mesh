#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import Cluster


@spi("mesh")
class MeshCluster(Cluster):

    def __init__(self):
        self.proxy = ServiceProxy.default_proxy(Cluster)

    def election(self, buff: bytes) -> bytes:
        return self.proxy.election(buff)

    def is_leader(self) -> bool:
        return self.proxy.is_leader()
