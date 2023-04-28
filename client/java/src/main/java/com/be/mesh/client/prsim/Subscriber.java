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

/**
 * Event subscriber.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Subscriber {

    /**
     * Subscribe the event with {@link com.be.mesh.client.annotate.Bindings} or {@link com.be.mesh.client.annotate.Binding}
     *
     * @param event uniform event payload.
     */
    @MPI("mesh.queue.subscribe")
    void subscribe(@Index(0) Event event);

}
