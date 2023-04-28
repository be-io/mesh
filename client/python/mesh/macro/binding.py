#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from typing import Any, List, Generic

from mesh.macro.ark import T, ark


class Binding(Generic[T]):

    def __init__(self, fn: T = None, topic='', code='', version='', proto='grpc', codec='json', flags=0, timeout=10000,
                 meshable=True):
        """
        Metadata annotation for Serial Peripheral Interface. Can be used with {@link ServiceLoader#load(Class)}
        or dependency injection at compile time and runtime time.
        """
        self.topic = topic
        self.code = code
        self.version = version
        self.proto = proto
        self.codec = codec
        self.flags = flags
        self.timeout = timeout
        self.meshable = meshable
        if topic != '':
            self.topic = topic
            return
        if fn is None:
            return
        if type(fn) is str:
            self.name = fn
            return
        self.topic = str(type(fn))
        self.kind = fn

    def __call__(self, *args, **kwargs) -> T:
        if not hasattr(self, 'kind'):
            self.kind = args[0]

        if hasattr(self.kind, '__bindings__'):
            self.kind.__bindings__.append(self)
        else:
            self.kind.__bindings__ = [self]
        ark.register(binding, f"{self.topic}.{self.code}", self.kind, self)
        return self.kind

    @staticmethod
    def get_binding_if_present(target: Any) -> List["Binding"]:
        if hasattr(target.__class__, '__bindings__'):
            return target.__class__.__bindings__

        return []

    def matches(self, bindings: List["Binding"]) -> bool:
        if not bindings:
            return False
        for b in bindings:
            if b.topic == self.topic and b.code == self.code:
                return True
        return False


def binding(fn: T = None, *, topic='', code='', version='', proto='grpc', codec='json', flags=0, timeout=10000,
            meshable=True) -> T:
    """
    Multi Provider Service. Mesh provider service.
    :param fn:
    :param topic:
    :param code:
    :param version:
    :param proto:
    :param codec:
    :param flags:
    :param timeout:
    :param meshable:
    :return:
    """
    return Binding(fn, topic, code, version, proto, codec, flags, timeout, meshable)
