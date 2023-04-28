/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.macro.*;
import io.be.mesh.struct.Event;

/**
 * Event subscriber.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Subscriber {

    /**
     * Subscribe the event with {@link Bindings} or {@link Binding}
     *
     * @param event uniform event payload.
     */
    @MPI("mesh.queue.subscribe")
    void subscribe(@Index(0) Event event);

}
