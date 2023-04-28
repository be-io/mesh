/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;

import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import org.testng.annotations.Test;

import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

/**
 * @author coyzeng@gmail.com
 */
public class CryptTest implements Mapper {

    @Test
    public void test() throws Exception {
        test(new RSA2(), new B64Codec(), "1111");
        test(new X509(true), new PemCodec(), "1111");
        test(new SM2(), new B64Codec(), "1111");

        pipeline();
    }

    private void pipeline() throws Exception {
        String platform = String.format("%s%s", "/Users/ducer/Desktop/", "cat0");
        String operator = String.format("%s%s", "/Users/ducer/Desktop/", "cat1");
        System.out.println("请求流程测试：");
        System.out.println(platform);
        System.out.println(operator);
        com.be.mesh.client.mpc.Codec codec = ServiceLoader.load(com.be.mesh.client.mpc.Codec.class).get("json");
        Crypts.generateX509Certificate(platform);
        Crypts.generateX509Certificate(operator);
        String platformPubKey = new String(Files.readAllBytes(Paths.get(String.format("%s%s%s", platform, File.separator, "public.key"))));
        String platformPriKey = new String(Files.readAllBytes(Paths.get(String.format("%s%s%s", platform, File.separator, "private.key"))));
        String operatorPubKey = new String(Files.readAllBytes(Paths.get(String.format("%s%s%s", operator, File.separator, "public.key"))));
        String operatorPriKey = new String(Files.readAllBytes(Paths.get(String.format("%s%s%s", operator, File.separator, "private.key"))));

        CryptTest test = new CryptTest();
        test.param.put("charset", "UTF-8");
        test.param.put("requestId", UUID.randomUUID().toString().replaceAll("-", ""));
        test.param.put("timestamp", System.currentTimeMillis());
        test.param.put("content", new HashMap<>());
        String encrypted = Crypts.pse(test, operatorPriKey, platformPubKey, codec);
        Map<String, Object> param1 = Crypts.dvs(encrypted, operatorPubKey, platformPriKey, codec, Types.MapObject);
        param1.forEach((k, v) -> System.out.printf("%s=%s\n", k, v));
    }

    private void test(Crypt ca, Codec codec, String data) {
        Map<String, String> pair = ca.generateKeys(codec);
        String privateKey = pair.get("private.key");
        String publicKey = pair.get("public.key");
        System.out.printf("%s公钥：%n", ca.getClass().getSimpleName());
        System.out.println(publicKey);
        System.out.printf("%s私钥：%n", ca.getClass().getSimpleName());
        System.out.println(privateKey);
        System.out.printf("%s签名：%n", ca.getClass().getSimpleName());
        String signature = ca.signature(data, privateKey, codec);
        System.out.println(signature);
        System.out.printf("%s验签：%n", ca.getClass().getSimpleName());
        System.out.println(ca.verify(data, signature, publicKey, codec));
        System.out.printf("%s原文：%n", ca.getClass().getSimpleName());
        System.out.println(data);
        System.out.printf("%s密文：%n", ca.getClass().getSimpleName());
        String enc = ca.encrypt(data, publicKey, codec);
        System.out.println(enc);
        System.out.printf("%s解密：%n", ca.getClass().getSimpleName());
        System.out.println(ca.decrypt(enc, privateKey, codec));
    }


    @Override
    public Map<String, Object> map() {
        return this.param;
    }

    private final Map<String, Object> param = new HashMap<>();
}