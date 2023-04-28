#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from mesh.codec.codec import Codec, T
from mesh.macro.spi import spi

Protobuf = "protobuf"


@spi(Protobuf)
class ProtobufCodec(Codec):

    def encode(self, value: T) -> bytes:
        pass

    def decode(self, value: bytes) -> T:
        pass
