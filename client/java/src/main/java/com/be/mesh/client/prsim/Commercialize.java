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
import com.be.mesh.client.struct.CommerceEnviron;
import com.be.mesh.client.struct.CommerceLicense;
import com.be.mesh.client.struct.License;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Commercialize {

    /**
     * License sign.
     */
    @MPI(value = "mesh.license.sign", flags = 2)
    CommerceLicense sign(@Index(0) License lsr);

    /**
     * History list the sign license in history, the latest is the first index.
     */
    @MPI(value = "mesh.license.history", flags = 2)
    List<CommerceLicense> history(@Index(value = 0, name = "inst_id") String instId);

    /**
     * Issued mesh node identity.
     */
    @MPI("mesh.net.issued")
    CommerceEnviron issued(@Index(0) String name, @Index(1) String kind, @Index(2) String cname);

    /**
     * Dump the node identity.
     */
    @MPI("mesh.net.dump")
    List<CommerceEnviron> dump(@Index(value = 0, name = "node_id") String nodeId);

}
