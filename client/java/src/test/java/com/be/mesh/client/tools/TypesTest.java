/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tools;

import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.struct.Versions;
import com.google.common.collect.ImmutableMap;
import lombok.extern.slf4j.Slf4j;
import org.junit.Assert;
import org.testng.annotations.Test;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class TypesTest {

    @Test
    public void typesTest() {
        Assert.assertSame(Types.of(String.class).getRawType(), String.class);
        Assert.assertTrue(Types.of(String.class).isAssignableFrom(String.class));
//        Assert.assertFalse(Types.list(String.class).isAssignableFrom(String.class));
//        Assert.assertTrue(Types.list(String.class).isAssignableFrom(List.class));
        Assert.assertTrue(Types.MapObject.isAssignableFrom(Map.class));
    }

    @Test
    public void versionTest() {
        Versions versions = new Versions();
        versions.setVersion("1.5.1.0");
        Assert.assertEquals(0, versions.compare("1.5.*"));
        Assert.assertEquals(-1, versions.compare("1.5.*.1"));
        Assert.assertEquals(1, versions.compare("1.4.*.1"));
    }

    @Test
    public void mapCodecTest() {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        log.info("{}", codec.decode(codec.encode(ImmutableMap.of("nodeId", "1")), Types.MapString));
    }
}
