#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import abstractmethod, ABC
from asyncio import futures
from typing import Any, Dict
from typing import Type, Generic

from mesh.macro import T, Inspector, Parameters, Returns
from mesh.macro import spi
from mesh.mpc.urn import URN


@spi(name='invoker', pattern="*")
class Invoker(Generic[T], ABC):
    """
    Invoke the next invoker.
    """

    @abstractmethod
    def run(self, invocation: "Invocation") -> Any:
        pass


class Invocation(ABC):

    @abstractmethod
    def get_proxy(self) -> Invoker[T]:
        """
         Get the delegate target object.
        :return: proxy target
        """
        pass

    @abstractmethod
    def get_inspector(self) -> Inspector:
        """
        Get the invocation inspector.
        :return: inspector
        """
        pass

    @abstractmethod
    def get_parameters(self) -> Parameters:
        """
        Invoke parameters. include arguments and parameters.
        :return: parameters
        """
        pass

    @abstractmethod
    def get_arguments(self) -> [Any]:
        """
        Invoke parameters.
        :return: arguments
        """
        pass

    @abstractmethod
    def get_attachments(self) -> Dict[str, str]:
        """
        Get the attachments. The attributes will be serialized.
        :return: attachments
        """
        pass

    @abstractmethod
    def get_execution(self) -> "Execution":
        """ Get the invocation execution. """
        pass

    @abstractmethod
    def is_futures(self) -> bool:
        """
        Is the method return future.
        :return:
        """
        pass

    @abstractmethod
    def get_urn(self) -> URN:
        """
        Get the invoke urn.
        :return:
        """
        pass


class Execution(Invoker[T], Generic[T], ABC):
    """

    """

    @abstractmethod
    def schema(self) -> T:
        """ Execution schema. """
        pass

    @abstractmethod
    def inspect(self) -> Inspector:
        """ Inspect execution. """
        pass

    @abstractmethod
    def intype(self) -> Type[Parameters]:
        """ Execution input type. """
        pass

    @abstractmethod
    def retype(self) -> Type[Returns]:
        """ Execution output return type. """
        pass

    @abstractmethod
    def inflect(self) -> Parameters:
        """ Reflect input type. """
        pass

    @abstractmethod
    def reflect(self) -> Returns:
        """ Reflect output return type. """
        pass


class ServiceInvocation(Generic[T], Invocation):

    def __init__(self, proxy: Invoker[T], inspector: Inspector, parameters: Parameters, execution: Execution, urn: URN):
        self.proxy = proxy
        self.inspector = inspector
        self.parameters = parameters
        self.execution = execution
        self.urn = urn

    def get_proxy(self) -> Invoker[T]:
        return self.proxy

    def get_inspector(self) -> Inspector:
        return self.inspector

    def get_parameters(self) -> Parameters:
        return self.parameters

    def get_arguments(self) -> [Any]:
        return self.parameters.get_arguments()

    def get_attachments(self) -> Dict[str, str]:
        return self.parameters.get_attachments()

    def get_execution(self) -> Execution:
        return self.execution

    def is_futures(self) -> bool:
        rt = self.inspector.get_return_type()
        if futures.isfuture(rt):
            return True
        return False

    def get_urn(self) -> URN:
        return self.urn
