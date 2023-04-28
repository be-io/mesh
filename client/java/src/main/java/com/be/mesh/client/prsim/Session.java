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

import java.nio.ByteBuffer;
import java.time.Duration;
import java.util.Map;

/**
 * Remote queue in async and blocking mode.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Session {

    /**
     * Retrieves, but does not remove, the head of this queue,
     * or returns None if this queue is empty.
     */
    @MPI("mesh.chan.peek")
    ByteBuffer peek(String topic);

    /**
     * Retrieves and removes the head of this queue,
     * or returns None if this queue is empty.
     */
    @MPI("mesh.chan.pop")
    ByteBuffer pop(@Index(0) Duration timeout, String topic);

    /**
     * Inserts the specified element into this queue if it is possible to do
     * so immediately without violating capacity restrictions.
     * When using a capacity-restricted queue, this method is generally
     * preferable to add, which can fail to insert an element only
     * by throwing an exception.
     */
    @MPI("mesh.chan.push")
    void push(@Index(0) ByteBuffer payload, Map<String, String> metadata, String topic);

    /**
     * Release the channel session.
     */
    @MPI("mesh.chan.release")
    void release(@Index(0) Duration timeout, String topic);
}
