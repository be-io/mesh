/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Cluster {

    /**
     * Election will election leader of instances.
     */
    @MPI("mesh.cluster.election")
    byte[] election(@Index(0) byte[] buff);

    /**
     * IsLeader if same level.
     */
    @MPI("mesh.cluster.leader")
    boolean isLeader();

}
