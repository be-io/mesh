/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.Event;
import com.be.mesh.client.struct.Principal;
import com.be.mesh.client.struct.Topic;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

/**
 * Event publisher.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Publisher {

    /**
     * Publish event to mesh.
     *
     * @param events events
     * @return results
     */
    @MPI("mesh.queue.publish")
    List<String> publish(@Index(0) List<Event> events);

    /**
     * Synchronized broadcast the event to all subscriber. This maybe timeout with to many subscriber.
     *
     * @param events Event payload
     * @return Synchronized subscriber return value
     */
    @MPI("mesh.queue.broadcast")
    List<String> broadcast(@Index(0) List<Event> events);

    /**
     * Publish message to local node.
     *
     * @param binding topic
     * @param payload structure
     * @return eventId
     */
    default String publish(Topic binding, Object payload) {
        Event event = Event.newInstance(payload, binding);
        return this.publish(Collections.singletonList(event)).get(0);
    }

    /**
     * Unicast will publish to another node.
     *
     * @param binding topic
     * @param payload payload
     * @return eventId
     */
    default String unicast(Topic binding, Object payload, Principal principal) {
        Event event = Event.newInstance(payload, binding, principal);
        return this.publish(Collections.singletonList(event)).get(0);
    }

    /**
     * Multicast will publish event to principal groups.
     *
     * @param binding    topic.
     * @param payload    payload
     * @param principals principal group.
     * @return event id
     */
    default List<String> multicast(Topic binding, Object payload, List<Principal> principals) {
        List<Event> events = principals.stream().map(principal -> Event.newInstance(payload, binding, principal)).collect(Collectors.toList());
        return this.publish(events);
    }

    /**
     * Synchronized broadcast the event to all subscriber. This maybe timeout with to many subscriber.
     *
     * @param payload Event payload
     * @return Synchronized subscriber return value
     */
    default List<String> broadcast(Topic binding, Object payload) {
        return this.broadcast(Collections.singletonList(Event.newInstance(payload, binding)));
    }
}
