/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Publisher;
import io.be.mesh.struct.Event;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshPublisher implements Publisher {

    private final Publisher publisher = ServiceProxy.proxy(Publisher.class);

    @Override
    public List<String> publish(List<Event> event) {
        return publisher.publish(event);
    }

    @Override
    public List<String> broadcast(List<Event> event) {
        return publisher.broadcast(event);
    }
}
