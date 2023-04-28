#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import mesh.log as log
import unittest

import mesh.codec as codec
import mesh.mpc as mpc


class TestGrpc(unittest.TestCase):

    def test_invoke_grpc(self):
        encoder = mpc.ServiceLoader.load(codec.Codec).get('json')
        inbound = encoder.encode("")
        execution = mpc.GenericExecution()
        self.consumer = mpc.ServiceLoader.load(mpc.Consumer).get_default()
        outbound = self.consumer.consume('https://10.12.0.83:572', execution, inbound)
        log.info(outbound)


if __name__ == '__main__':
    unittest.main()
