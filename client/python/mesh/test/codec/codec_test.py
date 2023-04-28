#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import unittest

import mesh.codec as codec
import mesh.codec.tools as tools
import mesh.log as log
import mesh.mpc as mpc
from mesh.kinds import Paging, Page


class A:
    x: str


class TestCodec(unittest.TestCase):

    def test_json_codec(self):
        cls = Page[bytes]
        log.info(tools.get_raw_type(cls))
        paging = Paging()
        paging.index = 1
        encoder = mpc.ServiceLoader.load(codec.Codec).get(codec.Json)
        inbound = encoder.encode_string(paging)
        log.info(inbound)
        log.info(encoder.encode_string({'x': 'y'}))
        a = A()
        a.x = 'y'
        log.info(encoder.encode_string(a))

        encoder.decode_string('{}', Page[bytes])


if __name__ == '__main__':
    unittest.main()
