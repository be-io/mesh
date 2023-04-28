#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import mesh.tool as tool
from mesh.macro import ServiceLoader
from mesh.mpc import Provider
from mesh.prsim import RuntimeHook


class Mooter(RuntimeHook):
    DEFAULT_PORT = 8864

    def __init__(self):
        self.http1: Provider = None
        self.http2: Provider = None

    def start(self):
        hooks = ServiceLoader.load(RuntimeHook).list('')
        for hook in hooks:
            hook.start()

        self.http1: Provider = ServiceLoader.load(Provider).get("http")
        self.http2: Provider = ServiceLoader.load(Provider).get("grpc")
        self.http1.start(f"0.0.0.0:{tool.get_mesh_runtime().port + 1}", None)
        self.http2.start(f"0.0.0.0:{tool.get_mesh_runtime().port}", None)

    def stop(self):
        if self.http1 is not None:
            self.http1.close()
        if self.http2 is not None:
            self.http2.close()
        hooks = ServiceLoader.load(RuntimeHook).list('')
        for hook in hooks:
            hook.stop()

    def refresh(self):
        self.stop()
        self.start()
