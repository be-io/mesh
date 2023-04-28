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
    os.environ.setdefault("mesh.address", "10.99.28.33:7304")
    mesh.start()

    transport = Routable.of(ServiceLoader.load(Transport).get("mesh"))
    session = transport.with_address("10.99.28.33:7304").local().open('session_id_001', {
        Metadata.MESH_VERSION.key(): '',
        Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
        Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
        Metadata.MESH_TOKEN.key(): 'x',
        Metadata.MESH_SESSION_ID.key(): 'session_id_001',
        Metadata.MESH_TARGET_INST_ID.key(): 'JG0100002700000000',
    })
    session2 = transport.with_address("10.99.28.33:7304").local().open('session_id_001', {
        Metadata.MESH_VERSION.key(): '',
        Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
        Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
        Metadata.MESH_TOKEN.key(): 'x',
        Metadata.MESH_SESSION_ID.key(): 'session_id_001',
        Metadata.MESH_TARGET_INST_ID.key(): 'JG0100002800000000',
    })
    for index in range(100):
        i2827 = f"28-27包{index}"
        log.info(f":{i2827}")
        session.push(i2827.encode('utf-8'), {})
        i2828 = f"28-28包{index}"
        log.info(f":{i2828}")
        session2.push(i2828.encode('utf-8'), {})
    session.release(0)
    session2.release(0)
    mesh.stop()


if __name__ == "__main__":
    main()
