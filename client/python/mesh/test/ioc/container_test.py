#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import unittest

import runtime


class TestMacro(unittest.TestCase):

    def test_start_container(self):
        container = runtime.Container()
        container.start()

    def test_stop_container(self):
        container = runtime.Container()
        container.start()
        container.stop()


if __name__ == '__main__':
    unittest.main()
