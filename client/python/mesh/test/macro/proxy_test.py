#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import mesh.log as log
import types
import unittest
from abc import ABC, abstractmethod
from typing import Any

from mesh.macro import Proxy, InvocationHandler


class Foo(ABC):

    @abstractmethod
    def bar(self) -> str:
        pass


class ProxyHandler(InvocationHandler):

    def invoke(self, proxy: Any, func: types.FunctionType, *args, **kwargs):
        return "bar"


class TestProxy(unittest.TestCase):

    def test_get_environ(self):
        foo = Proxy.new_proxy(Foo, ProxyHandler())
        bar = foo.bar()
        log.info(bar)


if __name__ == '__main__':
    unittest.main()
