/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.tool.Tool;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.MDC;

import java.io.Serializable;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@Getter
@Slf4j(topic = "digest")
public class Digest implements Serializable {

    private static final long serialVersionUID = -135075634815713279L;
    private final String traceId = Mesh.context().getTraceId();
    private final String spanId = Mesh.context().getSpanId();
    private final String mode = String.valueOf(Mesh.context().getRunMode().getMode());
    private final String cni = Optional.ofNullable(Mesh.context().getConsumer().getNodeId()).orElse("");
    private final String cii = Optional.ofNullable(Mesh.context().getConsumer().getInstId()).orElse("");
    private final String cip = Optional.ofNullable(Mesh.context().getConsumer().getIp()).orElse("");
    private final String chost = Optional.ofNullable(Mesh.context().getConsumer().getHost()).orElse("");
    private final String pni = Optional.ofNullable(Mesh.context().getProvider().getNodeId()).orElse("");
    private final String pii = Optional.ofNullable(Mesh.context().getProvider().getInstId()).orElse("");
    private final String pip = Optional.ofNullable(Mesh.context().getProvider().getIp()).orElse("");
    private final String phost = Optional.ofNullable(Mesh.context().getProvider().getHost()).orElse("");
    private final String urn = Mesh.context().getUrn();
    private final long now = System.currentTimeMillis();

    public void print(String pattern, String code) {
        log.info("{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{}",
                traceId,
                spanId,
                Mesh.context().getTimestamp(),
                now,
                System.currentTimeMillis() - Mesh.context().getTimestamp(),
                System.currentTimeMillis() - now,
                mode,
                pattern,
                cni,
                cii,
                pni,
                pii,
                cip,
                pip,
                chost,
                phost,
                Tool.anyone(Mesh.context().getAttribute(Mesh.REMOTE), Tool.MESH_ADDRESS.get().any()),
                urn,
                code);
    }

    public void mdc() {
        // Refresh slf4j MDC
        MDC.put("tid", traceId);
        MDC.put("sid", spanId);
        MDC.put("mod", mode);
        MDC.put("cni", cni);
        MDC.put("cii", cii);
        MDC.put("pni", pni);
        MDC.put("pii", pii);
    }
}
