#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod
from typing import Dict, List

from mesh.macro import spi


@spi("mesh")
class Hodor(ABC):

    @abstractmethod
    def stats(self, features: List[str]) -> Dict[str, str]:
        """
        Collect the system, application, process or thread status etc.
        :param features: customized parameters
        :return: quota stat
        """
        pass

    @abstractmethod
    def debug(self, features: Dict[str, str]):
        """
        Set debug features.
        :param features:
        :return:
        """
        pass
