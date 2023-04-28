/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.Entity;

import java.time.Duration;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
@SuppressWarnings("unchecked")
public interface Queue<T> {

    /**
     * Retrieves, but does not remove, the head of this queue, or returns {@code null} if this queue is empty.
     *
     * @param channel queue channel
     * @return the head of this queue, or None if this queue is empty
     */
    @SPI("mesh.queue.peek")
    Entity peek(@Index(0) String channel);

    /**
     * Retrieves and removes the head of this queue, or returns None if this queue is empty.
     *
     * @param channel queue channel
     * @param timeout pop timeout
     * @return the head of this queue, or None if this queue is empty
     */
    @SPI("mesh.queue.pop")
    Entity pop(@Index(0) String channel, @Index(2) Duration timeout);

    /**
     * Inserts the specified element into this queue if it is possible to do
     * so immediately without violating capacity restrictions.
     * When using a capacity-restricted queue, this method is generally
     * preferable to add, which can fail to insert an element only
     * by throwing an exception.
     *
     * @param channel queue channel
     * @param element pop timeout
     */
    @SPI("mesh.queue.push")
    void push(@Index(0) String channel, @Index(2) Entity element);

    /**
     * Peek the instance
     */
    default T pick(@Index(0) String channel) {
        Entity entity = this.peek(channel);
        return Optional.ofNullable(entity).map(x -> (T) x.readObject()).orElse(null);
    }

    /**
     * Pop the instance
     */
    default T poll(@Index(0) String channel, @Index(2) Duration timeout) {
        Entity entity = this.pop(channel, timeout);
        return Optional.ofNullable(entity).map(x -> (T) x.readObject()).orElse(null);
    }

    /**
     * Push the instance
     */
    default void offer(@Index(0) String channel, @Index(2) T element) {
        this.push(channel, Entity.wrap(element));
    }

}
