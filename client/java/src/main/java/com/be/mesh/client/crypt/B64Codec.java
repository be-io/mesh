/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import java.util.Base64;

/**
 * @author coyzeng@gmail.com
 */
class B64Codec implements Codec {

    private final Codec codec = new StringCodec();

    @Override
    public String encodeKey(byte[] key, String... headers) {
        return this.encodeSIG(key, headers);
    }

    @Override
    public byte[] decodeKey(String key) {
        return this.decodeSIG(key);
    }

    @Override
    public String encodeSIG(byte[] sig, String... headers) {
        return Base64.getEncoder().encodeToString(sig);
    }

    @Override
    public byte[] decodeSIG(String sig) {
        return Base64.getDecoder().decode(sig);
    }

    @Override
    public String encodeCipher(byte[] data, String... headers) {
        return Base64.getEncoder().encodeToString(data);
    }

    @Override
    public byte[] decodeCipher(String data) {
        return Base64.getDecoder().decode(data);
    }

    @Override
    public String encodePlain(byte[] data, String... headers) {
        return this.codec.encodePlain(data, headers);
    }

    @Override
    public byte[] decodePlain(String data) {
        return this.codec.decodePlain(data);
    }
}
