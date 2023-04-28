#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from typing import Any, Generic

from mesh.macro.ark import ark, T


class MPS(Generic[T]):

    def __init__(self, fn: T = None, name='', version='', proto='grpc', codec='json', flags=0, timeout=10000):
        """
        Metadata annotation for Serial Peripheral Interface. Can be used with {@link ServiceLoader#load(Class)}
        or dependency injection at compile time and runtime time.
        """
        self.urn = ""
        self.namespace = ""
        self.name = name
        self.version = version
        self.proto = proto
        self.codec = codec
        self.flags = flags
        self.timeout = timeout
        self.retries = 10
        self.node = ""
        self.inst = ""
        self.zone = ""
        self.cluster = ""
        self.cell = ""
        self.group = ""
        self.address = ""
        self.kind = ""
        self.lang = ""
        self.attrs = {}
        if name != '':
            self.name = name
            return
        if fn is None:
            return
        if type(fn) is str:
            self.name = fn
            return
        self.name = fn.__name__
        self.kind = fn
        ark.register(mps, self.name, self.kind, self)

    def __call__(self, *args, **kwargs) -> T:
        if not hasattr(self, 'kind') or self.kind == '':
            self.kind = args[0]

        ark.register(mps, self.name, self.kind, self)
        return self.kind

    @staticmethod
    def get_mps_if_present(target: Any) -> "MPS":
        if hasattr(target, '__mps__'):
            return target.__mps__

        return MPS()


def mps(fn: T = None, *, name='', version='', proto='grpc', codec='json', flags=0, timeout=10000) -> T:
    """
    Multi Provider Service. Mesh provider service.
    :param fn:
    :param name:
    :param version:
    :param proto:
    :param codec:
    :param flags: Service flag 1 asyncable 2 encrypt 4 communal.
    :param timeout:
    :return:
    """
    return MPS(fn, name, version, proto, codec, flags, timeout)
