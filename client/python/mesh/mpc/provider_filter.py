#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import time
from typing import Any, Dict

import mesh.log as log
import mesh.tool as tool
from mesh.cause import MeshCode, Codeable, MeshException
from mesh.context import Mesh
from mesh.macro import spi
from mesh.mpc.digest import Digest
from mesh.mpc.filter import Filter, Invoker, Invocation, PROVIDER
from mesh.prsim import Metadata
from mesh.prsim.context import RunMode


@spi(name="provider", pattern=PROVIDER, priority=0x7fffffff)
class ProviderFilter(Filter):

    def invoke(self, invoker: Invoker, invocation: Invocation) -> Any:
        attachments: Dict[str, str] = invocation.get_parameters().get_attachments()
        if not attachments:
            attachments = {}
        Mesh.context().decode(attachments)

        digest = Digest()
        try:
            ret = invoker.run(invocation)
            digest.write("P", MeshCode.SUCCESS.get_code())
            return ret
        except BaseException as e:
            if isinstance(e, Codeable):
                digest.write("P", e.get_code())
            else:
                digest.write("P", MeshCode.SYSTEM_ERROR.get_code())
            if isinstance(e, MeshException):
                log.error(f"{digest.trace_id},{Mesh.context().get_urn()},{e.get_message}")
            raise e

    @staticmethod
    def resolve_timestamp(v: str) -> int:
        if tool.optional(v):
            return int(time.time() * 1000)
        try:
            return int(v)
        except Exception as e:
            log.error("Parse timestamp failed. ", e)
            return int(time.time() * 1000)

    @staticmethod
    def resolve_run_mode(v: str) -> RunMode:
        if tool.optional(v):
            return RunMode.ROUTINE
        try:
            return RunMode.get_by_code(int(v))
        except Exception as e:
            log.error("Parse run mode failed.", e)
            return RunMode.ROUTINE
