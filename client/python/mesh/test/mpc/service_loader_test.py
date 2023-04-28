#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import unittest

from mesh.codec import Codec
from mesh.macro import ServiceLoader


class TestServiceLoad(unittest.TestCase):

    def test_load(self):
        codec = ServiceLoader.load(Codec).load()


if __name__ == '__main__':
    unittest.main()
