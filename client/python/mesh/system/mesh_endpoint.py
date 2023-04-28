#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.macro import mps
from mesh.prsim import Endpoint, EndpointSticker


@mps
class MeshEndpoint(Endpoint, EndpointSticker[bytes, bytes]):

    def fuzzy(self, buff: bytes) -> bytes:
        pass

    def stick(self, varg: bytes) -> bytes:
        pass
