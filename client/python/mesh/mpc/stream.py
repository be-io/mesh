#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import List, Dict, Any

import mesh.tool as tool
from mesh.cause import ValidationException
from mesh.context import Mesh
from mesh.kinds import Principal
from mesh.macro import T, spi, Proxy, InvocationHandler
from mesh.prsim import Routable


class MeshRouter(InvocationHandler):

    def __init__(self, reference: Any, principal: Principal, attachments: Dict[str, str], address: str):
        self.reference = reference
        self.principal = principal
        self.attachments = attachments
        self.address = address

    def invoke(self, proxy: Any, method: Any, *args, **kwargs):
        return Mesh.context_safe(lambda: self.safe_invoke(proxy, method, *args, **kwargs))

    def safe_invoke(self, proxy: Any, method: Any, *args, **kwargs):
        overrides = self.get_override_attachments()
        routable = tool.required(self.principal.inst_id) or tool.required(self.principal.node_id)
        try:
            if tool.required(self.address):
                Mesh.context().set_attribute(Mesh.REMOTE, self.address)
            if tool.required(self.attachments):
                for (k, v) in self.attachments.items():
                    Mesh.context().get_attachments().__setitem__(k, v)
            if routable:
                Mesh.context().get_principals().append(self.principal)

            executor = getattr(self.reference, method.__name__, method)
            if executor == method:
                return executor(self.reference, *args, **kwargs)
            return executor(*args, **kwargs)
        finally:
            for (k, v) in overrides.items():
                Mesh.context().get_attachments().__setitem__(k, v)
            if routable:
                Mesh.context().get_principals().pop()

    def get_override_attachments(self):
        if tool.optional(self.attachments):
            return {}

        overrides = {}
        for (k, v) in self.attachments.items():
            overrides[k] = Mesh.context().get_attachments().get(k, '')

        return overrides


@spi("mesh")
class MeshRoutable(Routable):

    def __init__(self, reference: T = None, attachments: Dict[str, str] = None, address: str = ""):
        self.ref = reference
        self.attachments = attachments if attachments is not None else {}
        self.address = address

    def within(self, key: str, value: str) -> Routable[T]:
        return self.with_map({key: value})

    def with_map(self, attachments: Dict[str, str]) -> Routable[T]:
        if not attachments or attachments.__len__() < 1:
            return self
        kvs = {}
        if self.attachments.__len__() > 0:
            for (k, v) in self.attachments.items():
                kvs[k] = v
        for (k, v) in attachments.items():
            kvs[k] = v

        return MeshRoutable(self.ref, kvs)

    def with_address(self, address: str) -> Routable[T]:
        return MeshRoutable(self.ref, self.attachments, address)

    def local(self) -> T:
        interfaces = Proxy.get_interfaces(self.ref.__class__)
        return Proxy(interfaces, MeshRouter(self.ref, Principal(), self.attachments, self.address))

    def any(self, principal: Principal) -> T:
        if tool.optional(principal.node_id) and tool.optional(principal.inst_id):
            raise ValidationException("Route key both cant be empty.")
        interfaces = Proxy.get_interfaces(self.ref.__class__)
        return Proxy(interfaces, MeshRouter(self.ref, principal, self.attachments, self.address))

    def any_inst(self, inst_id: str) -> T:
        principal = Principal()
        principal.inst_id = inst_id
        return self.any(principal)

    def many(self, principals: List[Principal]) -> List[T]:
        if tool.optional(principals):
            return []
        references = []
        for principal in principals:
            references.append(self.any(principal))
        return references

    def many_inst(self, inst_ids: List[str]) -> List[T]:
        if tool.optional(inst_ids):
            return []
        principals = []
        for inst_id in inst_ids:
            principal = Principal()
            principal.inst_id = inst_id
            principals.append(principal)
        return self.many(principals)
