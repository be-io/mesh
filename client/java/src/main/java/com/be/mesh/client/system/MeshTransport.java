/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.prsim.Session;
import com.be.mesh.client.prsim.Transport;
import com.be.mesh.client.tool.Tool;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;

import java.nio.ByteBuffer;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */

@Slf4j
@Setter
@SPI("mesh")
public class MeshTransport implements Transport, Session {

    private static final Map<String, Session> sessions = new ConcurrentHashMap<>();
    private final Session session = ServiceProxy.proxy(Session.class);
    private String sessionId;
    private Map<String, String> metadata;

    @Override
    public Session open(String sessionId, Map<String, String> metadata) {
        return sessions.computeIfAbsent(sessionId, sid -> {
            MeshTransport cs = new MeshTransport();
            cs.setSessionId(sessionId);
            cs.setMetadata(metadata);
            return cs;
        });
    }

    @Override
    public void close(Duration timeout) {
        sessions.forEach((sid, proxy) -> {
            try {
                proxy.release(timeout, "");
            } catch (Throwable e) {
                log.error(String.format("Close channel session %s with error.", sid), e);
            }
        });
        sessions.clear();
    }

    @Override
    public ByteBuffer peek(String topic) {
        withMetadata();
        return this.session.peek(topic);
    }

    @Override
    public ByteBuffer pop(Duration timeout, String topic) {
        withMetadata();
        return this.session.pop(timeout, topic);
    }

    @Override
    public void push(ByteBuffer payload, Map<String, String> metadata, String topic) {
        withMetadata();
        this.session.push(payload, metadata, topic);
    }

    @Override
    public void release(Duration timeout, String topic) {
        withMetadata();
        this.session.release(timeout, topic);
    }

    private void withMetadata() {
        if (Tool.required(this.metadata)) {
            Mesh.context().getAttachments().putAll(this.metadata);
        }
        if (Tool.required(this.sessionId)) {
            Context.Metadata.MESH_SESSION_ID.set(Mesh.context().getAttachments(), this.sessionId);
        }
    }
}
