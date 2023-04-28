/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Cluster;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshCluster implements Cluster {

    private final Cluster cluster = ServiceProxy.proxy(Cluster.class);

    @Override
    public byte[] election(byte[] buff) {
        return cluster.election(buff);
    }

    @Override
    public boolean isLeader() {
        return cluster.isLeader();
    }
}
