#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import ABC, abstractmethod
from typing import Generic, Any, Dict, List

from mesh.kinds import Script, Paging, Page
from mesh.macro import mpi, spi, T, A


@spi("mesh")
class Evaluator(ABC, Generic[T, A]):
    """
    Evaluate engine.
    """

    @abstractmethod
    @mpi("mesh.eval.compile")
    def compile(self, script: Script) -> str:
        """
        Compile the named rule.
        :param script:
        :return:
        """

    @abstractmethod
    @mpi("mesh.eval.exec")
    def exec(self, code: str, args: Any, dft: str) -> str:
        """
        Exec the script with name.
        :param code:
        :param args:
        :param dft:
        :return:
        """

    @abstractmethod
    @mpi("mesh.eval.dump")
    def dump(self, feature: Dict[str, str]) -> List[Script]:
        """
        Dump the scripts.
        :param feature:
        :return:
        """

    @abstractmethod
    @mpi("mesh.eval.index")
    def index(self, index: Paging) -> Page[Script]:
        """
        Dump the scripts.
        :param index:
        :return:
        """
