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

import java.time.Duration;
import java.util.Map;

/**
 * Private compute data channel in async and blocking mode.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Transport {

    /**
     * Open a channel session.
     */
    @MPI("mesh.chan.open")
    Session open(@Index(0) String sessionId, @Index(1) Map<String, String> metadata);

    /**
     * Close the channel.
     */
    @MPI("mesh.chan.close")
    void close(@Index(0) Duration timeout);
}
