#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import collections

import grpc

import mesh.tool as tool
from mesh.mpc import Mesh
from mesh.prsim import Metadata


class ClientCallDetails(
    collections.namedtuple('_ClientCallDetails',
                           ('method', 'timeout', 'metadata', 'credentials',
                            'wait_for_ready', 'compression')), grpc.ClientCallDetails):
    pass


class GrpcInterceptor(grpc.ServerInterceptor,
                      grpc.StreamStreamClientInterceptor,
                      grpc.StreamUnaryClientInterceptor,
                      grpc.UnaryStreamClientInterceptor,
                      grpc.UnaryUnaryClientInterceptor, ):
    TKeys = [
        Metadata.MESH_INCOMING_PROXY,
        Metadata.MESH_OUTGOING_PROXY,
        Metadata.MESH_SUBSET,
        Metadata.MESH_VERSION,
        Metadata.MESH_TIMESTAMP,
        Metadata.MESH_RUN_MODE,
        # INC
        Metadata.MESH_TECH_PROVIDER_CODE,
        Metadata.MESH_TOKEN,
        Metadata.MESH_TARGET_NODE_ID,
        Metadata.MESH_TARGET_INST_ID,
        Metadata.MESH_SESSION_ID
    ]

    @staticmethod
    def client_context(context):
        # python metadata must be lowercase
        attachments = Mesh.context().get_attachments()
        metadata = []
        Metadata.MESH_URN.append(metadata, Mesh.context().get_urn())
        Metadata.MESH_TRACE_ID.append(metadata, Mesh.context().get_trace_id())
        Metadata.MESH_SPAN_ID.append(metadata, Mesh.context().get_span_id())
        Metadata.MESH_FROM_INST_ID.append(metadata, Mesh.context().get_consumer().inst_id)
        Metadata.MESH_FROM_NODE_ID.append(metadata, Mesh.context().get_consumer().node_id)
        Metadata.MESH_INCOMING_HOST.append(metadata, f"{tool.get_mesh_name()}@{str(tool.get_mesh_runtime())}")
        Metadata.MESH_OUTGOING_HOST.append(metadata, attachments.get(Metadata.MESH_INCOMING_HOST.key(), ''))

        for mk in GrpcInterceptor.TKeys:
            mk.append(metadata, attachments.get(mk.key(), ''))

        if context.metadata is not None:
            for x in context.metadata:
                metadata.append(x)

        wait_for_ready = context.wait_for_ready if hasattr(context, 'wait_for_ready') else None
        compression = context.compression if hasattr(context, 'compression') else None
        credentials = context.credentials if hasattr(context, 'credentials') else None
        return ClientCallDetails(context.method, context.timeout, metadata, credentials, wait_for_ready, compression)

    def intercept_stream_stream(self, continuation, client_call_details, request_iterator):
        return continuation(self.client_context(client_call_details), request_iterator)

    def intercept_stream_unary(self, continuation, client_call_details, request_iterator):
        return continuation(self.client_context(client_call_details), request_iterator)

    def intercept_unary_stream(self, continuation, client_call_details, request):
        return continuation(self.client_context(client_call_details), request)

    def intercept_unary_unary(self, continuation, client_call_details, request):
        return continuation(self.client_context(client_call_details), request)

    def intercept_service(self, continuation, handler_call_details):
        self.server_context(handler_call_details)
        return continuation(handler_call_details)

    @staticmethod
    def server_context(handler_call_details):
        pass
