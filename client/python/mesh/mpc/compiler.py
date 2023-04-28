#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import inspect
from abc import abstractmethod, ABC
from typing import Any, Type, Dict, List

import mesh.tool as tool
from mesh.macro import Parameters, Returns, spi, Index, Cause, serializable

ATTACHMENTS: str = "attachments"
CODE: str = "code"
MESSAGE: str = "message"
CONTENT: str = "content"
CAUSE: str = "cause"


@spi(name="mpyc")
class Compiler(ABC):

    @abstractmethod
    def intype(self, method: Any) -> Type[Parameters]:
        """
        Get the parameter class type.
        :param method:
        :return:
        """
        pass

    @abstractmethod
    def retype(self, method: Any) -> Type[Returns]:
        """
        Get the return class type.
        :param method:
        :return:
        """
        pass

    @abstractmethod
    def compile(self, code: str, name: str, add_to_sys_modules=0):
        """
        Compile source code.
        :param code:
        :param name:
        :param add_to_sys_modules:
        :return:
        """
        pass


@spi(name="mpyc")
class PythonCompiler(Compiler):
    cache = {}

    @staticmethod
    def make_getter(name, kind: type):
        def getter(self) -> kind:
            if hasattr(self, name):
                return getattr(self, name)
            return None

        return getter

    @staticmethod
    def make_setter(name):
        def setter(self, value):
            setattr(self, name, value)

        return setter

    @staticmethod
    def make_kind():
        def kind(self) -> type:
            return self.__cls__

        return kind

    @staticmethod
    def make_map():
        def maps(self, names: List[str]) -> Dict[str, Any]:
            values = {}
            for name in names:
                if hasattr(self, name):
                    values[name] = getattr(self, name)
            return values

        return maps

    @staticmethod
    def make_get_arguments(names: List[str]):
        def getter(self) -> List[Any]:
            values = []
            for name in names:
                if hasattr(self, name):
                    values.append(getattr(self, name))
                else:
                    values.append(None)
            return values

        return getter

    @staticmethod
    def make_set_arguments(names: List[str]):
        def setter(self, *arguments):
            if not arguments:
                return
            for idx, name in enumerate(names):
                if idx < arguments.__len__():
                    setattr(self, name, arguments[idx])

        return setter

    @staticmethod
    def make_decode(names: List[str]):
        def setter(self, values: List[Any]):
            if not values:
                return
            for idx, name in enumerate(names):
                if idx < values.__len__():
                    setattr(self, name, values[idx])

        return setter

    @staticmethod
    def make_class_name(method: Any, suffix: str) -> str:
        qualname = tool.split(method.__qualname__, ".")
        return f"Mesh{qualname[0]}{qualname[1].title()}{suffix}".replace("_", "")

    def intype(self, method: Any) -> Type[Parameters]:
        signature = inspect.signature(method)
        class_name = self.make_class_name(method, "Parameters")
        intype = self.cache.get(class_name, None)
        if intype:
            return intype
        names = []
        annotations = {ATTACHMENTS: Dict[str, str]}
        for idx, (name, parameter) in enumerate(signature.parameters.items()):
            if idx > 0:
                names.append(name)
                annotations[name] = parameter.annotation
        variables = {
            "get_attachments": self.make_getter(ATTACHMENTS, Dict[str, str]),
            "set_attachments": self.make_setter(ATTACHMENTS),
            "map": self.make_map(),
            "kind": self.make_kind(),
            "get_arguments": self.make_get_arguments(names),
            "set_arguments": self.make_set_arguments(names),
            "__annotations__": annotations,
        }
        for idx, (name, parameter) in enumerate(signature.parameters.items()):
            if idx == 0:
                variables[ATTACHMENTS] = Index(idx=-1, name=ATTACHMENTS, kind=Dict)
                continue
            variables[name] = Index(idx=idx, name=name, kind=parameter.annotation)
        intype = serializable(type(class_name, (Parameters,), variables))
        self.cache.__setitem__(class_name, intype)
        return intype

    def retype(self, method: Any) -> Type[Returns]:
        signature = inspect.signature(method)
        class_name = self.make_class_name(method, "Returns")
        retype = self.cache.get(class_name, None)
        if retype:
            return retype
        variables = {
            CODE: Index(idx=0, name=CODE, kind=str),
            MESSAGE: Index(idx=5, name=MESSAGE, kind=str),
            CAUSE: Index(idx=10, name=CAUSE, kind=Cause),
            CONTENT: Index(idx=15, name=CONTENT, kind=signature.return_annotation),
            "get_code": self.make_getter(CODE, str),
            "set_code": self.make_setter(CODE),
            "get_message": self.make_getter(MESSAGE, str),
            "set_message": self.make_setter(MESSAGE),
            "get_cause": self.make_getter(CAUSE, Cause),
            "set_cause": self.make_setter(CAUSE),
            "get_content": self.make_getter(CONTENT, signature.return_annotation),
            "set_content": self.make_setter(CONTENT),
            "__annotations__": {
                CODE: str,
                MESSAGE: str,
                CAUSE: Cause,
                CONTENT: signature.return_annotation
            },
        }
        retype = serializable(type(class_name, (Returns,), variables))
        self.cache.__setitem__(class_name, retype)
        return retype

    def compile(self, code: str, name: str, add_to_sys_modules=0):
        """
        Import dynamically generated code as a module. code is the
        object containing the code (a string, a file handle or an
        actual compiled code object, same types as accepted by an
        exec statement). The name is the name to give to the module,
        and the final argument says whether to add it to sys.modules
        or not. If it is added, a subsequent import statement using
        name will return this module. If it is not added to sys.modules
        import will try to load it in the normal fashion.
        """
        import sys
        import types
        module = types.ModuleType(name)
        exec(code, module.__dict__)
        if add_to_sys_modules:
            sys.modules[name] = module
        return module
