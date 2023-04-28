#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import unittest

import mesh.log as log
import mesh.tool as tool
from mesh.kinds import Service


class TestTools(unittest.TestCase):

    def test_abstract_class(self):
        service = Service()
        service.address = '1'
        log.info(service.address)

    def test_new_trace_id(self):
        log.info(tool.new_trace_id())
        log.info(tool.get_ip())

    def test_required(self):
        log.info(f"{tool.required(True)}")
        log.info(f"{tool.required(True, False)}")
        log.info(f"{tool.required(0, 0)}")
        log.info(f"{tool.required([])}")
        log.info(f"{tool.required([], [])}")
        log.info("False begin")
        log.info(f"{tool.required()}")
        log.info(f"{tool.required(1, None)}")
        log.info(tool.get_ip())
        log.info(tool.get_ip())


if __name__ == '__main__':
    unittest.main()
