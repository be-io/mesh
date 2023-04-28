/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.mpc.Mesh;
import io.grpc.Context;
import io.grpc.Metadata;

/**
 * @author coyzeng@gmail.com
 */
public final class GrpcContextKey {

    private GrpcContextKey() {

    }

    public static final Context.Key<String> URN = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_URN.getKey());
    public static final Context.Key<String> CTX_MESH_URN = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_URN.getKey());
    public static final Context.Key<String> CTX_TRACE_ID = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_TRACE_ID.getKey());
    public static final Context.Key<String> CTX_SPAN_ID = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_SPAN_ID.getKey());
    public static final Context.Key<String> CTX_FROM_INST_ID = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_FROM_INST_ID.getKey());
    public static final Context.Key<String> CTX_FROM_NODE_ID = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_FROM_NODE_ID.getKey());
    public static final Context.Key<String> CTX_INCOMING_HOST = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_INCOMING_HOST.getKey());
    public static final Context.Key<String> CTX_OUTGOING_HOST = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_OUTGOING_HOST.getKey());
    public static final Context.Key<String> CTX_INCOMING_PROXY = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_INCOMING_PROXY.getKey());
    public static final Context.Key<String> CTX_OUTGOING_PROXY = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_OUTGOING_PROXY.getKey());
    public static final Context.Key<String> CTX_SUBSET = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_SUBSET.getKey());
    public static final Context.Key<String> CTX_VERSION = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_VERSION.getKey());
    public static final Context.Key<String> CTX_TIMESTAMP = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_TIMESTAMP.getKey());
    public static final Context.Key<String> CTX_RUN_MODE = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_RUN_MODE.getKey());
    // INC
    public static final Context.Key<String> CTX_TECH_PROVIDER_CODE = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_TECH_PROVIDER_CODE.getKey());
    public static final Context.Key<String> CTX_TOKEN = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_TOKEN.getKey());
    public static final Context.Key<String> CTX_TARGET_NODE_ID = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_TARGET_NODE_ID.getKey());
    public static final Context.Key<String> CTX_TARGET_INST_ID = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_TARGET_INST_ID.getKey());
    public static final Context.Key<String> CTX_SESSION_ID = Context.key(com.be.mesh.client.prsim.Context.Metadata.MESH_SESSION_ID.getKey());

    public static final Metadata.Key<String> MESH_URN = Metadata.Key.of(CTX_MESH_URN.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> TRACE_ID = Metadata.Key.of(CTX_TRACE_ID.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> SPAN_ID = Metadata.Key.of(CTX_SPAN_ID.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> FROM_INST_ID = Metadata.Key.of(CTX_FROM_INST_ID.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> FROM_NODE_ID = Metadata.Key.of(CTX_FROM_NODE_ID.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> INCOMING_HOST = Metadata.Key.of(CTX_INCOMING_HOST.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> OUTGOING_HOST = Metadata.Key.of(CTX_OUTGOING_HOST.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> INCOMING_PROXY = Metadata.Key.of(CTX_INCOMING_PROXY.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> OUTGOING_PROXY = Metadata.Key.of(CTX_OUTGOING_PROXY.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> SUBSET = Metadata.Key.of(CTX_SUBSET.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> VERSION = Metadata.Key.of(CTX_VERSION.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> TIMESTAMP = Metadata.Key.of(CTX_TIMESTAMP.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> RUN_MODE = Metadata.Key.of(CTX_RUN_MODE.toString(), Metadata.ASCII_STRING_MARSHALLER);
    // INC
    public static final Metadata.Key<String> TECH_PROVIDER_CODE = Metadata.Key.of(CTX_TECH_PROVIDER_CODE.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> TOKEN = Metadata.Key.of(CTX_TOKEN.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> TARGET_NODE_ID = Metadata.Key.of(CTX_TARGET_NODE_ID.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> TARGET_INST_ID = Metadata.Key.of(CTX_TARGET_INST_ID.toString(), Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> SESSION_ID = Metadata.Key.of(CTX_SESSION_ID.toString(), Metadata.ASCII_STRING_MARSHALLER);

    public static void setContext() {
        com.be.mesh.client.prsim.Context.Metadata.MESH_TRACE_ID.set(Mesh.context().getAttachments(), CTX_TRACE_ID.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_SPAN_ID.set(Mesh.context().getAttachments(), CTX_SPAN_ID.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_FROM_INST_ID.set(Mesh.context().getAttachments(), CTX_FROM_INST_ID.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_FROM_NODE_ID.set(Mesh.context().getAttachments(), CTX_FROM_NODE_ID.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_INCOMING_HOST.set(Mesh.context().getAttachments(), CTX_INCOMING_HOST.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_OUTGOING_HOST.set(Mesh.context().getAttachments(), CTX_OUTGOING_HOST.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_INCOMING_PROXY.set(Mesh.context().getAttachments(), CTX_INCOMING_PROXY.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_OUTGOING_PROXY.set(Mesh.context().getAttachments(), CTX_OUTGOING_PROXY.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_SUBSET.set(Mesh.context().getAttachments(), CTX_SUBSET.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_VERSION.set(Mesh.context().getAttachments(), CTX_VERSION.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_TIMESTAMP.set(Mesh.context().getAttachments(), CTX_TIMESTAMP.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_RUN_MODE.set(Mesh.context().getAttachments(), CTX_RUN_MODE.get());
        // INC
        com.be.mesh.client.prsim.Context.Metadata.MESH_TECH_PROVIDER_CODE.set(Mesh.context().getAttachments(), CTX_TECH_PROVIDER_CODE.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_TOKEN.set(Mesh.context().getAttachments(), CTX_TOKEN.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_TARGET_NODE_ID.set(Mesh.context().getAttachments(), CTX_TARGET_NODE_ID.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_TARGET_INST_ID.set(Mesh.context().getAttachments(), CTX_TARGET_INST_ID.get());
        com.be.mesh.client.prsim.Context.Metadata.MESH_SESSION_ID.set(Mesh.context().getAttachments(), CTX_SESSION_ID.get());
    }
}
