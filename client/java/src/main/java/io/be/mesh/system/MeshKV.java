/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.KV;
import io.be.mesh.struct.Entity;

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
