/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.prsim.Codeable;
import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.struct.Digest;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI(value = "consumer", pattern = Filter.CONSUMER, priority = Integer.MIN_VALUE)
public class ConsumerFilter implements Filter {

    @Override
    public Object invoke(Invoker<?> invoker, Invocation invocation) throws Throwable {
        Codec codec = ServiceLoader.load(Codec.class).get(MeshFlag.JSON.getName());
        Map<String, String> attachments = invocation.getParameters().attachments();
        attachments.putAll(Mesh.context().getAttachments());
        Context.Metadata.MESH_TRACE_ID.set(attachments, Mesh.context().getTraceId());
        Context.Metadata.MESH_SPAN_ID.set(attachments, Mesh.context().getSpanId());
        Context.Metadata.MESH_TIMESTAMP.set(attachments, String.valueOf(Mesh.context().getTimestamp()));
        Context.Metadata.MESH_RUN_MODE.set(attachments, String.valueOf(Mesh.context().getRunMode().getMode()));
        Context.Metadata.MESH_URN.set(attachments, String.valueOf(Mesh.context().getUrn()));
        Context.Metadata.MESH_CONSUMER.set(attachments, codec.encodeString(Mesh.context().getConsumer()));
        Context.Metadata.MESH_PROVIDER.set(attachments, codec.encodeString(Mesh.context().getProvider()));

        Digest digest = new Digest();

        try {
            // Refresh slf4j MDC
            digest.mdc();
            Object ret = invoker.invoke(invocation);
            digest.print("C", MeshCode.SUCCESS.getCode());
            return ret;
        } catch (Throwable e) {
            Cause.POS.ifPresent(pos -> log.error("MPC invoke fault at {} with {}", pos, digest.getTraceId()));
            if (e instanceof Codeable) {
                digest.print("C", ((Codeable) e).getCode());
            } else {
                digest.print("C", MeshCode.SYSTEM_ERROR.getCode());
            }
            if (e instanceof MeshException) {
                log.error("{},{},{}", digest.getTraceId(), Mesh.context().getUrn(), ((MeshException) e).getRootMessage());
            }
            throw e;
        }
    }

}
