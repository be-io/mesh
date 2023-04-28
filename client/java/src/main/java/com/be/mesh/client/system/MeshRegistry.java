/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Registry;
import com.be.mesh.client.struct.Registration;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshRegistry implements Registry {

    private final Registry registry = ServiceProxy.proxy(Registry.class);


    @Override
    public void register(Registration registration) {
        registry.register(registration);
    }

    @Override
    public void registers(List<Registration> registrations) {
        registry.registers(registrations);
    }

    @Override
    public void unregister(Registration registration) {
        registry.unregister(registration);
    }

    @Override
    public List<Registration> export(String kind) {
        return registry.export(kind);
    }
}
