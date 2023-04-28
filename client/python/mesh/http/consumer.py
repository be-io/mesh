#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from concurrent.futures import Future

from mesh.kinds import Reference
from mesh.macro import spi
from mesh.mpc import Consumer, Execution


@spi("http")
class HTTPConsumer(Consumer):

    def __enter__(self):
        pass

    def __exit__(self, exc_type, exc_val, exc_tb):
        pass

    def start(self):
        pass

    def consume(self, address: str, urn: str, execution: Execution[Reference], inbound: bytes) -> Future:
        pass
