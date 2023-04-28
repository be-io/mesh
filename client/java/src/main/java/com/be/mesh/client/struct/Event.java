/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.Network;
import com.be.mesh.client.tool.UUID;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Event implements Serializable {

    private static final long serialVersionUID = 3942244848716494415L;
    /**
     * Event version
     */
    @Index(0)
    private String version;
    /**
     * Trace ID.
     */
    @Index(5)
    private String tid;
    /**
     * Span ID.
     */
    @Index(10)
    private String sid;
    /**
     * Event ID
     */
    @Index(15)
    private String eid;
    /**
     * Message ID
     */
    @Index(20)
    private String mid;
    /**
     * Event published time as {@link System#currentTimeMillis()}
     */
    @Index(25)
    private String timestamp;
    /**
     * Event own principal information.
     */
    @Index(30)
    private Principal source;
    /**
     * Event own principal information.
     */
    @Index(35)
    private Principal target;
    /**
     * Event binding tuple.
     */
    @Index(40)
    private Topic binding;
    /**
     * Payload codec.
     */
    @Index(45)
    private Entity entity;

    /**
     * Create local event instance.
     */
    public static Event newInstance(Object payload, Topic topic) {
        Network network = ServiceLoader.load(Network.class).getDefault();
        Principal target = new Principal();
        target.setNodeId(network.getEnviron().getNodeId());
        target.setInstId(network.getEnviron().getInstId());
        return newInstance(payload, topic, target);
    }

    /**
     * Create any node event instance.
     */
    public static Event newInstance(Object payload, Topic topic, Principal target) {
        Network network = ServiceLoader.load(Network.class).getDefault();
        Principal source = new Principal();
        source.setNodeId(network.getEnviron().getNodeId());
        source.setInstId(network.getEnviron().getInstId());
        return newInstance(payload, topic, target, source);
    }

    /**
     * Create an event with source node and target node.
     */
    public static Event newInstance(Object payload, Topic topic, Principal target, Principal source) {
        Event event = new Event();
        event.setVersion("1.0.0");
        event.setTid(Mesh.context().getTraceId());
        event.setSid(Mesh.context().getSpanId());
        event.setEid(UUID.getInstance().shortUUID());
        event.setMid(UUID.getInstance().shortUUID());
        event.setTimestamp(String.valueOf(System.currentTimeMillis()));
        event.setSource(source);
        event.setTarget(target);
        event.setBinding(topic);
        event.setEntity(Entity.wrap(payload));
        return event;
    }


}
