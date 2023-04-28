/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

/**
 * 编解码.
 *
 * @author coyzeng@gmail.com
 */
public interface Codec {

    /**
     * 密钥字符编码.
     *
     * @param key     密钥
     * @param headers 编码头
     * @return 密钥字符编码
     */
    String encodeKey(byte[] key, String... headers);

    /**
     * 密钥字符解码.
     *
     * @param key 密钥
     * @return 密钥字符解码
     */
    byte[] decodeKey(String key);

    /**
     * 签名字符编码.
     *
     * @param sig     签名
     * @param headers 编码头
     * @return 字符编码
     */
    String encodeSIG(byte[] sig, String... headers);

    /**
     * 签名字符解码.
     *
     * @param sig 签名
     * @return 字符解码
     */
    byte[] decodeSIG(String sig);

    /**
     * 报文密文字符编码.
     *
     * @param data    数据
     * @param headers 编码头
     * @return 密文字符编码
     */
    String encodeCipher(byte[] data, String... headers);

    /**
     * 报文密文字符解码.
     *
     * @param data 数据
     * @return 密文字符解码
     */
    byte[] decodeCipher(String data);

    /**
     * 报文明文字符编码.
     *
     * @param data    数据
     * @param headers 编码头
     * @return 明文字符编码
     */
    String encodePlain(byte[] data, String... headers);

    /**
     * 报文明文字符解码.
     *
     * @param data 数据
     * @return 明文字符解码
     */
    byte[] decodePlain(String data);

}
