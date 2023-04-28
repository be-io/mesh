#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import asyncio
from functools import wraps
from typing import Any


class CallableMember:

    def __init__(self, name: str, v: Any):
        self.name = name
        self.v = v

    def __getitem__(self, item):
        v = object.__getattribute__(self, 'v')
        if v is None:
            return None
        return v[item]

    def __call__(self, *args, **kwargs):
        return object.__getattribute__(self, 'v')

    def __getattribute__(self, name: str):
        if name == object.__getattribute__(self, 'name'):
            return self
        v = object.__getattribute__(self, 'v')
        if v is None:
            return None
        return object.__getattribute__(v, name)


class Cache(object):
    """
    A property that is only computed once per instance and then replaces itself
    with an ordinary attribute. Deleting the attribute resets the property.
    Source: https://github.com/bottlepy/bottle/commit/fa7733e075da0d790d809aa3d2f53071897e6f76
    """  # noqa

    def __init__(self, fn):
        self.__doc__ = getattr(fn, "__doc__")
        self.fn = fn
        self.name = self.fn.__name__

    def __get__(self, obj, cls):
        if obj is None:
            return self

        if asyncio and asyncio.iscoroutinefunction(self.fn):
            return self._wrap_in_coroutine(obj)

        if not hasattr(self, self.name):
            setattr(self, self.name, self.fn(obj))

        return CallableMember(self.name, getattr(self, self.name))

    def _wrap_in_coroutine(self, obj):
        @wraps(obj)
        async def wrapper():
            future = asyncio.ensure_future(self.fn(obj))
            obj.__dict__[self.fn.__name__] = future
            return future

        return wrapper()
