#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#


from abc import ABC, abstractmethod
from typing import Dict, Generic, List

from mesh.kinds import Principal
from mesh.macro import T, ServiceLoader, spi


@spi("mesh")
class Routable(ABC, Generic[T]):

    @abstractmethod
    def __init__(self, reference: T = None, attachments: Dict[str, str] = None, address: str = ""):
        pass

    def __call__(self, *args, **kwargs):
        return self

    @abstractmethod
    def within(self, key: str, value: str) -> "Routable[T]":
        """
        Route with attachments.
        :param key: attachment key
        :param value: attachment value
        :return:
        """
        pass

    @abstractmethod
    def with_map(self, attachments: Dict[str, str]) -> "Routable[T]":
        """
        Invoke the service in local network.
        :param attachments: attachments
        :return:
        """
        pass

    @abstractmethod
    def with_address(self, address: str) -> "Routable[T]":
        """
        Invoke the service in many network, it may be local or others. Broadcast mode.
        :param address: Network address.
        :return: Service invoker.
        """
        pass

    @abstractmethod
    def local(self) -> T:
        """
        Invoke the service in local network.
        :return:
        """
        pass

    @abstractmethod
    def any(self, principal: Principal) -> T:
        """
        Invoke the service in a network, it may be local or others.
        :param principal: Network principal.
        :return:Service invoker.
        """
        pass

    @abstractmethod
    def any_inst(self, inst_id: str) -> T:
        """
        Invoke the service in a network, it may be local or others.
        :param inst_id: Network principal of inst_id.
        :return:Service invoker.
        """
        pass

    @abstractmethod
    def many(self, principals: List[Principal]) -> List[T]:
        """
        Invoke the service in many network, it may be local or others. Broadcast mode.
        :param principals: Network principals.
        :return: Service invoker.
        """
        pass

    @abstractmethod
    def many_inst(self, inst_ids: List[str]) -> List[T]:
        """
        Invoke the service in many network, it may be local or others. Broadcast mode.
        :param inst_ids: Network principals.
        :return: Service invoker.
        """
        pass

    @staticmethod
    def of(reference: T) -> "Routable[T]":
        """
        Wrap a service with streamable ability.
        :param reference:  Service reference.
        :return: Streamable program interface.
        """
        cls = ServiceLoader.load(Routable).get_default_cls()
        return cls(reference, {}, "")
