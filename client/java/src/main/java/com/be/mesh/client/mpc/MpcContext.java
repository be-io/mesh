/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.prsim.Network;
import com.be.mesh.client.struct.Environ;
import com.be.mesh.client.struct.Location;
import com.be.mesh.client.struct.Principal;
import com.be.mesh.client.tool.Tool;
import lombok.Setter;

import java.util.*;
import java.util.function.Supplier;

/**
 * @author coyzeng@gmail.com
 */
@Setter
public class MpcContext implements Context {

    private static final long serialVersionUID = 2668551432635194189L;
    private String traceId;
    private String spanId;
    private long timestamp;
    private Runmode runMode;
    private String urn;
    private Location consumer;
    private int calls;
    private final transient Map<String, String> attachments = new HashMap<>();
    private final transient Map<String, Object> attributes = new HashMap<>();
    private final transient Deque<Principal> principals = new ArrayDeque<>();

    /**
     * Create new invoke context.
     */
    public static MpcContext create() {
        MpcContext context = new MpcContext();
        context.setTraceId(Tool.newTraceId());
        context.setSpanId(Tool.newSpanId("", 0));
        context.setTimestamp(java.lang.System.currentTimeMillis());
        context.setRunMode(Runmode.ROUTINE);
        context.setUrn("");
        context.setConsumer(new Location());
        return context;
    }

    private static Location locale() {
        Network system = ServiceLoader.load(Network.class).getDefault();
        Environ environ = system.getEnviron();
        Location location = new Location();
        location.setInstId(environ.getInstId());
        location.setNodeId(environ.getNodeId());
        location.setIp(Tool.IP.get());
        location.setHost(Tool.HOST_NAME.get());
        location.setPort(String.valueOf(Tool.MESH_RUNTIME.get().getPort()));
        location.setName(Tool.MESH_NAME.get());
        return location;
    }

    @Override
    public String getTraceId() {
        return Optional.ofNullable(this.traceId).orElseGet(Tool::newTraceId);
    }

    @Override
    public String getSpanId() {
        return Optional.ofNullable(this.spanId).orElseGet(Tool::newTraceId);
    }

    @Override
    public long getTimestamp() {
        return 0 != this.timestamp ? this.timestamp : java.lang.System.currentTimeMillis();
    }

    @Override
    public Runmode getRunMode() {
        return Optional.ofNullable(this.runMode).orElse(Runmode.ROUTINE);
    }

    @Override
    public String getUrn() {
        return Optional.ofNullable(this.urn).orElse("");
    }

    @Override
    public Location getConsumer() {
        return Optional.ofNullable(this.consumer).orElseGet(this::getProvider);
    }

    @Override
    public Location getProvider() {
        return locale();
    }

    @Override
    public Map<String, String> getAttachments() {
        return this.attachments;
    }

    @Override
    public Deque<Principal> getPrincipals() {
        return this.principals;
    }

    @Override
    public Map<String, Object> getAttributes() {
        return this.attributes;
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T> T getAttribute(Key<T> key) {
        return (T) attributes.get(key.getName());
    }

    @Override
    public <T> void setAttribute(Key<T> key, T value) {
        attributes.put(key.getName(), value);
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T> T computeAttribute(Key<T> key, Supplier<T> supplier) {
        return (T) attributes.computeIfAbsent(key.getName(), k -> supplier.get());
    }

    @Override
    public void rewriteUrn(String urn) {
        this.urn = urn;
    }

    @Override
    public void rewriteContext(Context another) {
        if (Tool.required(another.getTraceId())) {
            this.setTraceId(another.getTraceId());
        }
        if (Tool.required(another.getSpanId())) {
            this.setSpanId(another.getSpanId());
        }
        if (0 != another.getTimestamp()) {
            this.setTimestamp(another.getTimestamp());
        }
        if (Runmode.ROUTINE != another.getRunMode()) {
            this.setRunMode(another.getRunMode());
        }
        if (Tool.required(another.getUrn())) {
            this.setUrn(another.getUrn());
        }
        if (Tool.required(another.getConsumer())) {
            this.setConsumer(another.getConsumer());
        }
        if (Tool.required(another.getAttachments())) {
            this.getAttachments().putAll(another.getAttachments());
        }
        if (Tool.required(another.getAttributes())) {
            this.getAttributes().putAll(another.getAttributes());
        }
        if (Tool.required(another.getPrincipals())) {
            this.getPrincipals().addAll(another.getPrincipals());
        }
    }

    @Override
    public MpcContext resume() {
        this.calls++;
        MpcContext context = new MpcContext();
        context.setTraceId(this.getTraceId());
        context.setSpanId(Tool.newSpanId(this.getSpanId(), this.calls));
        context.setTimestamp(this.getTimestamp());
        context.setRunMode(this.getRunMode());
        context.setUrn(this.getUrn());
        context.setConsumer(this.getConsumer());
        if (null != this.getAttachments()) {
            context.getAttachments().putAll(this.getAttachments());
        }
        if (null != this.getAttributes()) {
            context.getAttributes().putAll(this.getAttributes());
        }
        if (null != this.getPrincipals()) {
            context.getPrincipals().addAll(this.getPrincipals());
        }
        return context;
    }

    public static Context with(Context context) {
        return context.resume();
    }
}
