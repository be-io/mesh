/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Publisher;
import com.be.mesh.client.struct.Event;

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
