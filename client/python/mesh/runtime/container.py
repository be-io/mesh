#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from mesh.macro import ServiceLoader
from mesh.runtime.context import MeshContext

from mesh.ioc import PaaSProcessor, EnvsProcessor


class Container:

    def __init__(self):
        self.context = MeshContext()

    def start(self):
        self.scan()

    def stop(self):
        pass

    def wait(self):
        pass

    def scan(self):
        for processor in ServiceLoader.load(EnvsProcessor).list():
            self.context.register_processor(processor)
        for processor in ServiceLoader.load(PaaSProcessor).list():
            self.context.register_processor(processor)
