/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.tool.Tool;
import io.grpc.*;
import io.grpc.stub.MetadataUtils;

import javax.annotation.concurrent.ThreadSafe;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@ThreadSafe
public class GrpcInterceptor implements ServerInterceptor, ClientInterceptor {

    static final GrpcInterceptor INTERCEPTOR = new GrpcInterceptor();

    @Override
    public <O, I> ServerCall.Listener<O> interceptCall(ServerCall<O, I> invocation, Metadata metadata, ServerCallHandler<O, I> invoker) {
        Context context = Context.current()
                .withValue(GrpcContextKey.URN, Tool.anyone(metadata.get(GrpcContextKey.MESH_URN), invocation.getAuthority()))
                .withValue(GrpcContextKey.CTX_TRACE_ID, Optional.ofNullable(metadata.get(GrpcContextKey.TRACE_ID)).orElse(""))
                .withValue(GrpcContextKey.CTX_SPAN_ID, metadata.get(GrpcContextKey.SPAN_ID))
                .withValue(GrpcContextKey.CTX_FROM_INST_ID, metadata.get(GrpcContextKey.FROM_INST_ID))
                .withValue(GrpcContextKey.CTX_FROM_NODE_ID, metadata.get(GrpcContextKey.FROM_NODE_ID))
                .withValue(GrpcContextKey.CTX_INCOMING_HOST, metadata.get(GrpcContextKey.INCOMING_HOST))
                .withValue(GrpcContextKey.CTX_OUTGOING_HOST, metadata.get(GrpcContextKey.OUTGOING_HOST))
                .withValue(GrpcContextKey.CTX_INCOMING_PROXY, metadata.get(GrpcContextKey.INCOMING_PROXY))
                .withValue(GrpcContextKey.CTX_OUTGOING_PROXY, metadata.get(GrpcContextKey.OUTGOING_PROXY))
                .withValue(GrpcContextKey.CTX_SUBSET, metadata.get(GrpcContextKey.SUBSET))
                .withValue(GrpcContextKey.CTX_VERSION, metadata.get(GrpcContextKey.VERSION))
                .withValue(GrpcContextKey.CTX_TIMESTAMP, metadata.get(GrpcContextKey.TIMESTAMP))
                .withValue(GrpcContextKey.CTX_RUN_MODE, metadata.get(GrpcContextKey.RUN_MODE))
                .withValue(GrpcContextKey.CTX_TECH_PROVIDER_CODE, metadata.get(GrpcContextKey.TECH_PROVIDER_CODE))
                .withValue(GrpcContextKey.CTX_TOKEN, metadata.get(GrpcContextKey.TOKEN))
                .withValue(GrpcContextKey.CTX_TARGET_NODE_ID, metadata.get(GrpcContextKey.TARGET_NODE_ID))
                .withValue(GrpcContextKey.CTX_TARGET_INST_ID, metadata.get(GrpcContextKey.TARGET_INST_ID))
                .withValue(GrpcContextKey.CTX_SESSION_ID, metadata.get(GrpcContextKey.SESSION_ID));
        return Contexts.interceptCall(context, invocation, metadata, invoker);
    }

    @Override
    public <I, O> ClientCall<I, O> interceptCall(MethodDescriptor<I, O> descriptor, CallOptions options, Channel channel) {
        Metadata metadata = new Metadata();
        setMetadata(metadata, GrpcContextKey.TRACE_ID, Mesh.context().getTraceId());
        setMetadata(metadata, GrpcContextKey.SPAN_ID, Mesh.context().getSpanId());
        setMetadata(metadata, GrpcContextKey.FROM_INST_ID, Optional.ofNullable(Mesh.context().getConsumer().getInstId()).orElse(""));
        setMetadata(metadata, GrpcContextKey.FROM_NODE_ID, Optional.ofNullable(Mesh.context().getConsumer().getNodeId()).orElse(""));
        setMetadata(metadata, GrpcContextKey.INCOMING_HOST, String.format("%s@%s:%d", Tool.MESH_NAME.get(), Tool.MESH_RUNTIME.get().getHost(), Tool.MESH_RUNTIME.get().getPort()));
        setMetadata(metadata, GrpcContextKey.OUTGOING_HOST, com.be.mesh.client.prsim.Context.Metadata.MESH_INCOMING_HOST.get());
        setMetadata(metadata, GrpcContextKey.INCOMING_PROXY, com.be.mesh.client.prsim.Context.Metadata.MESH_INCOMING_PROXY.get());
        setMetadata(metadata, GrpcContextKey.OUTGOING_PROXY, com.be.mesh.client.prsim.Context.Metadata.MESH_OUTGOING_PROXY.get());
        setMetadata(metadata, GrpcContextKey.MESH_URN, Mesh.context().getUrn());
        setMetadata(metadata, GrpcContextKey.SUBSET, com.be.mesh.client.prsim.Context.Metadata.MESH_SUBSET.get());
        setMetadata(metadata, GrpcContextKey.VERSION, com.be.mesh.client.prsim.Context.Metadata.MESH_VERSION.get());
        setMetadata(metadata, GrpcContextKey.TIMESTAMP, com.be.mesh.client.prsim.Context.Metadata.MESH_TIMESTAMP.get());
        setMetadata(metadata, GrpcContextKey.RUN_MODE, com.be.mesh.client.prsim.Context.Metadata.MESH_RUN_MODE.get());
        // INC
        setMetadata(metadata, GrpcContextKey.TECH_PROVIDER_CODE, com.be.mesh.client.prsim.Context.Metadata.MESH_TECH_PROVIDER_CODE.get());
        setMetadata(metadata, GrpcContextKey.TOKEN, com.be.mesh.client.prsim.Context.Metadata.MESH_TOKEN.get());
        setMetadata(metadata, GrpcContextKey.TARGET_NODE_ID, com.be.mesh.client.prsim.Context.Metadata.MESH_TARGET_NODE_ID.get());
        setMetadata(metadata, GrpcContextKey.TARGET_INST_ID, com.be.mesh.client.prsim.Context.Metadata.MESH_TARGET_INST_ID.get());
        setMetadata(metadata, GrpcContextKey.SESSION_ID, com.be.mesh.client.prsim.Context.Metadata.MESH_SESSION_ID.get());
        return MetadataUtils.newAttachHeadersInterceptor(metadata).interceptCall(descriptor, options, channel);
    }

    private void setMetadata(Metadata metadata, Metadata.Key<String> key, String v) {
        if (Tool.required(v)) {
            metadata.put(key, v);
        }
    }

}
