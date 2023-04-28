#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import abstractmethod
from typing import Any

from mesh.macro import spi


@spi(name="grpc")
class Provider:
    """

    """

    @abstractmethod
    def start(self, address: str, tc: Any):
        """
        Start the mesh broker.
        :param address:
        :param tc:
        :return:
        """
        pass

    @abstractmethod
    def close(self):
        pass
