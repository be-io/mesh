/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import com.be.mesh.client.cause.CryptException;
import com.be.mesh.client.mpc.MeshCode;
import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemReader;
import org.bouncycastle.util.io.pem.PemWriter;

import java.io.StringReader;
import java.io.StringWriter;

/**
 * OpenSSL PEM 格式文件字符串的 SSL 证书请求 CSR 文件内容
 *
 * @author coyzeng@gmail.com
 */
class PemCodec implements Codec {

    private final Codec codec = new B64Codec();

    @Override
    public String encodeKey(byte[] key, String... headers) {
        String type = null != headers && headers.length > 0 ? headers[0] : "CERTIFICATE";
        PemObject pem = new PemObject(type, key);
        try (StringWriter writer = new StringWriter()) {
            PemWriter pw = new PemWriter(writer);
            pw.writeObject(pem);
            pw.close();
            return writer.toString();
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_ERROR, e);
        }
    }

    @Override
    public byte[] decodeKey(String key) {
        try (StringReader reader = new StringReader(key)) {
            PemReader pr = new PemReader(reader);
            return pr.readPemObject().getContent();
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_ERROR, e);
        }
    }

    @Override
    public String encodeSIG(byte[] sig, String... headers) {
        return codec.encodeKey(sig, headers);
    }

    @Override
    public byte[] decodeSIG(String sig) {
        return codec.decodeKey(sig);
    }

    @Override
    public String encodeCipher(byte[] data, String... headers) {
        return this.codec.encodeCipher(data, headers);
    }

    @Override
    public byte[] decodeCipher(String data) {
        return this.codec.decodeCipher(data);
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
