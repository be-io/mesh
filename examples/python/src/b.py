#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import mesh
from mesh import Mesh, log, Registry, MethodProxy, ServiceLoader, Metadata, Routable, Transport
from mesh.kinds import Registration, Service, Resource


def main():
    mesh.start()

    http1 = Service()
    http1.urn = '/a/b/c'
    http1.kind = 'Restful'
    http1.namespace = ''
    http1.name = ''
    http1.version = '1.0.0'
    http1.proto = 'http'
    http1.codec = 'json'
    http1.asyncable = False
    http1.timeout = 10000
    http1.retries = 0
    http1.address = "IP:PORT"
    http1.lang = 'Python3'
    http1.attrs = {}

    grpc = Service()
    grpc.urn = '/a/b/c'
    grpc.kind = 'Restful'
    grpc.namespace = ''
    grpc.name = ''
    grpc.version = '1.0.0'
    grpc.proto = 'grpc'
    grpc.codec = 'json'
    grpc.asyncable = False
    grpc.timeout = 10000
    grpc.retries = 0
    grpc.address = "IP:PORT"
    grpc.lang = 'Python3'
    grpc.attrs = {}

    metadata = Resource()
    metadata.services = [http1, grpc]

    registration = Registration()
    registration.name = 'tensorserving'
    registration.kind = "complex"
    registration.address = "IP:PORT"
    registration.instance_id = "IP:PORT"
    registration.timestamp = 5 * 60 * 1000
    registration.attachments = {}
    registration.content = metadata

    registry = ServiceLoader.load(MethodProxy).get_default().proxy(Registry)
    Routable.of(registry).with_address("10.99.27.33:570").local().register(registration)

    transport = Routable.of(ServiceLoader.load(Transport).get("mesh"))
    session = transport.with_address("10.99.27.33:7304").local().open('session_id_001', {
        Metadata.MESH_VERSION.key(): '',
        Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
        Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
        Metadata.MESH_TOKEN.key(): 'x',
        Metadata.MESH_SESSION_ID.key(): 'session_id_001',
        Metadata.MESH_TARGET_INST_ID.key(): 'JG0100002800000000',
    })
    for _ in range(100):
        peek = session.peek()
        if peek:
            log.info(f"PEEK:{peek.decode('utf-8')}")

    for _ in range(100):
        pop = session.pop(10000)
        if pop:
            log.info(f"POP:{pop.decode('utf-8')}")

    mesh.stop()


if __name__ == "__main__":
    main()
