#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import time
from typing import Dict, Any, Optional, Type

import mesh.tool as tool
from mesh.kinds import Principal, Location
from mesh.macro import ServiceLoader
from mesh.macro import T
from mesh.prsim import Context
from mesh.prsim import Network
from mesh.prsim.context import Key, RunMode, Queue


class MeshContext(Context):

    @staticmethod
    def create() -> Context:
        return MeshContext()

    def __init__(self):
        self.trace_id = tool.new_trace_id()
        self.span_id = tool.new_span_id("", 0)
        self.timestamp = int(time.time())
        self.run_mode = RunMode.ROUTINE
        self.urn = ''
        self.calls = 0
        self.consumer = Location()
        self.attachments = dict()
        self.attributes = dict()
        self.principals = Queue()

    def get_trace_id(self) -> str:
        return self.trace_id

    def get_span_id(self) -> str:
        return self.span_id

    def get_timestamp(self) -> int:
        return self.timestamp

    def get_run_mode(self) -> RunMode:
        return self.run_mode

    def get_urn(self) -> str:
        return self.urn

    def get_consumer(self) -> Location:
        if self.consumer and tool.required(self.consumer.node_id):
            return self.consumer
        return self.get_provider()

    def get_provider(self) -> Location:
        network = ServiceLoader.load(Network).get_default()
        environ = network.get_environ()
        location = Location()
        location.inst_id = environ.inst_id
        location.node_id = environ.node_id
        location.ip = tool.get_mesh_direct()
        location.host = tool.get_hostname()
        location.port = f"{tool.get_mesh_runtime().port}"
        location.name = tool.get_mesh_name()
        return location

    def get_attachments(self) -> Dict[str, str]:
        return self.attachments

    def get_principals(self) -> Queue[Principal]:
        return self.principals

    def get_attributes(self) -> Dict[str, Any]:
        return self.attributes

    def get_attribute(self, key: Key[T]) -> T:
        return self.attributes.get(key.name, None)

    def set_attribute(self, key: Key[T], value: T) -> None:
        self.attributes[key.name] = value

    def rewrite_urn(self, urn: str) -> None:
        self.urn = urn

    def rewrite_context(self, another: Context) -> None:
        if tool.required(another.get_trace_id()):
            self.trace_id = another.get_trace_id()
        if tool.required(another.get_span_id()):
            self.span_id = another.get_span_id()
        if another.get_timestamp() > 0:
            self.timestamp = another.get_timestamp()
        if RunMode.ROUTINE != another.get_run_mode():
            self.run_mode = another.get_run_mode()
        if tool.required(another.get_urn()):
            self.urn = another.get_urn()
        if tool.required(another.get_consumer()):
            self.consumer = another.get_consumer()
        if tool.required(another.get_attachments()):
            for key, value in another.get_attachments().items():
                self.attachments[key] = value
        if tool.required(another.get_attributes()):
            for key, value in another.get_attributes().items():
                self.attributes[key] = value
        if tool.required(another.get_principals()):
            for value in another.get_principals():
                self.principals.append(value)

    def resume(self) -> Context:
        self.calls = self.calls + 1
        context = MeshContext()
        context.rewrite_context(self)
        return context


class MeshKey(Key):

    def __init__(self, name: str, kind: Type[T]):
        self.name = name
        self.kind = kind

    def get_if_absent(self) -> T:
        pass

    def map(self, fn) -> Optional[Any]:
        pass

    def if_present(self, fn):
        pass

    def or_else(self, v: T) -> T:
        pass

    def is_present(self) -> bool:
        pass
