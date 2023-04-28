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
import com.be.mesh.client.struct.Versions;

import java.util.List;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Builtin {

    /**
     * Export the documents.
     *
     * @param name      mesh.name
     * @param formatter document formatter
     * @return formatted document
     */
    @MPI("${mesh.name}.builtin.doc")
    String doc(@Index(0) String name, @Index(1) String formatter);

    /**
     * Get the builtin application version.
     */
    @MPI("${mesh.name}.builtin.version")
    Versions version();

    /**
     * LogLevel set the application log level.
     */
    @MPI("${mesh.name}.builtin.debug")
    void debug(@Index(0) List<String> features);

    /**
     * Health check stats.
     */
    @MPI("${mesh.name}.builtin.stats")
    Map<String, String> stats(@Index(0) List<String> features);

}
