/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tools;


import com.be.mesh.client.mpc.URN;
import com.be.mesh.client.mpc.URNFlag;
import lombok.extern.slf4j.Slf4j;
import org.testng.Assert;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class URNTest {

    @Test
    public void urnConflictTest() {
        URNFlag flag = new URNFlag();
        flag.setV("00");
        flag.setProto("00");
        flag.setCodec("00");
        flag.setVersion("1.2.3");
        flag.setZone("00");
        flag.setCluster("00");
        flag.setCell("00");
        flag.setGroup("00");
        flag.setAddress("127.0.0.11");
        flag.setPort("4123");
        URN urn = new URN();
        urn.setNodeId("lx000000000000000000x");
        urn.setName("com.omega.network.edge.accessible");
        urn.setFlag(flag);
        log.info(urn.toString());
        Assert.assertEquals(urn.toString(), "accessible.edge.network.omega.com.0000000102030000000012700000001104123.lx000000000000000000x.trustbe.cn");
        flag.setAddress("xx.zz.test");
        log.info(urn.toString());
    }

    @Test
    public void urnLengthTest() {
        URNFlag flag = new URNFlag();
        flag.setV("00");
        flag.setProto("00");
        flag.setCodec("00");
        flag.setVersion("1.2.3");
        flag.setZone("00");
        flag.setCluster("00");
        flag.setCell("00");
        flag.setGroup("00");
        flag.setAddress("");
        flag.setPort("4123");
        URN urn = new URN();
        urn.setNodeId("lx000000000000000000x");
        urn.setName("com.omega.network.edge.accessible");
        urn.setFlag(flag);
        Assert.assertEquals(urn.getFlag().toString().length(), 37);
    }

    @Test
    public void urnParseTest() {
        log.info(Integer.valueOf("00001").toString());
        URN urn = URN.from("accessible.edge.network.omega.com.0000000102030000000012700000001104123.lx000000000000000000x.trustbe.cn");
        Assert.assertEquals(urn.getFlag().getAddress(), "127.0.0.11");
        Assert.assertEquals(urn.getFlag().getPort(), "4123");
        Assert.assertEquals(urn.getFlag().getVersion(), "1.2.3");
        Assert.assertEquals(urn.getName(), "com.omega.network.edge.accessible");
    }

    @Test
    public void urnParseBindingTest() {
        URN urn = URN.from("com.0000000102030000000012700000001104123.lx000000000000000000x.trustbe.cn");
        Assert.assertEquals(urn.getName(), "com");
    }
}
