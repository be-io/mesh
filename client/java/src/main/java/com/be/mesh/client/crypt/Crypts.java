/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import com.be.mesh.client.cause.CryptException;
import com.be.mesh.client.mpc.MeshCode;
import com.be.mesh.client.mpc.Types;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.*;
import java.util.stream.Collectors;

/**
 * Ca
 *
 * @author coyzeng@gmail.com
 */
public final class Crypts {

    private static final String SIGNATURE = "signature";
    private static final Crypt CA = new X509(true);
    private static final Codec CODEC = new PemCodec();

    private Crypts() {
    }

    /**
     * 生成基于国密RSA2算法PEM格式的密钥对和X509证书.
     *
     * @return PEM格式的密钥对和X509证书
     */
    public static Map<String, String> generateX509CertificatePemFmt() {
        return CA.generateKeys(new PemCodec());
    }

    /**
     * 生成基于国密SM2算法PEM格式的密钥对和X509证书到对应目录.
     */
    public static void generateX509Certificate(String path) {
        try {
            Files.createDirectories(Paths.get(path));
            Map<String, String> pem = generateX509CertificatePemFmt();
            for (Map.Entry<String, String> entry : pem.entrySet()) {
                Files.write(Paths.get(String.format("%s%s%s", path, File.separator, entry.getKey())), entry.getValue().getBytes(StandardCharsets.UTF_8));
            }
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_ERROR, e);
        }
    }

    /**
     * 为报文添加签名并加密. PaddingSignature+Encrypt
     *
     * @param mapper     请求报文
     * @param privateKey 平台私钥
     * @param publicKey  对方公钥
     * @param codec      平台自己的JSON序列化库实现
     * @return 填充好签名的加密报文
     */
    public static String pse(Mapper mapper, String privateKey, String publicKey, com.be.mesh.client.mpc.Codec codec) {
        Map<String, String> padding = new HashMap<>(mapper.map().size());
        mapper.map().forEach((k, v) -> {
            if (Objects.nonNull(v)) {
                padding.put(k, codec.encodeString(v));
            }
        });
        Map<String, String> sorted = new TreeMap<>(padding);
        String chain = sorted.entrySet().stream().
                map(e -> String.format("%s=%s", e.getKey(), e.getValue())).
                collect(Collectors.joining("&"));
        String signature = CA.signature(chain, privateKey, CODEC);
        padding.put(SIGNATURE, signature);
        return CA.encrypt(codec.encodeString(padding), publicKey, CODEC);
    }

    /**
     * 解密并验证报文中的签名. Decrypt+VerifySignature
     *
     * @param data       请求报文密文（包含签名）
     * @param publicKey  对方公钥
     * @param privateKey 平台私钥
     * @param codec      平台自己的JSON序列化库实现
     * @return 原始报文
     */
    public static <T> T dvs(String data, String publicKey, String privateKey, com.be.mesh.client.mpc.Codec codec, Types<T> type) {
        String param = decrypt(data, privateKey);
        if (!verify(param, codec, publicKey)) {
            throw new java.lang.SecurityException("签名非法");
        }
        return codec.decodeString(param, type);
    }

    /**
     * 私钥解密.
     *
     * @param data       密文
     * @param privateKey 私钥
     * @return 原文
     */
    public static String decrypt(String data, String privateKey) {
        return CA.decrypt(data, privateKey, CODEC);
    }

    /**
     * Encrypt data with public key.
     *
     * @param data      explain data
     * @param publicKey public key
     * @return encrypted data
     */
    public static String encrypt(String data, String publicKey) {
        return CA.encrypt(data, publicKey, CODEC);
    }

    /**
     * 验签
     *
     * @param data      原文
     * @param codec     编解码
     * @param publicKey 公钥
     * @return true通过
     */
    public static boolean verify(String data, com.be.mesh.client.mpc.Codec codec, String publicKey) {
        Map<String, Object> sorted = new TreeMap<>(codec.decodeString(data, Types.MapObject));
        String signature = Optional.ofNullable(sorted.remove(SIGNATURE)).map(String::valueOf).orElse("");
        String chain = sorted.entrySet().stream().
                filter(e -> Objects.nonNull(e.getValue())).
                map(e -> String.format("%s=%s", e.getKey(), codec.encodeString(e.getValue()))).
                collect(Collectors.joining("&"));
        return CA.verify(chain, signature, publicKey, CODEC);
    }


}
