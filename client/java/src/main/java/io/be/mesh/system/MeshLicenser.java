/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Licenser;
import io.be.mesh.struct.License;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshLicenser implements Licenser {

    private final Licenser licenser = ServiceProxy.proxy(Licenser.class);

    @Override
    public void imports(String text) {
        licenser.imports(text);
    }

    @Override
    public String exports() {
        return licenser.exports();
    }

    @Override
    public License explain() {
        return licenser.explain();
    }

    @Override
    public long verify() {
        return licenser.verify();
    }

    @Override
    public Map<String, String> features() {
        return licenser.features();
    }
}
