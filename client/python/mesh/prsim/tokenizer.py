#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import ABC, abstractmethod

from mesh.macro import spi, mpi


@spi("mesh")
class Tokenizer(ABC):

    @abstractmethod
    @mpi("mesh.trust.apply")
    def apply(self, kind: str, duration: int) -> str:
        """
        Apply a node token.
        :param kind:
        :param duration:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.trust.verify")
    def verify(self, token: str) -> bool:
        """
        Verify some token verifiable.
        :param token:
        :return:
        """
        pass
