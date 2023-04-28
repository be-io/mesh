/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.SPI;

import java.util.List;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Hodor {

    /**
     * Collect the system, application, process or thread status etc.
     *
     * @param features customized parameters
     * @return quota stat
     */
    Map<String, String> stats(List<String> features);

    /**
     * Set debug features.
     */
    void debug(List<String> features);

}
