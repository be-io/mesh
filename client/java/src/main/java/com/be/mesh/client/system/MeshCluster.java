/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Cluster;

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
