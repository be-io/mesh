#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#


from typing import Any, Type

from mesh.codec import Codec, Json
from mesh.macro import index, ServiceLoader, T, serializable


@serializable
class Entity:

    @staticmethod
    def wrap(value: Any) -> "Entity":
        if value is None:
            return Entity()
        cdc = ServiceLoader.load(Codec)
        entity = Entity()
        entity.codec = cdc.default_name()
        entity.schema = ""
        entity.buffer = cdc.get_default().encode(value)
        return entity

    def present(self) -> bool:
        return self.schema is not None and self.schema != ""

    def read_object(self) -> T:
        return

    def try_read_object(self, kind: Type[T]) -> T:
        return self.load_codec().decode(self.buffer, kind) if self.buffer else None

    def load_codec(self) -> Codec:
        if not self.codec or self.codec == "":
            return ServiceLoader.load(Codec).get(Json)
        return ServiceLoader.load(Codec).get(self.codec)

    @index(0)
    def codec(self) -> str:
        return ""

    @index(5)
    def schema(self) -> str:
        return ""

    @index(10)
    def buffer(self) -> bytes:
        return bytes()


@serializable
class CacheEntity:

    @index(0)
    def version(self) -> str:
        return ""

    @index(5)
    def entity(self) -> Entity:
        return Entity()

    @index(10)
    def timestamp(self) -> int:
        return 0

    @index(15)
    def duration(self) -> int:
        return 0

    @index(20)
    def key(self) -> str:
        return ""
