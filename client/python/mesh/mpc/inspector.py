#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Type, Any

from mesh.macro import Inspector, mpi, T


class MethodInspector(Inspector):

    def __init__(self, macro: mpi, kind: Type[T], method: Any):
        self.macro = macro
        self.kind = kind
        self.method = method

    def get_signature(self) -> str:
        return self.method.__name__

    def get_type(self) -> Type[T]:
        return self.kind

    def get_name(self) -> str:
        return self.method

    def get_annotation(self, kind: Type[T]) -> T:
        return self.method.__annotations__.get(kind)

    def get_return_type(self) -> Type[T]:
        return self.method.__annotations__.get("return")

    def invoke(self, obj: Any, args: [Any]) -> Any:
        return self.method(*args)
