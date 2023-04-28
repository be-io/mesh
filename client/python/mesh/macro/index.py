#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import inspect
from typing import Any, Callable, Generic

from mesh.macro.ark import T


class Index(Generic[T], object):

    def __init__(self, fn=-1, idx=-1, name='', dft='', doc=None, kind=Any, value=None, initializer=None):
        """
        Emulate PyProperty_Type() in Objects/descrobject.c
        Index for protobuf or thrift etc.
        :param fn: decorate function
        :param idx: Index position.
        """
        self.idx = idx
        self.name = name
        self.dft = dft
        self.doc = doc
        self.kind = kind
        self.value = value
        self.initializer = initializer
        if self.idx != -1:
            return
        if fn == -1:
            return
        if type(fn) is int:
            self.idx = fn

    def __call__(self, *args, **kwargs) -> Callable[..., T]:
        self.initializer = args[0]
        if hasattr(self.initializer, "__name__"):
            self.name = self.initializer.__name__
        if hasattr(self.initializer, "__doc__"):
            self.doc = self.initializer.__doc__
        if hasattr(self.initializer, "__annotations__"):
            signature = inspect.signature(self.initializer)
            self.kind = signature.return_annotation

        return self

    def __get__(self, ref, kind) -> T:
        if ref is None:
            return None
        value = ref.__dict__.get(self.name, None)
        if value is None and callable(self.initializer):
            value = self.initializer(ref)
            ref.__dict__[self.name] = value
        return value

    def __set__(self, ref, value: T):
        if ref is None:
            return
        ref.__dict__[self.name] = value

    def __delete__(self, ref):
        if ref is None:
            return
        if ref.__dict__.get(self.name):
            del ref.__dict__[self.name]


def index(fn=-1, *, idx=-1, name='', dft=None) -> Callable[[Callable[..., T]], T]:
    """
    Index for protobuf or thrift/protobuf etc.
    """
    return Index[T](fn, idx, name, dft)
