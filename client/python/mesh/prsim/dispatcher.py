#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod
from typing import Type, Dict, Any

from mesh.macro import spi, T


@spi("mesh")
class Dispatcher(ABC):

    @abstractmethod
    def reference(self, mpi: Type[T]) -> T:
        """
        Refer an generic dispatcher.
        :param mpi:
        :return:
        """
        pass

    @abstractmethod
    def invoke(self, urn: str, param: Dict[str, Any]) -> Any:
        """
        Dispatch generic invoke with urn.
        :param urn:
        :param param:
        :return:
        """
        pass

    @abstractmethod
    def invoke_generic(self, urn: str, param: Any) -> Any:
        """
        Dispatch generic invoke with urn and generic params.
        :param urn:
        :param param:
        :return:
        """
        pass
