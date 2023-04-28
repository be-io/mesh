#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

# https://stackoverflow.com/questions/33533148/how-do-i-type-hint-a-method-with-the-type-of-the-enclosing-class

from typing import Generic, Dict, List, Type, Any

from mesh.cause import NotFoundException
from mesh.macro.ark import T, ark
from mesh.macro.spi import spi, SPI


class Instance(Generic[T]):
    values: Dict[type, Any] = {}

    def __init__(self, name: str, kind: type, prototype=False, pattern='', priority=0) -> None:
        self.name = name
        self.kind = kind
        self.prototype = prototype
        self.pattern = pattern
        self.priority = priority

    def __lt__(self, other):
        return self.priority > other.priority

    def get_if_absent(self) -> T:
        if self.prototype:
            return self.create()

        value = self.values.get(self.kind, None)

        if value is None:
            value = self.values[self.kind] = self.create()

        return value

    def create(self) -> T:
        return self.kind()


class ServiceLoader(Generic[T]):
    loaders: Dict[Type[T], "ServiceLoader[T]"] = {}
    resources = {}

    def __init__(self, kind: type) -> None:
        self.providers: Dict[str, List[Instance]] = dict()
        self.kind = kind
        self.first = SPI.get_macro(kind).name
        self.types = {}

    @staticmethod
    def load(kind: Type[T]) -> "ServiceLoader[T]":
        if not ServiceLoader.loaders.get(kind):
            ServiceLoader.loaders[kind] = ServiceLoader(kind)
        return ServiceLoader.loaders.get(kind)

    @staticmethod
    def resource(self, name: str) -> bytes:
        if hasattr(self.providers, name):
            return self.providers[name]
        return bytes("")

    def default_name(self) -> str:
        return self.first

    def get_default(self) -> T:
        return self.get(self.first)

    def get_default_cls(self) -> Type[T]:
        instance = self.get(self.first)
        return instance.__class__

    def get(self, name: str) -> T:
        instance = self.get_instance(name)
        if not instance:
            raise NotFoundException(f'No {self.kind.__name__} named {name} exist.')
        return instance.get_if_absent()

    def list(self, pattern: str) -> List[T]:
        instances: List[Instance] = []
        for _, iis in self.get_instances().items():
            for instance in iis:
                if '' == pattern or not pattern or instance.pattern == pattern:
                    instances.append(instance)
        instances = sorted(instances, key=lambda x: x.priority)
        services = []
        for instance in instances:
            services.append(instance.get_if_absent())
        return services

    def map(self) -> Dict[str, T]:
        instances = {}
        for name, iis in self.get_instances():
            for instance in iis:
                instances[name] = instance.get_if_absent()
        return instances

    def get_instance(self, name: str) -> Instance[T]:
        instances = self.get_instances().get(name, []) if name else self.get_instances().get(self.first, [])
        return instances[0] if instances.__len__() > 0 else None

    def get_instances(self) -> Dict[str, List[Instance[T]]]:
        if self.providers.__len__() > 0:
            return self.providers
        for spec in ark.providers(spi, self.kind):
            if not isinstance(spec.metadata, SPI):
                continue
            if spec.kind is self.kind:
                self.first = spec.metadata.name
                continue
            instance = Instance(spec.metadata.name, spec.kind, spec.metadata.prototype, spec.metadata.pattern,
                                spec.metadata.priority)
            instances = self.providers.get(spec.metadata.name, [])
            instances.append(instance)
            self.providers[spec.metadata.name] = instances

        return self.providers
