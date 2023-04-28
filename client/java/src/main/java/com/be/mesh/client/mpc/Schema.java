/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("dynamic")
public interface Schema {

    /**
     * Import services to parse schema.
     *
     * @param types service interfaces.
     */
    void imports(List<Class<?>> types);

    /**
     * Import schema definition json to parse schema.
     *
     * @param schema schema definition.
     */
    void imports(String schema);

    /**
     * Export schema definition as json string.
     *
     * @return schema definition json schema.
     */
    String exports();

    /**
     * Search the by urn.
     *
     * @param urn Service urn.
     * @return
     */
    Class<?> search(String urn);
}
