/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import com.be.mesh.client.cause.CryptException;
import com.be.mesh.client.mpc.MeshCode;

import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.HashMap;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
public interface Crypt {

    /**
     * 生成密钥对.
     *
     * @return 密钥对
     */
    KeyPair generateKeyPair();

    /**
     * 生成RSA公私钥
     *
     * @param codec 编解码
     * @return key公钥, value私钥
     */
    default Map<String, String> generateKeys(Codec codec) {
        try {
            KeyPair keyPair = this.generateKeyPair();
            PublicKey pubKey = keyPair.getPublic();
            PrivateKey priKey = keyPair.getPrivate();
            Map<String, String> keyStore = new HashMap<>(2);
            keyStore.put("public.key", codec.encodeKey(pubKey.getEncoded()));
            keyStore.put("private.key", codec.encodeKey(priKey.getEncoded()));
            return keyStore;
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_ERROR, e);
        }
    }

    /**
     * 私钥签名.
     *
     * @param data       加签数据
     * @param privateKey 私钥
     * @param codec      编解码
     * @return 签名
     */
    String signature(String data, String privateKey, Codec codec);

    /**
     * 公钥验签.
     *
     * @param data      加签数据
     * @param signature 签名
     * @param publicKey 公钥
     * @param codec     编解码
     * @return 验签通过与否
     */
    boolean verify(String data, String signature, String publicKey, Codec codec);

    /**
     * 公钥加密.
     *
     * @param data      待加密数据
     * @param publicKey 公钥
     * @param codec     编解码
     * @return 已加密数据
     */
    String encrypt(String data, String publicKey, Codec codec);

    /**
     * 私钥解密.
     *
     * @param data       已加密数据
     * @param privateKey 公钥
     * @param codec      编解码
     * @return 待加密数据
     */
    String decrypt(String data, String privateKey, Codec codec);
}
