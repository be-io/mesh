/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.macro.MPI;
import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.Types;
import io.be.mesh.struct.Entity;

import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface KV {

    // Get the value from kv store.
    @MPI("mesh.kv.get")
    Entity get(String key);

    // Put the value to kv store.
    @MPI("mesh.kv.put")
    void put(String key, Entity value);

    // Remove the kv store.
    @MPI("mesh.kv.remove")
    void remove(String key);

    /**
     * Get by codec.
     */
    default <T> T get(String key, Types<T> type) {
        return Optional.ofNullable(this.get(key)).orElseGet(Entity::empty).tryReadObject(type);
    }

    /**
     * Put by codec.
     */
    default void put(String key, Object value) {
        this.put(key, Entity.wrap(value));
    }
}
