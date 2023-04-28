/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

/**
 * @author coyzeng@gmail.com
 */
public interface RuntimeHook {

    /**
     * Trigger when mesh runtime is start.
     */
    default void start() throws Throwable {

    }

    /**
     * Trigger when mesh runtime is stop.
     */
    default void stop() throws Throwable {

    }

    /**
     * Trigger then mesh runtime context is refresh or metadata is refresh.
     */
    default void refresh() throws Throwable {

    }

}
