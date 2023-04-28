#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import json
from typing import Type

from mesh.codec.codec import Codec, T
from mesh.macro import spi, Compatible

Json = "json"


@spi(Json)
class JsonCodec(Codec):
    native_types = (complex, int, float, bool, str, set, tuple, list, iter, bytes, dict)

    def encode(self, value: T) -> bytes:
        default = value.encode if value and hasattr(value, 'encode') else None
        return json.dumps(value, default=default).encode('UTF-8')

    def decode(self, value: bytes, kind: Type[T]) -> T:
        if value is None:
            return None

        vt = Compatible.get_args(kind)
        if vt.__len__() < 1:
            return json.loads(value, cls=kind)

        if self.is_native_type(vt[0]):
            return json.loads(value)

        return json.loads(value, cls=vt[0] if kind else None)

    def is_native_type(self, kind: type) -> bool:
        for k in self.native_types:
            if kind is k:
                return True

        return False
