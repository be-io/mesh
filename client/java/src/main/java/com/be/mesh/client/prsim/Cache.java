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
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.struct.CacheEntity;
import com.be.mesh.client.struct.Entity;

import java.time.Duration;
import java.util.List;
import java.util.Optional;
import java.util.function.Function;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Cache {

    @MPI("mesh.cache.get")
    CacheEntity get(@Index(0) String key);

    @MPI("mesh.cache.put")
    void put(@Index(0) CacheEntity cell);

    @MPI("mesh.cache.remove")
    void remove(@Index(0) String key);

    @MPI("mesh.cache.incr")
    long incr(@Index(0) String key, @Index(1) long value);

    @MPI("mesh.cache.decr")
    long decr(@Index(0) String key, @Index(1) long value);

    /**
     * Keys the cache key set.
     */
    @MPI("mesh.cache.keys")
    List<String> keys(@Index(0) String pattern);

    /**
     * HGet get value in hash
     */
    @MPI("mesh.cache.hget")
    CacheEntity hget(String key, String name);

    /**
     * HSet put value in hash
     */
    @MPI("mesh.cache.hset")
    void hset(String key, CacheEntity cell);

    /**
     * HDel put value in hash
     */
    @MPI("mesh.cache.hdel")
    void hdel(String key, String name);

    /**
     * HKeys get the hash keys
     */
    @MPI("mesh.cache.hkeys")
    List<String> hkeys(String key);

    /**
     * Get by codec.
     */
    default <T> T get(String key, Types<T> type) {
        return Optional.ofNullable(this.get(key)).map(entity -> entity.getEntity().tryReadObject(type)).orElse(null);
    }

    /**
     * Put by codec.
     */
    default void put(String key, Object value, Duration duration) {
        CacheEntity cell = new CacheEntity();
        cell.setVersion("1.0.0");
        cell.setEntity(Entity.wrap(value));
        cell.setTimestamp(System.currentTimeMillis());
        cell.setDuration(duration.toMillis());
        cell.setKey(key);
        this.put(cell);
    }

    /**
     * Default compute and put if absent.
     */
    default <T> T computeIfAbsent(String key, Types<T> type, Duration duration, Function<String, T> fn) {
        T value = get(key, type);
        if (null != value) {
            return value;
        }
        value = fn.apply(key);
        if (null != value) {
            put(key, value, duration);
        }
        return value;
    }
}
