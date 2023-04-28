/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.tools;

import io.be.mesh.cause.MeshException;
import io.be.mesh.cause.TimeoutException;
import io.be.mesh.grpc.GrpcContextKey;
import io.be.mesh.mpc.Filter;
import io.be.mesh.mpc.Mesh;
import io.be.mesh.mpc.ServiceLoader;
import io.be.mesh.tool.Addrs;
import io.be.mesh.tool.Tool;
import io.grpc.Metadata;
import lombok.extern.slf4j.Slf4j;
import org.junit.Assert;
import org.testng.annotations.Test;

import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class ToolTest {

    @Test
    public void metadataSetTest() {
        Mesh.contextSafeUncheck(() -> {
            Mesh.context().getAttachments().put(GrpcContextKey.SUBSET.toString(), null);
            Metadata metadata = new Metadata();
            metadata.put(GrpcContextKey.SUBSET, Mesh.context().getAttachments().getOrDefault(GrpcContextKey.SUBSET.toString(), ""));
        });
    }

    @Test
    public void clusterAddressTest() {
        Assert.assertEquals("", new Addrs("").any());
        Assert.assertEquals("1", new Addrs("1").any());
        Assert.assertNotEquals("0", new Addrs("1,2,3").any());
        log.info(new Addrs("1,2,3").any());
    }

    @Test
    public void traceIdTest() {
        log.info(Tool.newTraceId());
        log.info(Tool.IP_HEX.get());
    }

    @Test
    public void optionalTest() {
        Assert.assertTrue(Tool.optional());
        Assert.assertTrue(Tool.optional("", ""));
        Assert.assertTrue(Tool.optional("", "1"));
        Assert.assertEquals("", "1".substring(1));
    }

    @Test
    public void addressTest() {
        System.setProperty("mesh.address", "gaia-mesh");
        String[] hosts = Tool.split(Tool.MESH_ADDRESS.get().any(), ":");
        log.info("{}", (Object) hosts);
        if (hosts.length < 2) {
            Assert.assertTrue(true);
            return;
        }
        Assert.fail();
    }

    @Test
    public void parseLongTest() {
        log.info("{}", Long.parseLong("127000000001"));
        log.info("{}", Long.parseLong("000000000001"));
        log.info("{}", Long.parseLong("000000000000"));
    }

    @Test
    public void numericTest() {
        Assert.assertTrue(Tool.isNumeric("19999999"));
        Assert.assertTrue(Tool.isNumeric("1"));
        Assert.assertTrue(Tool.isNumeric("9"));
        Assert.assertTrue(Tool.isNumeric("0"));
        Assert.assertFalse(Tool.isNumeric("0.1"));
        Assert.assertFalse(Tool.isNumeric(""));
        Assert.assertFalse(Tool.isNumeric(" "));
        Assert.assertFalse(Tool.isNumeric("a"));
    }

    @Test
    public void startWithTest() {
        Assert.assertTrue(Tool.startWith("19999999", "19999999"));
    }

    @Test
    public void instanceOfTest() {
        Assert.assertTrue(TimeoutException.class.isAssignableFrom(TimeoutException.class));
        Assert.assertTrue(MeshException.class.isAssignableFrom(TimeoutException.class));
        Assert.assertFalse(TimeoutException.class.isAssignableFrom(MeshException.class));
    }

    @Test
    public void featuresTest() {
        System.setProperty("MESH_FEATURE", String.format("%s:inactive=session,exclude-test-session", Filter.class.getCanonicalName()));
        List<String> excludes = Tool.MESH_FEATURE.get().getFeatures().stream().
                filter(x -> Tool.equals(Filter.class.getCanonicalName(), x.getName())).
                filter(x -> Tool.equals("inactive", x.getKind())).
                map(x -> Tool.split(x.getValue(), ",")).
                flatMap(Arrays::stream).
                filter(Tool::required).collect(Collectors.toList());
        log.info("{}", excludes);

        log.info("{}", ServiceLoader.load(Filter.class).list().size());
    }
}
