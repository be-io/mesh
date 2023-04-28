#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import time
from typing import Any, Dict, Type

import mesh.log as log
import mesh.tool as tool
from mesh.cause import MeshCode, Codeable, MeshException
from mesh.codec import Codec
from mesh.context import Mesh, MeshContext
from mesh.kinds import MeshFlag, Location
from mesh.macro import spi, ServiceLoader
from mesh.mpc.digest import Digest
from mesh.mpc.filter import Filter, Invoker, Invocation, PROVIDER
from mesh.prsim import Metadata
from mesh.prsim.context import RunMode


@spi(name="provider", pattern=PROVIDER, priority=0x7fffffff)
class ProviderFilter(Filter):

    def invoke(self, invoker: Invoker, invocation: Invocation) -> Any:
        codec: Codec = ServiceLoader.load(Codec).get(MeshFlag.JSON.get_name())
        attachments: Dict[str, str] = invocation.get_parameters().get_attachments()
        if not attachments:
            attachments = {}
        trace_id = attachments.get(Metadata.MESH_TRACE_ID.key(), tool.new_trace_id())
        span_id = attachments.get(Metadata.MESH_SPAN_ID.key(), tool.new_span_id("", 0))
        timestamp = attachments.get(Metadata.MESH_TIMESTAMP.key(), '')
        run_mode = attachments.get(Metadata.MESH_RUN_MODE.key(), '')
        urn = attachments.get(Metadata.MESH_URN.key())
        provider = attachments.get(Metadata.MESH_PROVIDER.key(), '{}')
        pp: Location = codec.decode_string(provider, Type[Location])

        # Refresh context
        context = MeshContext()
        context.trace_id = trace_id
        context.span_id = span_id
        context.timestamp = ProviderFilter.resolve_timestamp(timestamp)
        context.run_mode = ProviderFilter.resolve_run_mode(run_mode)
        context.urn = urn
        context.consumer = pp
        for key, value in attachments.items():
            context.attachments[key] = value
        Mesh.context().rewrite_context(context)

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
            return int(time.time())
        try:
            return int(v)
        except Exception as e:
            log.error("Parse timestamp failed. ", e)
            return int(time.time())

    @staticmethod
    def resolve_run_mode(v: str) -> RunMode:
        if tool.optional(v):
            return RunMode.ROUTINE
        try:
            return RunMode.get_by_code(int(v))
        except Exception as e:
            log.error("Parse run mode failed.", e)
            return RunMode.ROUTINE
