#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any, Dict, Type

from mesh.macro import spi, T
from mesh.prsim import Dispatcher


@spi("mesh")
class MeshDispatcher(Dispatcher):

    def reference(self, mpi: Type[T]) -> T:
        pass

    def invoke(self, urn: str, param: Dict[str, Any]) -> Any:
        pass

    def invoke_generic(self, urn: str, param: Any) -> Any:
        pass
