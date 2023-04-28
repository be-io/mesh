#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import time

import mesh.log as log
import mesh.tool as tool
from mesh.context import Mesh


class Digest:

    def __init__(self):
        self.trace_id = Mesh.context().get_trace_id()
        self.span_id = Mesh.context().get_span_id()
        self.mode = Mesh.context().get_run_mode().value
        self.cni = tool.anyone(Mesh.context().get_consumer().node_id, '')
        self.cii = tool.anyone(Mesh.context().get_consumer().inst_id, '')
        self.cip = tool.anyone(Mesh.context().get_consumer().ip, '')
        self.chost = tool.anyone(Mesh.context().get_consumer().host, '')
        self.pni = tool.anyone(Mesh.context().get_provider().node_id, '')
        self.pii = tool.anyone(Mesh.context().get_provider().inst_id, '')
        self.pip = tool.anyone(Mesh.context().get_provider().ip, '')
        self.phost = tool.anyone(Mesh.context().get_provider().host, '')
        self.urn = Mesh.context().get_urn()
        self.now = int(time.time() * 1000)

    def write(self, pattern: str, code: str):
        now = int(time.time() * 1000)
        log.info(
            f"{self.trace_id},"
            f"{self.span_id},"
            f"{Mesh.context().get_timestamp()},"
            f"{self.now},"
            f"{now - Mesh.context().get_timestamp()},"
            f"{now - self.now},"
            f"{self.mode},"
            f"{pattern},"
            f"{self.cni},"
            f"{self.cii},"
            f"{self.pni},"
            f"{self.pii},"
            f"{self.cip},"
            f"{self.pip},"
            f"{self.chost},"
            f"{self.phost},"
            f"{tool.anyone(Mesh.context().get_attribute(Mesh.REMOTE), tool.get_mesh_address().any())},"
            f"{self.urn},{code}")
