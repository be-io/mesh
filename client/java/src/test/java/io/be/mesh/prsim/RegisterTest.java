/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.mpc.ServiceLoader;
import io.be.mesh.struct.Registration;
import io.be.mesh.tool.UUID;
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
