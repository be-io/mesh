/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Licenser;
import com.be.mesh.client.struct.License;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshLicenser implements Licenser {

    private final Licenser licenser = ServiceProxy.proxy(Licenser.class);

    @Override
    public void imports(String license) {
        licenser.imports(license);
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
