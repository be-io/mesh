/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.KV;
import com.be.mesh.client.struct.Entity;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshKV implements KV {

    private final KV kv = ServiceProxy.proxy(KV.class);

    @Override
    public Entity get(String key) {
        return kv.get(key);
    }

    @Override
    public void put(String key, Entity value) {
        kv.put(key, value);
    }

    @Override
    public void remove(String key) {
        kv.remove(key);
    }
}
