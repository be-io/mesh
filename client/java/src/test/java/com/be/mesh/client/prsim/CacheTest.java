/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.testng.Assert;
import org.testng.annotations.Test;

import java.time.Duration;
import java.util.concurrent.locks.LockSupport;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class CacheTest<T> {

    @Test
    public void testMeshCache() {
        Target<String, String> target = new Target<>(Duration.ofMillis(100), Types.of(String.class), "");
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        String text = codec.encodeString(target);
        log.info(text);
        Target<String, String> next = codec.decodeString(text, new Types<Target<String, String>>() {
        });
        log.info("{}", next.duration);
        System.setProperty("mesh.address", "127.0.0.1");
        Cache cache = ServiceLoader.load(Cache.class).getDefault();
        cache.put("mesh-cache-ping", "pong", Duration.ofSeconds(5));
        String pong = cache.get("mesh-cache-ping", Types.of(String.class));
        cache.get("mesh-cache-ping", Types.of(String.class));
        cache.get("mesh-cache-ping", Types.of(String.class));
        Assert.assertEquals(pong, "pong");
        LockSupport.parkNanos(Duration.ofSeconds(10).toNanos());
        Assert.assertNull(pong);
    }

    private interface BB<X, V> {

    }

    @NoArgsConstructor
    @AllArgsConstructor
    private static final class Target<X, V> implements BB<X, V> {
        private Duration duration;
        private Types<X> type;
        private V x;
    }
}
