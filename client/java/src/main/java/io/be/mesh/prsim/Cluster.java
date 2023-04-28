/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.macro.Index;
import io.be.mesh.macro.MPI;
import io.be.mesh.macro.SPI;

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
