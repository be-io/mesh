#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import List

from mesh.context import MeshKey
from mesh.kinds import Route, Environ, Versions, Paging, Page, Institution
from mesh.macro import spi
from mesh.mpc import ServiceProxy, Mesh
from mesh.prsim import Network


@spi("mesh")
class MeshNetwork(Network):

    def __init__(self):
        self.proxy = ServiceProxy.default_proxy(Network)
        self.environ: Environ = None
        self.environ_key = MeshKey("mesh-environ", Environ)

    def get_environ(self) -> Environ:
        if self.environ:
            return self.environ
        if Mesh.context().get_attribute(self.environ_key):
            return Mesh.context().get_attribute(self.environ_key)
        environ = Environ()
        environ.node_id = "LX0000000000000"
        environ.inst_id = "JG0000000000000000"
        Mesh.context().set_attribute(self.environ_key, environ)

        # self.environ = Mesh.context_safe(lambda: self.get_environ_safe())
        # return self.environ
        return environ

    def get_environ_safe(self) -> Environ:
        Mesh.context().get_principals().clear()
        return self.proxy.get_environ()

    def accessible(self, route: Route) -> bool:
        return self.proxy.accessible(route)

    def refresh(self, routes: List[Route]):
        return self.proxy.refresh(routes)

    def get_route(self, node_id: str) -> Route:
        return self.proxy.get_route(node_id)

    def get_routes(self) -> List[Route]:
        return self.proxy.get_routes()

    def get_domains(self) -> List[Route]:
        return self.proxy.get_domains()

    def put_domains(self, domains: List[Route]):
        return self.proxy.put_domains(domains)

    def weave(self, route: Route) -> None:
        return self.proxy.weave(route)

    def ack(self, route: Route) -> None:
        return self.proxy.ack(route)

    def disable(self, node_id: str) -> None:
        return self.proxy.disable(node_id)

    def enable(self, node_id: str) -> None:
        return self.proxy.enable(node_id)

    def index(self, index: Paging) -> Page[Route]:
        return self.proxy.index(index)

    def version(self, node_id: str) -> Versions:
        return self.proxy.version(node_id)

    def instx(self, index: Paging) -> Page[Institution]:
        return self.proxy.instx(index)

    def instr(self, institutions: List[Institution]) -> str:
        return self.proxy.instr(institutions)

    def ally(self, node_ids: List[str]) -> None:
        self.proxy.ally(node_ids)

    def disband(self, node_ids: List[str]) -> None:
        self.proxy.disband(node_ids)
