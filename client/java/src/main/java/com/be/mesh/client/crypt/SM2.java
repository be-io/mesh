/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import com.be.mesh.client.cause.CryptException;
import com.be.mesh.client.mpc.MeshCode;
import org.bouncycastle.jcajce.spec.SM2ParameterSpec;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.util.Strings;

import javax.crypto.Cipher;
import java.security.*;
import java.security.spec.ECGenParameterSpec;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

/**
 * @author coyzeng@gmail.com
 */
public class SM2 implements Crypt {

    private static final String DEFAULT_USER_ID = "1234567812345678"; // SM2算法默认用户ID
    private static final int DEFAULT_KEY_SIZE = 128; // SM4算法目前只支持128位（即密钥16字节）
    public static final String KEY_ALGORITHM = "EC";
    private static final String ALGORITHM_ECB_PKCS5PADDING = "SM2";
    public static final String SIGNATURE_ALGORITHM = "SM3withSM2";

    static {
        if (null == Security.getProvider(BouncyCastleProvider.PROVIDER_NAME)) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    @Override
    public KeyPair generateKeyPair() {
        try {
            KeyPairGenerator generator = KeyPairGenerator.getInstance(KEY_ALGORITHM, BouncyCastleProvider.PROVIDER_NAME);
            generator.initialize(new ECGenParameterSpec("sm2p256v1")); // ECC
            return generator.generateKeyPair();
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_ERROR, e);
        }
    }

    @Override
    public String signature(String data, String privateKey, Codec codec) {
        try {
            PrivateKey key = getPrivateKey(privateKey, codec);
            Signature signature = Signature.getInstance(SIGNATURE_ALGORITHM, BouncyCastleProvider.PROVIDER_NAME);
            signature.setParameter(new SM2ParameterSpec(Strings.toByteArray(DEFAULT_USER_ID)));
            signature.initSign(key);
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
            Signature sm2sig = Signature.getInstance(SIGNATURE_ALGORITHM);
            sm2sig.setParameter(new SM2ParameterSpec(Strings.toByteArray(DEFAULT_USER_ID)));
            sm2sig.initVerify(pk);
            sm2sig.update(data.getBytes());
            return sm2sig.verify(Base64.getDecoder().decode(signature));
        } catch (Exception e) {
            throw new CryptException(MeshCode.SIGNATURE_ERROR, e);
        }
    }

    @Override
    public String encrypt(String data, String publicKey, Codec codec) {
        try {
            byte[] dataBytes = codec.decodePlain(data);
            PublicKey key = getPublicKey(publicKey, codec);
            Cipher sm2Cipher = Cipher.getInstance(ALGORITHM_ECB_PKCS5PADDING, BouncyCastleProvider.PROVIDER_NAME);
            sm2Cipher.init(Cipher.ENCRYPT_MODE, key);
            byte[] encryptBytes = sm2Cipher.doFinal(dataBytes);
            return codec.encodeCipher(encryptBytes);
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_CODEC_ERROR, e);
        }
    }

    @Override
    public String decrypt(String data, String privateKey, Codec codec) {
        try {
            byte[] encryptedBytes = codec.decodeCipher(data);
            PrivateKey key = getPrivateKey(privateKey, codec);
            Cipher sm2Cipher = Cipher.getInstance(ALGORITHM_ECB_PKCS5PADDING, BouncyCastleProvider.PROVIDER_NAME);
            sm2Cipher.init(Cipher.DECRYPT_MODE, key);
            byte[] decryptedBytes = sm2Cipher.doFinal(encryptedBytes);
            return codec.encodePlain(decryptedBytes);
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
        KeyFactory keyFactory = KeyFactory.getInstance(KEY_ALGORITHM);
        X509EncodedKeySpec keySpec = new X509EncodedKeySpec(byteKey);
        return keyFactory.generatePublic(keySpec);
    }

    /**
     * 获取私钥.
     *
     * @param privateKey 私钥串
     * @return 公钥
     * @throws InvalidKeySpecException  公钥非法
     * @throws NoSuchAlgorithmException 算法SPI实现不存在
     */
    private PrivateKey getPrivateKey(String privateKey, Codec codec) throws NoSuchAlgorithmException, InvalidKeySpecException {
        byte[] pkBytes = codec.decodeKey(privateKey);
        KeyFactory keyFactory = KeyFactory.getInstance(KEY_ALGORITHM);
        PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(pkBytes);
        return keyFactory.generatePrivate(keySpec);
    }
}
