/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import java.nio.charset.StandardCharsets;

/**
 * @author coyzeng@gmail.com
 */
class StringCodec implements Codec {

    @Override
    public String encodeKey(byte[] key, String... headers) {
        return this.encodeCipher(key, headers);
    }

    @Override
    public byte[] decodeKey(String key) {
        return this.decodeCipher(key);
    }

    @Override
    public String encodeSIG(byte[] sig, String... headers) {
        return this.encodeCipher(sig, headers);
    }

    @Override
    public byte[] decodeSIG(String sig) {
        return this.decodeCipher(sig);
    }

    @Override
    public String encodeCipher(byte[] data, String... headers) {
        return new String(data, StandardCharsets.UTF_8);
    }

    @Override
    public byte[] decodeCipher(String data) {
        return data.getBytes(StandardCharsets.UTF_8);
    }

    @Override
    public String encodePlain(byte[] data, String... headers) {
        return this.encodeCipher(data, headers);
    }

    @Override
    public byte[] decodePlain(String data) {
        return data.getBytes(StandardCharsets.UTF_8);
    }
}
