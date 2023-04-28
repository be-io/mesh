/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Commercialize;
import com.be.mesh.client.struct.CommerceEnviron;
import com.be.mesh.client.struct.CommerceLicense;
import com.be.mesh.client.struct.License;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshCommercialize implements Commercialize {

    private final Commercialize commercialize = ServiceProxy.proxy(Commercialize.class);

    @Override
    public CommerceLicense sign(License lsr) {
        return commercialize.sign(lsr);
    }

    @Override
    public List<CommerceLicense> history(String instId) {
        return commercialize.history(instId);
    }

    @Override
    public CommerceEnviron issued(String name, String kind, String cname) {
        return commercialize.issued(name, kind, cname);
    }

    @Override
    public List<CommerceEnviron> dump(String nodeId) {
        return commercialize.dump(nodeId);
    }
}
