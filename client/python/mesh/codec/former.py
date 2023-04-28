#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import abstractmethod, ABC
from typing import Type, Any, get_type_hints, List

import mesh.codec.tools as tools
from mesh.macro import Index, spi, T, ServiceLoader


class Context:

    def __init__(self):
        self.offset = 0

    def form(self, value: Any, cls: Type[T]) -> T:
        """ Transform given value. """
        self.offset = self.offset + 1

    def inform(self, value: Any, cls: Type[T]) -> Any:
        self.offset = self.offset + 1

    @staticmethod
    def var_type(cls: Type[T], name: str):
        ft = get_type_hints(cls).get(name, None)
        if ft:
            return ft
        cls = getattr(cls, '__dict__', {}).get(name, None)
        return bytes if cls is Index and cls.kind is bytes else str


class ChainContext(Context):

    def __init__(self, ctx: Context):
        self.ctx = ctx

    def form(self, value: Any, cls: Type[T]) -> T:
        return super().form(value, cls)

    def inform(self, value: Any, cls: Type[T]) -> Any:
        return super().inform(value, cls)


@spi("mesh")
class Former(ABC):

    @abstractmethod
    def form(self, ctx: Context, value: Any, cls: Type[T]) -> T:
        """ Transform given value. """
        pass

    @abstractmethod
    def inform(self, ctx: Context, value: Any, cls: Type[T]) -> Any:
        """ Transform from given value. """
        pass


@spi("object")
class ObjectFormer(Former):

    def form(self, ctx: Context, value: Any, cls: Type[T]) -> T:
        if cls is not dict:
            return ctx.form(value, cls)

    def inform(self, ctx: Context, value: Any, cls: Type[T]) -> Any:
        return self.restore(value, cls)


@spi("bytes")
class BytesFormer(Former):

    def form(self, ctx: Context, value: Any, cls: Type[T]) -> T:
        if type(value) is bytes and (cls is bytes or (cls is Index and cls.kind is bytes)):
            return tools.b64encode(value).encode('utf-8')
        return ctx.form(value, cls)

    def inform(self, ctx: Context, value: Any, cls: Type[T]) -> Any:
        if type(value) is str and (cls is bytes or (cls is Index and cls.kind is bytes)):
            return tools.b64decode(value.encode('utf-8'))
        return ctx.inform(value, cls)


@spi("list")
class ListFormer(Former):

    def form(self, ctx: Context, value: Any, cls: Type[T]) -> T:
        if cls is not list:
            return ctx.form(value, cls)
        pass

    def inform(self, ctx: Context, value: Any, cls: Type[T]) -> Any:
        if cls is not list:
            return ctx.form(value, cls)


class PuppetFormer(Former):

    def __init__(self, informers: List[Former]):
        self.informers = informers

    def form(self, ctx: Context, value: Any, cls: Type[T]) -> T:
        if ctx.offset >= self.informers.__len__():
            return value
        for (index, informer) in self.informers:
            v = informer.form(ctx, value, cls)
            if index == ctx.offset:
                return v
        return value

    def inform(self, ctx: Context, value: Any, cls: Type[T]) -> Any:
        if ctx.offset >= self.informers.__len__():
            return value
        for (index, informer) in self.informers:
            v = informer.inform(ctx, value, cls)
            if index == ctx.offset:
                return v
        return value


@spi("mesh")
class MeshFormer(Former):

    def get_once(self) -> Former:
        if not hasattr(self, '__former__'):
            setattr(self, '__former__', PuppetFormer(ServiceLoader.load(Former).list('')))
        return self.__transformer__

    def matches(self, ctx: Context, value: Any, cls: Type[T]) -> bool:
        return self.get_once().matches(ctx, value, cls)

    def form(self, ctx: Context, value: Any, cls: Type[T]) -> T:
        return self.get_once().form(ctx, value, cls)

    def inform(self, ctx: Context, value: Any, cls: Type[T]) -> Any:
        return self.get_once().inform(ctx, value, cls)
