/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Registry;
import io.be.mesh.struct.Registration;

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
