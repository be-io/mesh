/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import com.be.mesh.client.cause.CryptException;
import com.be.mesh.client.mpc.MeshCode;

import javax.crypto.Cipher;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.security.*;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;

/**
 * @author coyzeng@gmail.com
 */
public class RSA2 implements Crypt {

    public static final int KEY_SIZE = 2048;
    public static final String KEY_ALGORITHM = "RSA";
    public static final String SIGNATURE_ALGORITHM = "SHA256withRSA";
    private static final String RSA_TYPE = "RSA/ECB/PKCS1Padding"; //"RSA/ECB/OAEPWITHSHA-256ANDMGF1PADDING";
    private static final int MAX_ENCRYPT_BLOCK_SIZE = 244; // RSA2最大加密明文大小(2048/8-11=244)
    private static final int MAX_DECRYPT_BLOCK_SIZE = 256; // RSA2最大解密密文大小(2048/8=256)

    @Override
    public KeyPair generateKeyPair() {
        try {
            KeyPairGenerator generator = KeyPairGenerator.getInstance(KEY_ALGORITHM);
            generator.initialize(KEY_SIZE, new SecureRandom());
            return generator.generateKeyPair();
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_ERROR, e);
        }
    }

    @Override
    public String signature(String data, String privateKey, Codec codec) {
        try {
            PrivateKey pk = getPrivateKey(privateKey, codec);
            Signature signature = Signature.getInstance(SIGNATURE_ALGORITHM);
            signature.initSign(pk);
            signature.update(codec.decodePlain(data));
            byte[] signed = signature.sign();
            return codec.encodeSIG(signed);
        } catch (Exception e) {
            throw new CryptException(MeshCode.SIGNATURE_ERROR, e);
        }
    }

    @Override
    public boolean verify(String data, String signature, String publicKey, Codec codec) {
        try {
            PublicKey pk = getPublicKey(publicKey, codec);
            Signature signed = Signature.getInstance(SIGNATURE_ALGORITHM);
            signed.initVerify(pk);
            signed.update(codec.decodePlain(data));
            return signed.verify(codec.decodeSIG(signature));
        } catch (Exception e) {
            throw new CryptException(MeshCode.SIGNATURE_ERROR, e);
        }
    }

    @Override
    public String encrypt(String data, String publicKey, Codec codec) {
        try {
            byte[] dataBytes = codec.decodePlain(data);
            Key key = getPublicKey(publicKey, codec);
            Cipher cipher = Cipher.getInstance(RSA_TYPE);
            cipher.init(Cipher.ENCRYPT_MODE, key);

            ByteArrayInputStream input = new ByteArrayInputStream(dataBytes);
            ByteArrayOutputStream output = new ByteArrayOutputStream();
            byte[] buffer = new byte[MAX_ENCRYPT_BLOCK_SIZE];
            for (int size = input.read(buffer); size > 0; size = input.read(buffer)) {
                byte[] encryptedBytes = cipher.doFinal(buffer, 0, size);
                output.write(encryptedBytes);
            }
            return codec.encodeCipher(output.toByteArray());
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_CODEC_ERROR, e);
        }
    }

    @Override
    public String decrypt(String data, String privateKey, Codec codec) {
        try {
            Key key = getPrivateKey(privateKey, codec);
            Cipher cipher = Cipher.getInstance(RSA_TYPE);
            cipher.init(Cipher.DECRYPT_MODE, key);

            byte[] encryptedBytes = codec.decodeCipher(data);
            ByteArrayInputStream input = new ByteArrayInputStream(encryptedBytes);
            ByteArrayOutputStream output = new ByteArrayOutputStream();
            byte[] buffer = new byte[MAX_DECRYPT_BLOCK_SIZE];
            for (int size = input.read(buffer); size > 0; size = input.read(buffer)) {
                byte[] decryptedBytes = cipher.doFinal(buffer, 0, size);
                output.write(decryptedBytes);
            }
            return codec.encodePlain(output.toByteArray());
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_CODEC_ERROR, e);
        }
    }

    /**
     * 获取RSA公钥.
     *
     * @param publicKey 公钥串
     * @return 公钥
     * @throws InvalidKeySpecException  公钥非法
     * @throws NoSuchAlgorithmException 算法SPI实现不存在
     */
    public PublicKey getPublicKey(String publicKey, Codec codec) throws InvalidKeySpecException, NoSuchAlgorithmException {
        byte[] byteKey = codec.decodeKey(publicKey);
        X509EncodedKeySpec spec = new X509EncodedKeySpec(byteKey);
        KeyFactory keyFactory = KeyFactory.getInstance(KEY_ALGORITHM);
        return keyFactory.generatePublic(spec);
    }

    /**
     * 获取RSA私钥.
     *
     * @param privateKey 私钥串
     * @return 公钥
     * @throws InvalidKeySpecException  公钥非法
     * @throws NoSuchAlgorithmException 算法SPI实现不存在
     */
    private PrivateKey getPrivateKey(String privateKey, Codec codec) throws NoSuchAlgorithmException, InvalidKeySpecException {
        byte[] byteKey = codec.decodeKey(privateKey);
        PKCS8EncodedKeySpec spec = new PKCS8EncodedKeySpec(byteKey);
        KeyFactory keyFactory = KeyFactory.getInstance(KEY_ALGORITHM);
        return keyFactory.generatePrivate(spec);
    }

}