/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.License;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Licenser {

    /**
     * Import the license.
     */
    @MPI(value = "mesh.license.imports", flags = 2)
    void imports(String license);

    /***
     * Exports the license.
     */
    @MPI(value = "mesh.license.exports", flags = 2)
    String exports();

    /**
     * Explain the license.
     */
    @MPI(value = "mesh.license.explain", flags = 2)
    License explain();

    /**
     * Verify the license.
     */
    @MPI(value = "mesh.license.verify", flags = 2)
    long verify();

    /**
     * License features.
     */
    @MPI(value = "mesh.license.features", flags = 2)
    Map<String, String> features();
}
