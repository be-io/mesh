/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.macro.Index;
import io.be.mesh.macro.SPI;

import java.util.Map;

/**
 * Generic reference.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Dispatcher {

    /**
     * Refer an generic dispatcher.
     */
    <T> T reference(Class<T> mpi);

    Object invoke(@Index(0) String urn, @Index(1) Map<String, Object> param);

    Object invoke(@Index(0) String urn, @Index(1) Object param);
}
