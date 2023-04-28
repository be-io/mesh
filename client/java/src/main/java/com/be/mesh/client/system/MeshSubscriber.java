/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Subscriber;
import com.be.mesh.client.struct.Event;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshSubscriber implements Subscriber {

    private final Subscriber subscriber = ServiceProxy.proxy(Subscriber.class);

    @Override
    public void subscribe(Event event) {
        subscriber.subscribe(event);
    }
}
