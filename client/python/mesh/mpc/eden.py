#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import abstractmethod, ABC
from typing import Type, Any, List

from mesh.kinds import Service, Reference
from mesh.macro import mpi, T, spi, Inspector
from mesh.mpc.invoker import Execution


@spi("mesh")
class Eden(ABC):

    @abstractmethod
    def define(self, metadata: mpi, reference: Type[T]) -> T:
        """
        Define the reference object.

        :param reference:
        :param metadata: Object reference type.
        :param c: Meta custom annotation.
        :return: Reference proxy
        """
        pass

    @abstractmethod
    def refer(self, metadata: mpi, reference: Type[T], method: Inspector) -> Execution[Reference]:
        """
        Refer the service reference by method.

        :param metadata:
        :param reference:
        :param method:
        :return:
        """
        pass

    @abstractmethod
    def store(self, kind: Type[T], service: Any):
        """
        Store the service object.

        :param kind:
        :param service:
        :return:
        """
        pass

    @abstractmethod
    def infer(self, urn: str) -> Execution[Service]:
        """
        Infer the reference service by domain.

        :param urn:
        :return:
        """
        pass

    @abstractmethod
    def refer_types(self) -> List[Type[T]]:
        """
        Get all reference types.
        :return: All reference types.
        """
        pass

    @abstractmethod
    def infer_types(self) -> List[Type[T]]:
        """
        Get all service types.
        :return: All service types.
        """
        pass
