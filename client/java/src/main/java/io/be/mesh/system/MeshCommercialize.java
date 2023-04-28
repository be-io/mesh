/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Commercialize;
import io.be.mesh.struct.CommerceEnviron;
import io.be.mesh.struct.CommerceLicense;
import io.be.mesh.struct.License;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshCommercialize implements Commercialize {

    private final Commercialize commercialize = ServiceProxy.proxy(Commercialize.class);

    @Override
    public String sign(License lsr) {
        return commercialize.sign(lsr);
    }

    @Override
    public List<CommerceLicense> history(String instId) {
        return commercialize.history(instId);
    }

    @Override
    public CommerceEnviron issued(String name, String kind) {
        return commercialize.issued(name, kind);
    }

    @Override
    public List<CommerceEnviron> dump(String nodeId) {
        return commercialize.dump(nodeId);
    }
}
