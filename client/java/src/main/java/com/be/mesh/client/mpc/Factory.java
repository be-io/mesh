/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;

import java.io.InputStream;
import java.util.stream.Stream;

/**
 * @author coyzeng@gmail.com
 */
@SPI("java")
public interface Factory {

    String WARN = "No implementation defined in /META-INF/services/%s, please check whether the file exists and has the right implementation class!";

    /**
     * Gets Extension.
     *
     * @param <T> the type parameter
     * @param key the key
     * @param spi the clazz
     * @return the extension
     */
    <T> T getProvider(String key, Class<T> spi);

    /**
     * Get the java spi extension.
     *
     * @param <T> the type
     * @param spi the spi interface
     * @return the clazz
     */
    <T> Stream<T> getProvider(Class<T> spi);

    /**
     * Get the resource from the plugin META-INF/janus.xxx.yaml or remote.
     *
     * @param name plugin name
     * @return profile object
     */
    InputStream getResource(String name);

}
