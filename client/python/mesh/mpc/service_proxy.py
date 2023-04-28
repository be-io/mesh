#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import inspect
from typing import Type

from mesh.cause import ValidationException
from mesh.macro import spi, mpi, T, Proxy, MethodProxy
from mesh.mpc.reference import ReferenceInvokeHandler


@spi("mpc")
class ServiceProxy(MethodProxy):

    def proxy(self, kind: Type[T]) -> T:
        return ServiceProxy.default_proxy(kind)

    @staticmethod
    def default_proxy(kind: Type[T]) -> T:
        return ServiceProxy.static_proxy(kind, mpi())

    @staticmethod
    def static_proxy(kind: Type[T], metadata: mpi) -> T:
        interfaces = [kind]
        if hasattr(kind, "__args__") and getattr(kind, "__args__") is not None:
            interfaces = kind.__args__

        for interface in interfaces:
            if not inspect.isabstract(interface):
                raise ValidationException(f'{str(interface)} is not abstract class. ')

        if metadata is None:
            metadata = mpi()

        return Proxy(interfaces, ReferenceInvokeHandler(metadata, interfaces))
