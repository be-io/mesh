#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod

from mesh.kinds import Paging, Page
from mesh.macro import spi, mpi


@spi("mesh")
class IOStream(ABC):

    @abstractmethod
    @mpi(name="${mesh.name}.io.page.read", timeout=60000)
    def read(self, index: Paging) -> Page[bytes]:
        """
        Read storage with page.
        :param index:
        :return:
        """
        pass

    @abstractmethod
    @mpi(name="${mesh.name}.io.page.write", timeout=60000)
    def write(self, index: Page[bytes]):
        """
        Write storage with page.
        :param index:
        :return:
        """
        pass
