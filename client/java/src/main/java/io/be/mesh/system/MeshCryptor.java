/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Cryptor;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshCryptor implements Cryptor {

    private final Cryptor cryptor = ServiceProxy.proxy(Cryptor.class);

    @Override
    public byte[] encrypt(byte[] buff, Map<String, String> features) {
        return cryptor.encrypt(buff, features);
    }

    @Override
    public byte[] decrypt(byte[] buff, Map<String, String> features) {
        return cryptor.decrypt(buff, features);
    }

    @Override
    public byte[] hash(byte[] buff, Map<String, String> features) {
        return cryptor.hash(buff, features);
    }

    @Override
    public String sign(byte[] buff, Map<String, String> features) {
        return cryptor.sign(buff, features);
    }

    @Override
    public boolean verify(byte[] buff, Map<String, String> features) {
        return cryptor.verify(buff, features);
    }
}
