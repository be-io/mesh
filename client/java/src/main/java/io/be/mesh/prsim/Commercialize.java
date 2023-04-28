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
import io.be.mesh.struct.CommerceEnviron;
import io.be.mesh.struct.CommerceLicense;
import io.be.mesh.struct.License;

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
    String sign(@Index(0) License lsr);

    /**
     * History list the sign license in history, the latest is the first index.
     */
    @MPI(value = "mesh.license.history", flags = 2)
    List<CommerceLicense> history(@Index(value = 0, name = "inst_id") String instId);

    /**
     * Issued mesh node identity.
     */
    @MPI("mesh.net.issued")
    CommerceEnviron issued(@Index(0) String name, @Index(1) String kind);

    /**
     * Dump the node identity.
     */
    @MPI("mesh.net.dump")
    List<CommerceEnviron> dump(@Index(value = 0, name = "node_id") String nodeId);

}
