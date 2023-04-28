#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod
from typing import Dict

from mesh.kinds.captcha import Captcha
from mesh.macro import spi, mpi


@spi("mesh")
class Graphics(ABC):

    @abstractmethod
    @mpi("mesh.graphics.captcha.apply")
    def captcha(self, kind: str, features: Dict[str, str]) -> Captcha:
        """
        Apply a graphics captcha.
        :param kind:
        :param features:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.graphics.captcha.verify")
    def verify(self, mno: str, value: str) -> bool:
        """
        Verify a graphics captcha value.
        :param mno:
        :param value:
        :return:
        """
        pass
