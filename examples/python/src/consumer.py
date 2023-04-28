#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import abstractmethod, ABC

from mesh import mpi


class Consumer(ABC):

    @abstractmethod
    @mpi("mesh.who")
    def who(self) -> str:
        pass
