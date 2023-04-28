#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import ABC, abstractmethod
from typing import List

from mesh.kinds import License, CommerceLicense, CommerceEnviron
from mesh.macro import mpi


class Commercialize(ABC):

    @abstractmethod
    @mpi(name="mesh.license.sign", flags=2)
    def sign(self, lsr: License) -> str:
        """
        Sign the license.
        :param lsr:
        :return:
        """
        pass

    @abstractmethod
    @mpi(name="mesh.license.history", flags=2)
    def history(self, inst_id: str) -> List[CommerceLicense]:
        """
        History list the licenses.
        :param inst_id:
        :return:
        """
        pass

    @abstractmethod
    @mpi(name="mesh.net.issued", flags=2)
    def issued(self, name: str, kind: str) -> CommerceEnviron:
        """
        Issued mesh node identity.
        :param name:
        :param kind:
        :return:
        """
        pass

    @abstractmethod
    @mpi(name="mesh.net.dump", flags=2)
    def issued(self, node_id: str) -> List[CommerceEnviron]:
        """
        Dump the node identity.
        :param node_id:
        :return:
        """
        pass
