#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod
from typing import Any, Dict, Type

from mesh.macro.cause import Cause
from mesh.macro.types import T


class Returns(ABC):

    @abstractmethod
    def get_code(self) -> str:
        """ Return code. """
        pass

    @abstractmethod
    def set_code(self, code: str):
        """ Return code. """
        pass

    @abstractmethod
    def get_message(self) -> str:
        """ Return message. """
        pass

    @abstractmethod
    def set_message(self, message: str):
        """ Return message. """
        pass

    @abstractmethod
    def get_cause(self) -> Cause:
        """ Return cause. """
        pass

    @abstractmethod
    def set_cause(self, cause: Cause):
        """ Return cause. """
        pass

    @abstractmethod
    def get_content(self) -> Any:
        """ Return content. """
        pass

    @abstractmethod
    def set_content(self, content: Any):
        """ Return content. """
        pass


class Parameters(ABC):

    @abstractmethod
    def map(self) -> Dict[str, Any]:
        """ Convert parameters to map. """
        pass

    @abstractmethod
    def kind(self) -> type:
        """ Parameters declared type. """
        pass

    @abstractmethod
    def get_arguments(self) -> [Any]:
        """ Generic arguments array. """
        pass

    @abstractmethod
    def set_arguments(self, *arguments):
        """ Generic arguments array. """
        pass

    @abstractmethod
    def get_attachments(self) -> Dict[str, str]:
        """ Get the generic attachments. The attributes will be serialized. The attachments are mutable. """
        pass

    @abstractmethod
    def set_attachments(self, attachments: Dict[str, str]):
        """ Attachment arguments. """
        pass


class Inspector(ABC):

    @abstractmethod
    def get_signature(self) -> str:
        """ Get the inspector signature. """
        pass

    @abstractmethod
    def get_type(self) -> Type[T]:
        """ Get the declared type of inspect object declared. """
        pass

    @abstractmethod
    def get_name(self) -> str:
        """ Get the name of inspector. """
        pass

    @abstractmethod
    def get_annotation(self, kind: Type[T]) -> T:
        """ Get the annotations of inspector. """
        pass

    @abstractmethod
    def get_return_type(self) -> Type[T]:
        """ Get the return type of inspector. """
        pass

    @abstractmethod
    def invoke(self, obj: Any, args: [Any]) -> Any:
        """ Invoke the inspector object. """
        pass
