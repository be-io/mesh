/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.ServiceLoader;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class CryptorTest {

    @Test
    public void testSign() {
        //System.setProperty("mesh.address", "10.90.21.174");
        byte[] explain = "123456".getBytes(StandardCharsets.UTF_8);
        Cryptor cryptor = ServiceLoader.load(Cryptor.class).getDefault();
        String signature = cryptor.sign(explain, new HashMap<>());
        log.info(signature);
        Map<String, String> features = new HashMap<>();
        features.put("signature", signature);
        log.info("{}", cryptor.verify(explain, features));
    }
}
