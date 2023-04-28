#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import mesh.log as log
import typing
import unittest

from mesh.macro import ServiceLoader
from mesh.mpc import Compiler


class TestCompile(unittest.TestCase):

    def test_compile_code(self):
        compiler = ServiceLoader.load(Compiler).get_default()
        proxy = compiler.compile("""
from mesh.macro import index

class Terminator:

    baz: str

    def get_name(self)->str:
        return "I am back!"
    
    @index(0)
    def get_index(self)->int:
        return 9
""", "proxy", 1)
        terminator = proxy.Terminator()
        name = terminator.get_name()
        log.info(name)
        members = typing.get_type_hints(proxy.Terminator)
        log.info(members)


if __name__ == '__main__':
    unittest.main()
