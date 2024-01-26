#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import math
from typing import Dict, Any, Callable

import mesh.log as log
import mesh.tool as tool
from mesh.kinds import Principal
from mesh.macro import spi
from mesh.mpc import ServiceProxy, Mesh
from mesh.prsim import Transport, Session, Metadata


@spi("mesh")
class MeshTransport(Transport):
    """
    Mesh transport is stateless transport depends on mesh server.
    """

    def __init__(self):
        self.sessions: Dict[str, Session] = dict()

    def open(self, session_id: str, metadata: Dict[str, str]) -> Session:
        skey = self.session_key(session_id, metadata)
        session = self.sessions.get(skey, None)
        if not session:
            attachments = Mesh.context().get_attachments()
            addr = Mesh.context().get_attribute(Mesh.REMOTE)
            principal = Mesh.context().get_principals().peek()
            session = MeshSession(session_id, metadata, attachments, addr, principal, lambda: self.finalize(skey))
            if self.sessions.__len__() < 500:
                self.sessions[skey] = session
        return session

    def close(self, timeout: int):
        for k, v in self.sessions.items():
            try:
                v.release(timeout, "")
            except BaseException as e:
                log.error("", e)
        self.sessions.clear()

    def roundtrip(self, payload: bytes, metadata: Dict[str, str]) -> bytes:
        """ bid """
        pass

    def finalize(self, session_key: str):
        if self.sessions.get(session_key, None):
            self.sessions.__delitem__(session_key)

    @staticmethod
    def session_key(session_id: str, metadata: Dict[str, str]) -> str:
        target_node_id = Metadata.MESH_TARGET_NODE_ID.get(metadata)
        target_inst_id = Metadata.MESH_TARGET_INST_ID.get(metadata)
        return f"{session_id}.{target_node_id}.{target_inst_id}"


@spi("mesh")
class MeshSession(Session):

    def __init__(self, session_id: str, metadata: Dict[str, str], attachments: Dict[str, str], address: str,
                 principal: Principal, finalizer: Callable):
        self.session_id = session_id
        self.metadata = metadata
        self.attachments = attachments
        self.address = address
        self.principal = principal
        self.session = SplitSession(ServiceProxy.default_proxy(Session))
        self.finalizer = finalizer

    def with_context(self, fn: Any, remote: bool, timeout: int) -> Any:
        return Mesh.context_safe(lambda: self.context_safe(fn, remote, timeout))

    def context_safe(self, fn: Any, remote: bool, timeout: int) -> Any:
        if self.attachments:
            for k, v in self.attachments.items():
                Mesh.context().get_attachments()[k] = v

        if self.metadata:
            for k, v in self.metadata.items():
                Mesh.context().get_attachments()[k] = v

        if tool.required(self.address):
            Mesh.context().set_attribute(Mesh.REMOTE, self.address)

        if timeout > 0:
            Mesh.context().set_attribute(Mesh.TIMEOUT, timeout)

        Metadata.MESH_SESSION_ID.set(Mesh.context().get_attachments(), self.session_id)

        if not remote:
            nk = Metadata.MESH_TARGET_NODE_ID.key()
            ik = Metadata.MESH_TARGET_INST_ID.key()
            if Mesh.context().get_attachments().__contains__(nk):
                Mesh.context().get_attachments().__delitem__(nk)
            if Mesh.context().get_attachments().__contains__(ik):
                Mesh.context().get_attachments().__delitem__(ik)
            return fn()

        target_node_id = Metadata.MESH_TARGET_NODE_ID.get(Mesh.context().get_attachments())
        target_inst_id = Metadata.MESH_TARGET_INST_ID.get(Mesh.context().get_attachments())
        if '' == target_node_id and '' == target_inst_id:
            return fn()
        try:
            principal = Principal()
            principal.node_id = target_node_id
            principal.inst_id = target_inst_id
            Mesh.context().get_principals().append(principal)
            return fn()
        finally:
            Mesh.context().get_principals().pop()

    def peek(self, topic: str = "") -> bytes:
        return self.with_context(lambda: self.session.peek(topic), False, 0)

    def pop(self, timeout: int, topic: str = "") -> bytes:
        return self.with_context(lambda: self.session.pop(timeout, topic), False, timeout)

    def push(self, payload: bytes, metadata: Dict[str, str], topic: str = ""):
        return self.with_context(lambda: self.session.push(payload, metadata, topic), True, 0)

    def release(self, timeout: int, topic: str = ""):
        self.finalizer()
        return self.with_context(lambda: self.session.release(timeout, topic), False, timeout)


class SplitSession(Session):

    def __init__(self, session: Session):
        self.session = session

    def peek(self, topic: str = "") -> bytes:
        buff = bytes()
        packet_size = self.session.peek(self.split_key(topic, 0)).decode()
        if "" == packet_size:
            return buff
        for idx in range(int(packet_size)):
            buff += self.session.peek(self.split_key(topic, idx + 1))
        return buff

    def pop(self, timeout: int, topic: str = "") -> bytes:
        packet_size = int(self.session.pop(timeout, self.split_key(topic, 0)).decode())
        buff = bytes()
        for idx in range(packet_size):
            buff += self.session.pop(timeout, self.split_key(topic, idx + 1))
        return buff

    def push(self, payload: bytes, metadata: Dict[str, str], topic: str = ""):
        packet_length = tool.get_packet_size()
        packet_size = math.ceil(payload.__len__() / packet_length)
        self.session.push(f"{packet_size}".encode(), metadata, self.split_key(topic, 0))
        for idx in range(packet_size):
            buff = payload[packet_length * idx:min(packet_length * (idx + 1), payload.__len__())]
            self.session.push(buff, metadata, self.split_key(topic, idx + 1))

    def release(self, timeout: int, topic: str = ""):
        return self.session.release(timeout, topic)

    @staticmethod
    def split_key(topic: str, idx: int, ) -> str:
        return f"{topic}-mesh-split-{idx}"
