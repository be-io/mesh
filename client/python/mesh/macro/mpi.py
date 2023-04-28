#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#


import inspect
from abc import ABC, abstractmethod
from typing import Any, Type

from mesh.macro.ark import ark, T
from mesh.macro.compatible import Compatible
from mesh.macro.loader import ServiceLoader
from mesh.macro.spi import spi


@spi("mpc")
class MethodProxy(ABC):

    @abstractmethod
    def proxy(self, kind: Type[T]) -> T:
        """Proxy facade"""
        pass


class MPI:

    def __init__(self, fn: str = None, name='', version='', proto='grpc', codec='json', flags=0, timeout=10000,
                 retries=3, node='', inst='', zone='', cluster='', cell='', group='', address=''):
        """
        Multi Provider Interface. Mesh provider interface.

        :param name: service name.
        :param version: service version.
        :param proto: invoke net protocol.
        :param codec: invoke serialize protocol.
        :param flags: service flag 1 asyncable 2 encrypt 4 communal.
        :param timeout: invoke timeout with mills.
        :param retries: invoke retry times.
        :param node: invoke to any network node.
        :param inst: alias of node.
        :param zone: invoke to network zone.
        :param cluster: invoke to network cluster.
        :param cell: invoke to network cell.
        :param group: invoke to network group.
        :param address: invoke to direct address.
        """
        self.__proxy__ = None
        self.urn = ""
        self.namespace = ""
        self.name = name
        self.version = version
        self.proto = proto
        self.codec = codec
        self.flags = flags
        self.timeout = timeout
        self.retries = retries
        self.node = node
        self.inst = inst
        self.zone = zone
        self.cluster = cluster
        self.cell = cell
        self.group = group
        self.address = address
        if fn:
            fn_is_str = type(fn) is str
            fn_qualname = getattr(fn, '__qualname__', '').split(".")
            fn_name = fn_qualname[0] if fn_qualname.__len__() < 2 else fn_qualname[1]
            fn_namespace = '' if fn_qualname.__len__() < 2 else fn_qualname[0]
            self.fn = None if fn_is_str else fn
            self.declared_kind = None if fn_is_str else ark.get_declared_class(self.fn)
            self.namespace = '' if fn_is_str else fn_namespace
            self.name = fn if fn_is_str else fn_name

    def __call__(self, *args, **kwargs):
        """ Interpreter is executing when called, so some metadata must get delay """
        if not getattr(self, 'fn', None):
            self.fn = args[0]
            self.declared_kind = ark.get_declared_class(self.fn)
        if self.declared_kind or inspect.isabstract(self.declared_kind):
            signature = inspect.signature(self.fn)
            self.__proxy__ = ServiceLoader.load(MethodProxy).get_default().proxy(signature.return_annotation)

        ark.register(mpi, self.name, self.fn, self)
        self.fn.__mpi__ = self
        return self.fn

    def __get__(self, ref, kind):
        return self.proxy_if_absent()

    def proxy_if_absent(self):
        if not self.declared_kind or not inspect.isabstract(self.declared_kind):
            from mesh.prsim import Routable
            self.declared_kind = ark.get_declared_class(self.fn)
            signature = inspect.signature(self.fn)
            origin = Compatible.get_origin(signature.return_annotation)
            if origin == Routable:
                parameters = Compatible.get_args(signature.return_annotation)
                self.__proxy__ = Routable.of(ServiceLoader.load(MethodProxy).get_default().proxy(parameters[0]))
            else:
                self.__proxy__ = ServiceLoader.load(MethodProxy).get_default().proxy(signature.return_annotation)

        return getattr(self, '__proxy__', self.fn)

    @staticmethod
    def get_mpi_if_present(target: Any) -> "MPI":
        if hasattr(target, '__mpi__'):
            return target.__mpi__

        return MPI()


def mpi(fn: str = None, *, name='', version='', proto='grpc', codec='json', flags=0,
        timeout=10000, retries=3, node='', inst='', zone='', cluster='', cell='', group='', address='') -> T:
    return MPI(fn, name, version, proto, codec, flags, timeout, retries, node, inst, zone, cluster, cell, group,
               address)
