/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.struct.Registration;
import com.be.mesh.client.tool.UUID;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
public class RegisterTest {

    private final Registry registry = ServiceLoader.load(Registry.class).getDefault();

    @Test
    public void testRegister() {
        Registration registration = new Registration();
        registration.setInstanceId(UUID.getInstance().shortUUID());
        registration.setContent("Say hello.");
        registration.setKind("Metadata");
        registry.register(registration);
    }

    @Test
    public void testUnRegister() {
        Registration registration = new Registration();
        registration.setInstanceId(UUID.getInstance().shortUUID());
        registry.unregister(registration);
    }
}
