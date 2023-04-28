#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import functools
import inspect
from typing import Any, TypeVar, Dict, List, Type

import mesh.log as log

T = TypeVar('T')
A = TypeVar('A')
B = TypeVar('B')
C = TypeVar('C')


class Species:
    macro = Any
    name = ''
    kind = Type[Any]
    metadata = Any
    trait = Type[Any]
    abstract = False


class Ark:
    def __init__(self):
        self.tinder: Dict[str, List[Species]] = dict()

    def register(self, macro: Any, name: str, kind: type, metadata: Any):
        """
        Register a annotated object as a metadata.
        :param macro: annotated decorator
        :param name: object name
        :param kind: object type
        :param metadata: decorated object
        :return: none
        """
        traits = [kind]
        if hasattr(kind, '__bases__'):
            for base in kind.__bases__:
                traits.append(base)

        for trait in traits:
            key = f'{macro.__name__}-{name}-{trait.__name__}'
            specs = self.tinder.get(key, [])
            for spec in specs:
                if spec.kind == kind:
                    log.warn(f'Object of {trait.__name__} named {name} has been register already.')

            spec = Species()
            spec.macro = macro
            spec.name = name
            spec.kind = kind
            spec.metadata = metadata
            spec.trait = trait
            spec.abstract = False
            if inspect.isabstract(kind):
                spec.abstract = True
            specs.append(spec)

            self.tinder[key] = specs

    def unregister(self, macro: Any, trait: type):
        matches = []
        for key, species in self.tinder.items():
            if species.macro == macro and species.trait == trait:
                matches.append(key)
        for _, key in matches:
            del self.tinder[key]

    def export(self, macro: Any) -> Dict[Type[T], Dict[str, List[Species]]]:
        matches: Dict[type[T], Dict[str, List[Species]]] = dict()
        for _, specs in self.tinder.items():
            for spec in specs:
                if spec.macro != macro or spec.abstract:
                    continue
                if not matches.get(spec.trait):
                    matches[spec.trait] = dict()
                if not matches.get(spec.trait).get(spec.name):
                    matches[spec.trait][spec.name] = []
                matches[spec.trait][spec.name].append(spec)

        return matches

    def providers(self, macro: Any, kind: Type[T]) -> List[Species]:
        matches = list()
        for _, specs in self.export(macro).get(kind, {}).items():
            for spec in specs:
                matches.append(spec)

        return matches

    def trait(self, macro: Any, kind: Type[T], name: str) -> List[Type[T]]:
        key = f'{macro.__name__}-{name}-{kind.__name__}'
        pvs: List[Type[T]] = []
        for spec in self.tinder.get(key, []):
            if not spec.abstract:
                pvs.append(spec.kind)

        return pvs

    @staticmethod
    def get_declared_class(method: Any):
        if isinstance(method, functools.partial):
            return Ark.get_declared_class(method.func)
        if inspect.ismethod(method) or (
                inspect.isbuiltin(method) and hasattr(method, '__self__') and hasattr(method.__self__, '__class__')):
            for cls in inspect.getmro(method.__self__.__class__):
                if method.__name__ in cls.__dict__:
                    return cls
            method = getattr(method, '__func__', method)  # fallback to __qualname__ parsing
        if inspect.isfunction(method):
            class_name = method.__qualname__.split('.<locals>', 1)[0].rsplit('.', 1)[0]
            cls = getattr(inspect.getmodule(method), class_name, None)
            if not cls:
                cls = method.__globals__.get(class_name, None)
            if isinstance(cls, type):
                return cls
        return getattr(method, '__objclass__', None)  # handle special descriptor objects


ark = Ark()
