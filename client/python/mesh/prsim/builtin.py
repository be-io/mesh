#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import ABC, abstractmethod
from typing import Dict, List

from mesh.kinds import Versions
from mesh.macro import spi, mpi


@spi("mesh")
class Builtin(ABC):

    @abstractmethod
    @mpi("${mesh.name}.builtin.doc")
    def doc(self, name: str, formatter: str) -> str:
        """
        Export the documents.
        :param name:
        :param formatter:
        :return:
        """
        pass

    @abstractmethod
    @mpi("${mesh.name}.builtin.version")
    def version(self) -> Versions:
        """
        Get the builtin application version.
        :return:
        """
        pass

    @abstractmethod
    @mpi("${mesh.name}.builtin.debug")
    def debug(self, features: Dict[str, str]):
        """
        LogLevel set the application log level.
        :return:
        """
        pass

    @abstractmethod
    @mpi("${mesh.name}.builtin.stats")
    def stats(self, features: List[str]) -> Dict[str, str]:
        """
        Health check stats.
        :param features:
        :return:
        """
        pass
