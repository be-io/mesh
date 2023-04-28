/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.prsim.Cache;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class ProxyTest {

    @Test
    public void testProxy() {
        Cache cache = ServiceProxy.proxy(Cache.class);
        Map<String, Object> ret = cache.get("x", Types.MapObject);
        log.info("{}", ret);
    }

    @Test
    public void testProxyDefault() {
        Cache cache = ServiceProxy.proxy(Cache.class);
        String v = cache.get("t", Types.of(String.class));
        log.info(v);
    }
}
