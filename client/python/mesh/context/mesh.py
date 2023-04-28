#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import threading

__context__ = threading.local()

from collections import deque

from typing import Any, Deque, Type

from mesh.context import MeshContext, MeshKey
from mesh.prsim import Context, Key
from mesh.macro import T


class Mesh:
    # Mesh invoke timeout
    TIMEOUT: Key[int] = MeshKey("mesh.mpc.timeout", int)

    # Mesh invoke mpi name attributes.
    UNAME: Key[str] = MeshKey("mesh.mpc.uname", str)

    # Mesh mpc remote address.
    REMOTE: Key[str] = MeshKey("mesh.mpc.address", str)

    # Remote app name.
    REMOTE_NAME: Key[str] = MeshKey("mesh.mpc.remote.name", str)

    # Generic codec type
    GENERIC_CODEC_TYPE: Key[str] = MeshKey("mesh.codec.generic.type", Type[T])

    @staticmethod
    def context() -> Context:
        if not hasattr(__context__, 'mesh_context'):
            return MeshContext.create()
        mtx: Deque[Context] = __context__.mesh_context
        if not mtx or mtx.__len__() < 1:
            return MeshContext.create()
        return mtx[-1]

    @staticmethod
    def release():
        if hasattr(__context__, 'mesh_context'):
            del __context__.mesh_context

    @staticmethod
    def reset(ctx: Context):
        mtx = deque()
        mtx.append(ctx)
        __context__.mesh_context = mtx

    @staticmethod
    def is_empty() -> bool:
        return not hasattr(__context__, 'mesh_context')

    @staticmethod
    def push(ctx: Context):
        mtx: Deque[Context] = __context__.mesh_context
        if not mtx:
            return
        mtx.append(ctx)

    @staticmethod
    def pop():
        mtx: Deque[Context] = __context__.mesh_context
        if not mtx or mtx.__len__() < 1:
            return
        mtx.pop()

    @staticmethod
    def context_safe(routine) -> Any:
        """ Execute with context safety. """
        context_disable = Mesh.is_empty()
        try:
            if context_disable:
                Mesh.reset(MeshContext.create())
            else:
                Mesh.push(Mesh.context().resume())
            return routine()
        finally:
            if context_disable:
                Mesh.release()
            else:
                Mesh.pop()

    @staticmethod
    def context_id() -> str:
        return f"{threading.currentThread().ident}-{threading.currentThread().name}"
