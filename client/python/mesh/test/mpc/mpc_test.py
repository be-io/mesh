#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import unittest
from abc import ABC, abstractmethod
from typing import Any

from mesh.kinds import Principal
from mesh.macro import mpi, mps
from mesh.mpc import Mesh


class Service(ABC):

    @abstractmethod
    @mpi(name='Service.foo')
    def foo(self, hi: str) -> str:
        pass


@mps
class Implement(Service):

    def foo(self, hi: str) -> str:
        return f'I am {str}'


class TestServiceInvoke(unittest.TestCase):

    @mpi
    @property
    def service(self) -> Service:
        return self.service

    def test_foo(self):
        ret = self.service.foo("Terminator")
        self.assertEqual(ret, "I am Terminator")

    def test_context_safe(self):
        Mesh.context_safe(
            lambda: TestServiceInvoke.phase(
                lambda: TestServiceInvoke.phase(lambda: TestServiceInvoke.phase(lambda: print(-1)))))
        print(len(Mesh.context().get_principals()))

    @staticmethod
    def phase(fn: Any):
        Mesh.context().get_principals().append(Principal())
        Mesh.context_safe(fn)
        print(len(Mesh.context().get_principals()))


if __name__ == '__main__':
    unittest.main()
