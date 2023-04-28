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
import com.be.mesh.client.struct.Digest;
import com.be.mesh.client.struct.Location;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.MDC;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j(topic = "digest")
@SPI(value = "provider", pattern = Filter.PROVIDER, priority = Integer.MAX_VALUE)
public class ProviderFilter implements Filter {

    @Override
    public Object invoke(Invoker<?> invoker, Invocation invocation) throws Throwable {
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        Map<String, String> attachments = Optional.ofNullable(invocation.getParameters().attachments()).orElseGet(HashMap::new);
        String traceId = Context.Metadata.MESH_TRACE_ID.get(attachments, Tool::newTraceId);
        String spanId = Context.Metadata.MESH_SPAN_ID.get(attachments, () -> Tool.newSpanId("", 0));
        String timestamp = Context.Metadata.MESH_TIMESTAMP.get(attachments);
        String runMode = Context.Metadata.MESH_RUN_MODE.get(attachments);
        String urn = Context.Metadata.MESH_URN.get(attachments);
        String provider = Context.Metadata.MESH_PROVIDER.get(attachments);
        Location pp = codec.decodeString(provider, Types.of(Location.class));

        // Refresh context
        MpcContext context = new MpcContext();
        context.setTraceId(traceId);
        context.setSpanId(spanId);
        context.setTimestamp(resolveTimestamp(timestamp));
        context.setRunMode(resolveRunMode(runMode));
        context.setUrn(urn);
        context.setConsumer(pp);
        context.getAttachments().putAll(attachments);
        Mesh.context().rewriteContext(context);

        Digest digest = new Digest();

        try {
            // Refresh slf4j MDC
            digest.mdc();
            Object ret = invoker.invoke(invocation);
            digest.print("P", MeshCode.SUCCESS.getCode());
            return ret;
        } catch (Throwable e) {
            if (e instanceof Codeable) {
                digest.print("P", ((Codeable) e).getCode());
            } else {
                digest.print("P", MeshCode.SYSTEM_ERROR.getCode());
            }
            if (e instanceof MeshException) {
                log.error("{},{},{}", digest.getTraceId(), Mesh.context().getUrn(), ((MeshException) e).getRootMessage());
            }
            throw e;
        } finally {
            MDC.clear();
        }
    }

    private long resolveTimestamp(String v) {
        if (Tool.optional(v)) {
            return System.currentTimeMillis();
        }
        try {
            return Long.parseLong(v);
        } catch (Throwable e) {
            log.error("Parse timestamp failed.", e);
            return System.currentTimeMillis();
        }
    }

    private Runmode resolveRunMode(String v) {
        if (Tool.optional(v)) {
            return Runmode.ROUTINE;
        }
        try {
            return Runmode.from(Integer.parseInt(v));
        } catch (Throwable e) {
            log.error("Parse run mode failed.", e);
            return Runmode.ROUTINE;
        }
    }
}
