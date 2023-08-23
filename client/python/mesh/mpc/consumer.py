#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import abstractmethod, ABC
from concurrent.futures import Future
from typing import Any

from mesh.kinds import Reference
from mesh.macro import spi
from mesh.mpc.invoker import Execution


@spi(name="grpc")
class Consumer(ABC):
    HTTP = "http"
    GRPC = "grpc"
    TCP = "tcp"
    MQTT = "mqtt"

    """
    Service consumer with any protocol and codec.
    """

    @abstractmethod
    def __enter__(self):
        pass

    @abstractmethod
    def __exit__(self, exc_type, exc_val, exc_tb):
        pass

    @abstractmethod
    def start(self):
        """
        Start the mesh broker.
        :return:
        """
        pass

    @abstractmethod
    def consume(self, address: str, urn: str, execution: Execution[Reference], inbound: bytes, args: Any) -> Future:
        """
        Consume the input payload.
        :param address: Remote address.
        :param urn: Actual uniform resource domain name.
        :param execution: Service reference.
        :param inbound: Input arguments.
        :param args: Input typed arguments.
        :return: Output payload
        """
        pass
