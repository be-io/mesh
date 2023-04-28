#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#


import inspect
from abc import abstractmethod
from typing import Any, Generic, Type, List, Dict, Callable

from mesh.macro.ark import T


class InvocationHandler:
    """
    InvocationHandler is the interface implemented by the invocation handler of a proxy instance.
    Each proxy instance has an associated invocation handler. When a method is invoked on a proxy instance,
    the method invocation is encoded and dispatched to the invoke method of its invocation handler.
    """

    @abstractmethod
    def invoke(self, proxy: Any, method: Any, *args, **kwargs):
        pass


class InvocationException(Exception):
    """
    Invocation exceptions.
    """

    def __init__(self, cls):
        super(InvocationException, self).__init__(cls, 'is not a InvocationHandler')


class Proxy(Generic[T]):
    """
    Abstract class dynamic proxy
    https://code.activestate.com/recipes/496741-object-proxying/

    Proxy provides static methods for creating objects that act like instances of interfaces but allow for customized
    method invocation. To create a proxy instance for some interface Foo:

    handler = new MyInvocationHandler(...)
    foo = (Foo) Proxy(cls, handler)

    A proxy class is a class created at runtime that implements a specified list of interfaces, known as
    proxy interfaces. A proxy instance is an instance of a proxy class. Each proxy instance has an
    associated invocation handler object, which implements the interface InvocationHandler. A method invocation
    on a proxy instance through one of its proxy interfaces will be dispatched to the invoke method of the instance's
    invocation handler, passing the proxy instance, a java.lang.reflect.Method object identifying the method that was
    invoked, and an array of type Object containing the arguments. The invocation handler processes the encoded method
    invocation as appropriate and the result that it returns will be returned as the result of the method invocation
    on the proxy instance.
    """

    __slots__ = ["_target_", "__weakref__"]

    def __target__(self):
        return object.__getattribute__(self, "_target_")

    def __handler__(self):
        return object.__getattribute__(self, "_handler_")

    def __nonzero__(self):
        return bool(self.__target__())

    def __init__(self, target: Any, handler: InvocationHandler):
        object.__setattr__(self, "_target_", target)
        object.__setattr__(self, "_handler_", handler)

    def __getattribute__(self, name):
        if object.__getattribute__(self, name):
            return object.__getattribute__(self, name)
        if self.__target__():
            return object.__getattribute__(self.__target__(), name)
        return None

    def __delattr__(self, name):
        delattr(self.__target__(), name)

    def __setattr__(self, name, value):
        setattr(self.__target__(), name, value)

    def __str__(self):
        return str(self.__target__())

    def __repr__(self):
        return repr(self.__target__())

    def __hash__(self):
        return hash(self.__target__())

    def __call__(self, *args, **kwargs):
        return self

    def __new__(cls, kinds: List[Type], handler: InvocationHandler, *args, **kwargs):
        """
        Creates an proxy instance referencing `interface`. (interface, *args, **kwargs) are
        passed to this class' __init__, so deriving classes can define an
        __init__ method of their own.
        note: _proxy_class_cache_ is unique per deriving class (each deriving
        class must hold its own cache)
        """
        proxy_class = cls.create_proxy_class(kinds)
        instance = object.__new__(proxy_class)
        instance.__init__(instance, handler)
        return instance

    _special_names_ = [
        '__abs__', '__add__', '__and__', '__cmp__', '__coerce__',
        '__contains__', '__delitem__', '__delslice__', '__div__', '__divmod__',
        '__eq__', '__float__', '__floordiv__', '__ge__', '__getitem__',
        '__getslice__', '__gt__', '__hex__', '__iadd__', '__iand__',
        '__idiv__', '__idivmod__', '__ifloordiv__', '__ilshift__', '__imod__',
        '__imul__', '__int__', '__invert__', '__ior__', '__ipow__', '__irshift__',
        '__isub__', '__iter__', '__itruediv__', '__ixor__', '__le__', '__len__',
        '__long__', '__lshift__', '__lt__', '__mod__', '__mul__', '__ne__',
        '__neg__', '__oct__', '__or__', '__pos__', '__pow__', '__radd__',
        '__rand__', '__rdiv__', '__rdivmod__', '__reduce__', '__reduce_ex__',
        '__repr__', '__reversed__', '__rfloorfiv__', '__rlshift__', '__rmod__',
        '__rmul__', '__ror__', '__rpow__', '__rrshift__', '__rshift__', '__rsub__',
        '__rtruediv__', '__rxor__', '__setitem__', '__setslice__', '__sub__',
        '__truediv__', '__xor__', 'next',
    ]

    @classmethod
    def create_proxy_class(cls, kinds: List[Type]) -> Type["Proxy"]:
        """
        Creates a proxy for the given class
        """
        qname = ",".join(sorted(map(lambda x: str(x), filter(lambda k: k != Proxy, kinds))))
        if hasattr(cls, '__proxy_classes__') and cls.__proxy_classes__.get(qname):
            return cls.__proxy_classes__.get(qname)

        if not hasattr(cls, '__proxy_classes__'):
            cls.__proxy_classes__ = {}

        namespace = {}
        for kind in kinds:
            for name in cls._special_names_:
                if hasattr(kind, name) and not hasattr(cls, name) and not namespace.get(name):
                    namespace[name] = Proxy.make_special(name)

        for kind in kinds:
            for name, method in Proxy.get_abstract_methods(kind).items():
                namespace[name] = Proxy.make_delegate(method)

        cls.__proxy_classes__[qname] = type("%s(%s)" % (cls.__name__, qname), (*kinds, cls), namespace)
        return cls.__proxy_classes__[qname]

    @staticmethod
    def new_proxy(kind: Type[T], handler: InvocationHandler) -> T:
        interface = [kind]
        if hasattr(kind, "__args__"):
            interface = kind.__args__

        return Proxy(interface, handler)

    @staticmethod
    def get_interfaces(kind: Type) -> List[Type]:
        if inspect.isabstract(kind):
            return [kind]
        interfaces = []
        if hasattr(kind, '__bases__'):
            for base in kind.__bases__:
                interfaces.append(base)
        if interfaces.__len__() < 1:
            return [kind]
        return interfaces

    @staticmethod
    def get_abstract_methods(kind: Type) -> Dict[str, Callable]:
        interfaces = Proxy.get_interfaces(kind)
        methods: Dict[str, Callable] = {}
        for interface in interfaces:
            if hasattr(interface, '__abstractmethods__'):
                for name in interface.__abstractmethods__:
                    if hasattr(interface, name):
                        methods[name] = interface.__dict__[name]

        return methods

    @staticmethod
    def make_special(method: str):
        def invoke(self, *args, **kw):
            return getattr(self.__target__(), method)(*args, **kw)

        return invoke

    @staticmethod
    def make_delegate(method):
        def invoke(self, *args, **kwargs):
            return self.__handler__().invoke(self, method, *args, **kwargs)

        return invoke
