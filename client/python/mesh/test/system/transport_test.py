#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import os
import unittest

import mesh.log as log
from mesh.context import Mesh
from mesh.macro import ServiceLoader
from mesh.prsim import Transport, Routable, Metadata


class TestServiceLoad(unittest.TestCase):

    def test_transport_split(self):
        os.environ.setdefault("mesh.address", "127.0.0.1:17304")

        transport = Routable.of(ServiceLoader.load(Transport).get("mesh"))

        session = transport.with_address("127.0.0.1:7304").local().open('session_id_001', {
            Metadata.MESH_VERSION.key(): '',
            Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
            Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
            Metadata.MESH_TOKEN.key(): 'x',
            Metadata.MESH_SESSION_ID.key(): 'session_id_001',
            Metadata.MESH_TARGET_NODE_ID.key(): 'LX0000000000000',
        })
        session2 = transport.with_address("127.0.0.1:7304").local().open('session_id_001', {
            Metadata.MESH_VERSION.key(): '',
            Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
            Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
            Metadata.MESH_TOKEN.key(): 'x',
            Metadata.MESH_SESSION_ID.key(): 'session_id_001',
            Metadata.MESH_TARGET_NODE_ID.key(): 'LX0000000000000',
        })
        for index in range(1):
            wb = b"ABC" * 30000
            session.push(wb, {}, "1")
            rb = session2.pop(1000, "1")
            if wb.decode('utf-8') != rb.decode('utf-8'):
                log.info("PUSH-POP读写不匹配")
            else:
                log.info("PUSH-POP读写匹配")

            session.push(wb, {}, "1")
            rb = session2.peek("1")
            if wb.decode('utf-8') != rb.decode('utf-8'):
                log.info("PUSH-PEEK读写不匹配")
            else:
                log.info("PUSH-PEEK读写匹配")

        session.release(0)
        session2.release(0)


if __name__ == '__main__':
    unittest.main()
