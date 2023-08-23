#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import os

import mesh
from mesh import Mesh, log, ServiceLoader, Transport, Routable
from mesh.prsim import Metadata


def main():
    os.environ.setdefault("mesh.address", "127.0.0.1:17304")
    mesh.start()

    transport = Routable.of(ServiceLoader.load(Transport).get("mesh"))
    session = transport.with_address("127.0.0.1:17304").local().open('session_id_001', {
        Metadata.MESH_VERSION.key(): '',
        Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
        Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
        Metadata.MESH_TOKEN.key(): 'x',
        Metadata.MESH_SESSION_ID.key(): 'session_id_001',
        Metadata.MESH_TARGET_NODE_ID.key(): '2',
    })
    session2 = transport.with_address("127.0.0.1:27304").local().open('session_id_001', {
        Metadata.MESH_VERSION.key(): '',
        Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
        Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
        Metadata.MESH_TOKEN.key(): 'x',
        Metadata.MESH_SESSION_ID.key(): 'session_id_001',
        Metadata.MESH_TARGET_NODE_ID.key(): '2',
    })
    for index in range(100):
        b = f"{index}"
        log.info(f"写1-2: {b}")
        session.push(b.encode('utf-8'), {}, "1")
        r12 = session2.peek("1")
        log.info(f"读1-2: {r12.decode('utf-8')}")

    session.release(0)
    session2.release(0)
    mesh.stop()


if __name__ == "__main__":
    main()
