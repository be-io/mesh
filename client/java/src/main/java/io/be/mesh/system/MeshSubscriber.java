/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Subscriber;
import io.be.mesh.struct.Event;

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
