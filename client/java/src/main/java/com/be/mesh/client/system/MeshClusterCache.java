/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Cache;
import com.be.mesh.client.struct.CacheEntity;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshClusterCache implements Cache {

    private final Cache cache = ServiceProxy.proxy(Cache.class);


    @Override
    public CacheEntity get(String key) {
        return cache.get(key);
    }

    @Override
    public void put(CacheEntity cell) {
        cache.put(cell);
    }

    @Override
    public void remove(String key) {
        cache.remove(key);
    }

    @Override
    public long incr(String key, long value) {
        return cache.incr(key, value);
    }

    @Override
    public long decr(String key, long value) {
        return cache.decr(key, value);
    }

    @Override
    public List<String> keys(String pattern) {
        return cache.keys(pattern);
    }

    @Override
    public CacheEntity hget(String key, String name) {
        return cache.hget(key, name);
    }

    @Override
    public void hset(String key, CacheEntity cell) {
        cache.hset(key, cell);
    }

    @Override
    public void hdel(String key, String name) {
        cache.hdel(key, name);
    }

    @Override
    public List<String> hkeys(String key) {
        return cache.hkeys(key);
    }

}
