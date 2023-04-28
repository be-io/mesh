#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import queue
from typing import Dict, Iterator, Any

import grpc

from mesh.grpx.marshaller import GrpcMarshaller
from mesh.macro import ServiceLoader
from mesh.mpc import Transporter, PROVIDER, MeshContext
from mesh.prsim import Context, Metadata


class GrpcBindableService(grpc.GenericRpcHandler):

    def __init__(self):
        self.marshaller = GrpcMarshaller()
        self.handlers = grpc.method_handlers_generic_handler("mesh-rpc", self.grpc_handlers())

    def grpc_handlers(self) -> Dict[str, grpc.RpcMethodHandler]:
        return {
            "v1": grpc.stream_stream_rpc_method_handler(
                self.stream_stream,
                request_deserializer=self.marshaller.deserialize,
                response_serializer=self.marshaller.serialize,
            ),
        }

    def service(self, handler_call_details):
        return self.handlers.service(handler_call_details)

    @staticmethod
    def stream_stream(iterator, ctx: grpc.ServicerContext) -> Iterator[Any]:
        return Transformer(iterator, ctx)


class Transformer(Iterator[Any]):

    def __init__(self, iterator: Iterator[Any], ctx: grpc.ServicerContext):
        self.iterator = iterator
        self.grpc_ctx = ctx
        self.transporter = ServiceLoader.load(Transporter).get(PROVIDER)
        self.mesh_ctx = self.context(ctx)

    def __next__(self) -> Any:
        buff = next(self.iterator)
        return self.transporter.transport(self.mesh_ctx, self.mesh_ctx.get_urn(), buff)

    def on_error(self, e: BaseException):
        """"""
        pass

    def on_complete(self):
        """"""
        pass

    @staticmethod
    def context(ctx: grpc.ServicerContext) -> Context:
        metadata = ctx.invocation_metadata()
        if metadata is None:
            return MeshContext.create()
        mtx = MeshContext()
        for (name, value) in metadata:
            std_name = name.replace("_", "-").lower()
            if std_name == Metadata.MESH_URN.key() or std_name == 'authority':
                mtx.urn = value
                mtx.attachments[Metadata.MESH_URN.key()] = value
                continue
            if std_name == Metadata.MESH_TRACE_ID.key():
                mtx.trace_id = value
                mtx.attachments[Metadata.MESH_TRACE_ID.key()] = value
                continue
            if std_name == Metadata.MESH_SPAN_ID.key():
                mtx.span_id = value
                mtx.attachments[Metadata.MESH_SPAN_ID.key()] = value
                continue
            mtx.attachments[std_name] = value
        return mtx
