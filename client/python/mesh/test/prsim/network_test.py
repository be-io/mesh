#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import unittest

import mesh.asm as asm
import mesh.log as log
from mesh.macro import mpi
from mesh.prsim import Network


class TestGrpc(unittest.TestCase):

    @mpi
    def network(self) -> Network:
        """"""
        pass

    def test_get_environ(self):
        asm.init()
        environ = self.network.get_environ()
        assert environ
        log.info(environ.node_id)


if __name__ == '__main__':
    unittest.main()
