#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from queue import Queue
from typing import Dict

from grpc import StreamStreamMultiCallable

from mesh.grpx.channels import GrpcChannels
from mesh.macro import spi
from mesh.prsim import Transport, Session
from mesh.tool import get_mesh_address


@spi("grpc")
class GrpcTransport(Transport):

    def __init__(self):
        self.channel = GrpcChannels()
        self.sessions: Dict[str, Session] = {}

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close(5000)
        self.channel.__exit__(exc_type, exc_val, exc_tb)

    def open(self, session_id: str, metadata: Dict[str, str]) -> Session:
        session = self.sessions.get(session_id, None)
        if not session:
            stream = self.channel.create_stream(get_mesh_address().any())
            session = self.sessions[session_id] = GrpcSession(stream)
        return session

    def close(self, timeout: int):
        for session_id, session in self.sessions.items():
            session.release(timeout)

    def roundtrip(self, payload: bytes, metadata: Dict[str, str]) -> bytes:
        """"""
        pass


@spi("grpc")
class GrpcSession(Session):

    def __init__(self, stream: StreamStreamMultiCallable):
        self.stream = stream
        self.inbounds = Queue()
        self.closed = False
        self.outbounds = self.stream(self.iter, 10000, [('urn', '')], None, True, False)

    def peek(self, topic: str = "") -> bytes:
        return self.outbounds.__next__()

    def pop(self, timeout: int, topic: str = "") -> bytes:
        return self.outbounds.__next__()

    def push(self, payload: bytes, metadata: Dict[str, str], topic: str = ""):
        self.inbounds.put(payload)

    def release(self, timeout: int, topic: str = ""):
        """"""
        pass

    def iter(self):
        while True:
            if self.closed:
                return
            payload = self.inbounds.get()
            if not payload:
                continue
            yield payload
