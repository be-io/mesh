#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from typing import Type, Any, Generic

from mesh.macro.ark import ark, T


class SPI(Generic[T]):

    def __init__(self, fn: T = None, name='', pattern='', priority=0, prototype=False):
        """
        Metadata annotation for Serial Peripheral Interface. Can be used with {@link ServiceLoader#load(Class)}
        or dependency injection at compile time and runtime time.
        """
        self.name = name
        self.pattern = pattern
        self.priority = priority
        self.prototype = prototype
        if name != '':
            self.name = name
            return
        if fn is None:
            return
        if type(fn) is str:
            self.name = fn
            return
        self.name = str(type(fn))
        self.kind = fn

    def __call__(self, *args, **kwargs) -> T:
        if not hasattr(self, 'kind'):
            self.kind = args[0]
        self.kind.__spi__ = self
        ark.register(spi, self.name, self.kind, self)
        return self.kind

    @staticmethod
    def get_macro(kind: Type[Any]) -> "SPI":
        if hasattr(kind, '__spi__'):
            return kind.__spi__
        return SPI()


def spi(fn: T = None, *, name='', pattern='', priority=0, prototype=False):
    return SPI(fn, name, pattern, priority, prototype)
