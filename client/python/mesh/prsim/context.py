#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import mesh.log as log
import threading
from abc import ABC, abstractmethod
from enum import Enum
from mesh.kinds.location import Location, Principal
from mesh.macro import T
from typing import Any, Dict, Generic, Deque, Optional, List


class RunMode(Enum):
    ROUTINE = 1
    """
    正常模式
    """
    PERFORM = 2
    """
    评测模式
    """
    DEFENSE = 4
    """
    高防模式
    """
    DEBUG = 8
    """
    调试模式
    """
    LOAD_TEST = 16
    """
    压测模式
    """
    MOCK = 32
    """
    Mock模式
    """

    def is_debug(self) -> bool:
        return RunMode.matches(RunMode.DEBUG, self.value)

    def is_load_test(self) -> bool:
        return RunMode.matches(RunMode.LOAD_TEST, self.value)

    def is_routine(self) -> bool:
        return RunMode.matches(RunMode.ROUTINE, self.value)

    def is_perform(self) -> bool:
        return RunMode.matches(RunMode.PERFORM, self.value)

    def is_defense(self) -> bool:
        return RunMode.matches(RunMode.DEFENSE, self.value)

    def is_mock(self) -> bool:
        return RunMode.matches(RunMode.MOCK, self.value)

    @staticmethod
    def matches(mode: Any, code: int):
        return isinstance(mode, RunMode) and (mode.value & code) is code

    @staticmethod
    def get_by_code(code: int) -> "RunMode":
        for e in RunMode:
            if (e.value & code) is e.value:
                return e
        return RunMode.ROUTINE


class Key(Generic[T], ABC):
    """
    Context key.
    """

    def __init__(self, name: str):
        self.name = name

    @abstractmethod
    def get_if_absent(self) -> T:
        """ Get the attribute if present by key """
        pass

    @abstractmethod
    def map(self, fn) -> Optional[Any]:
        """ Map new value """
        pass

    @abstractmethod
    def if_present(self, fn):
        """ Consume attribute if present """
        pass

    @abstractmethod
    def or_else(self, v: T) -> T:
        """ Get attribute default if optional """
        pass

    @abstractmethod
    def is_present(self) -> bool:
        """ Check the attribute is present """
        pass


class Metadata(Enum):
    """
    https://www.rfc-editor.org/rfc/rfc7540#section-8.1.2
    """

    MESH_SPAN_ID = "mesh-span-id"
    MESH_TIMESTAMP = "mesh-timestamp"
    MESH_RUN_MODE = "mesh-run-mode"
    MESH_CONSUMER = "mesh-consumer"
    MESH_PROVIDER = "mesh-provider"
    MESH_URN = "mesh-urn"
    #
    MESH_INCOMING_HOST = "mesh-incoming-host"
    MESH_OUTGOING_HOST = "mesh-outgoing-host"
    MESH_INCOMING_PROXY = "mesh-incoming-proxy"
    MESH_OUTGOING_PROXY = "mesh-outgoing-proxy"
    MESH_SUBSET = "mesh-subset"
    # PTP
    MESH_VERSION = "x-ptp-version"
    MESH_TECH_PROVIDER_CODE = "x-ptp-tech-provider-code"
    MESH_TRACE_ID = "x-ptp-trace-id"
    MESH_TOKEN = "x-ptp-token"
    MESH_URI = "x-ptp-uri"
    MESH_FROM_NODE_ID = "x-ptp-source-node-id"
    MESH_FROM_INST_ID = "x-ptp-source-inst-id"
    MESH_TARGET_NODE_ID = "x-ptp-target-node-id"
    MESH_TARGET_INST_ID = "x-ptp-target-inst-id"
    MESH_SESSION_ID = "x-ptp-session-id"
    MESH_TOPIC = "x-ptp-topic"
    MESH_TIMEOUT = "x-ptp-timeout"

    def key(self) -> str:
        return self.value

    def set(self, attachments: Dict[str, str], v: str):
        if attachments and '' != v and v != attachments.get(self.key(), ''):
            attachments[self.key()] = v

    def get(self, attachments: Dict[str, str]) -> str:
        return attachments.get(self.key(), '')

    def append(self, attachments: List[Any], v: str):
        if attachments is not None and '' != v:
            attachments.append((self.key(), v))


class Context(ABC):
    """
    MPC invoke context.
    """

    @abstractmethod
    def get_trace_id(self) -> str:
        """
        Get the request trace id.
        :return: The request trace id.
        """
        pass

    @abstractmethod
    def get_span_id(self) -> str:
        """
        Get the request span id.
        :return: The request span id.
        """
        pass

    @abstractmethod
    def get_timestamp(self) -> int:
        """
        Get the request create time.
        :return: The request create time.
        """
        pass

    @abstractmethod
    def get_run_mode(self) -> RunMode:
        """
        Get the request run mode. RunMode
        :return: The request run mode.
        """
        pass

    @abstractmethod
    def get_urn(self) -> str:
        """
        Mesh resource uniform name. Like: create.tenant.omega.json.http2.lx000001.mpi.trustbe.net
        :return: Uniform name.
        """
        pass

    @abstractmethod
    def get_consumer(self) -> Location:
        """
        Get the consumer network principal.
        :return: Consumer network principal.
        """
        pass

    @abstractmethod
    def get_provider(self) -> Location:
        """
        Get the provider network principal.
        :return: Provider network principal.
        """
        pass

    @abstractmethod
    def get_attachments(self) -> Dict[str, str]:
        """
        Dispatch attachments.
        :return: Dispatch attachments.
        """
        pass

    @abstractmethod
    def get_principals(self) -> "Queue[Principal]":
        """
        Get the mpc broadcast network principals.
        :return: Broadcast principals.
        """
        pass

    @abstractmethod
    def get_attributes(self) -> Dict[str, Any]:
        """
        Get the context attributes. The attributes don't be serialized.
        :return: attributes
        """
        pass

    @abstractmethod
    def get_attribute(self, key: Key[T]) -> T:
        """
        Like getAttachments, but attribute wont be transfer in invoke chain.
        :param key: attribute key
        :return: attribute value
        """
        pass

    @abstractmethod
    def set_attribute(self, key: Key[T], value: T) -> None:
        """
        Like putAttachments, but attribute won't be transfer in invoke chain.
        :param key: attribute key
        :param value: attribute
        :return: None
        """
        pass

    @abstractmethod
    def rewrite_urn(self, urn: str) -> None:
        """
        Rewrite the urn.
        :param urn: urn
        :return:
        """
        pass

    @abstractmethod
    def rewrite_context(self, another: "Context") -> None:
        """
        Rewrite the context by another context.
        :param another: another context
        :return:
        """
        pass

    @abstractmethod
    def resume(self) -> "Context":
        """
        Open a new context.
        :return:
        """
        pass

    @abstractmethod
    def encode(self) -> Dict[str, str]:
        """
        Encode context to string.
        :return:
        """
        pass

    @abstractmethod
    def decode(self, attachments: Dict[str, str]) -> None:
        """
        Decode context from string.
        :return:
        """
        pass


class Queue(Deque, Generic[T]):
    def pop(self) -> T:
        if self.__len__() < 1:
            log.warn(f"{threading.currentThread().ident}({threading.currentThread().name}) pop empty queue.")
            return None
        return super().pop()

    def peek(self) -> T:
        if self.__len__() < 1:
            return None
        return self[-1]
