#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import base64
import inspect
import json
from typing import Any, Type, get_type_hints

from mesh.macro.ark import T
from mesh.macro.compatible import Compatible


class Serializable(object):
    """ Serializable """

    def __init__(self, cls=None):
        self.cls = cls
        annotations = getattr(cls, '__annotations__') if hasattr(cls, '__annotations__') else {}
        for attr in inspect.classify_class_attrs(cls):
            if not attr.name or attr.name.startswith("_") or not attr.object or not hasattr(attr.object, 'kind'):
                continue
            annotations[attr.name] = getattr(attr.object, 'kind')
        setattr(cls, 'decode', self.decode)
        setattr(cls, 'encode', self.serialize)
        setattr(cls, '__annotations__', annotations)

    def serialize(self, value: Any) -> Any:
        if isinstance(value, bytes):
            return base64.b64encode(value).decode("utf-8")

        if isinstance(value, (int, float)):
            return value

        return vars(value)

    def deserialize(self, instance: Any, **kwargs) -> Any:
        if kwargs is None:
            return instance
        types = get_type_hints(instance.__class__)
        if types and types.__len__() > 0:
            for name, kind in types.items():
                value = kwargs.get(name, None)
                if value is None:
                    continue
                setattr(instance, name, self.construct(value, kind))

            return instance

        members = inspect.getmembers(instance)
        for name, _ in members:
            value = kwargs.get(name, None)
            if value is None or not name or name.startswith('_'):
                continue
            setattr(instance, name, self.construct(value, type(value)))

        return instance

    def construct(self, value: Any, vtp: type) -> Any:
        kind = Compatible.get_origin(vtp) if Compatible.get_origin(vtp) else vtp
        if kind is bool:
            return value if value else False
        if kind is complex:
            return value
        if kind is int or kind is float:
            return value if value else 0
        if kind is str:
            return value if value else ''
        if kind is list:
            sut = dict if Compatible.get_args(vtp).__len__() < 1 else Compatible.get_args(vtp)[0]
            return [self.construct(v, sut) for v in value]
        if kind is tuple:
            sut = dict if Compatible.get_args(vtp).__len__() < 1 else Compatible.get_args(vtp)[0]
            return tuple([self.construct(v, sut) for v in value])
        if kind is set:
            sut = dict if Compatible.get_args(vtp).__len__() < 1 else Compatible.get_args(vtp)[0]
            return {self.construct(v, sut) for v in value}
        if kind is iter:
            sut = dict if Compatible.get_args(vtp).__len__() < 1 else Compatible.get_args(vtp)[0]
            return iter([self.construct(v, sut) for v in value])
        if kind is dict:
            return value
        if kind is bytes:
            if isinstance(value, bytes):
                return value
            return base64.b64decode(value)
        if isinstance(value, (bool, complex, int, float, str,)):
            return value
        sub = kind()
        return self.deserialize(sub, **value)

    def decode(self, text: str):
        struct = json.loads(text)
        if struct is None:
            return None
        instance = self.cls()
        return self.deserialize(instance, **struct)


def serializable(cls: Type[T] = None) -> Type[T]:
    if cls is not None:
        Serializable(cls)
        return cls

    def classic(icls: Type[T] = None) -> Type[T]:
        Serializable(icls)
        return icls

    return classic
