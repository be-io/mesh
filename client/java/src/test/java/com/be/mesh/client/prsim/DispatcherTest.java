/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.*;
import com.be.mesh.client.struct.Route;
import com.google.common.collect.ImmutableMap;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.util.HashMap;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class DispatcherTest {

    @Test
    public void networkTest() {
        System.setProperty("mesh.address", "10.99.2.33");
        Network network = ServiceLoader.load(Network.class).getDefault();
        URNFlag flag = new URNFlag();
        flag.setProto(MeshFlag.GRPC.getCode());
        flag.setCodec(MeshFlag.JSON.getCode());
        Dispatcher dispatcher = ServiceLoader.load(Dispatcher.class).getDefault();
        URN urn = new URN();
        urn.setDomain(URN.MESH_DOMAIN);
        urn.setNodeId(network.getEnviron().getNodeId());
        urn.setFlag(flag);
        urn.setName("mesh.net.environ");

        log.info("{}", dispatcher.invoke(urn.toString(), new HashMap<>()));

        urn.setName("mesh.net.edges");
        log.info("{}", dispatcher.invoke(urn.toString(), new HashMap<>()));

        urn.setName("mesh.net.accessible");
        log.info("{}", dispatcher.invoke(urn.toString(), ImmutableMap.of("route", new Route())));

        urn.setName("mesh.net.version");
        log.info("{}", dispatcher.invoke(urn.toString(), ImmutableMap.of("nodeId", "LX0000010000030")));
    }

    @Test
    private void edgeTest() {
        String input = "{\"attachments\":{\"mesh_span_id\":\"0.170.1\",\"mesh_run_mode\":\"1\",\"mesh_timestamp\":\"1646204932606\",\"omega.tenant.id\":\"JG0100000100000000\",\"mesh_urn\":\"hostModelDeployCallback.model.serving.edge.0000000000000000000000000000000000000.jg0100000500000000.trustbe.cn\",\"mesh_consumer\":\"{}\",\"mesh_provider\":\"{\\\"ip\\\":\\\"10.90.40.253\\\",\\\"port\\\":\\\"7706\\\",\\\"host\\\":\\\"edge-serving-78ccc74cd9-tvf46\\\",\\\"node_id\\\":\\\"LX0000010000010\\\",\\\"inst_id\\\":\\\"JG0100000100000000\\\"}\",\"mesh_trace_id\":\"a5a28fd1498918244106014720\",\"omega.inst.id\":\"JG0100000500000000\"},\"request\":{\"bizId\":\"573bcadf5fdf4220bbbca56b297c7907\",\"sourcePartyId\":\"JG0100000100000000\",\"status\":\"SUCCESS\",\"partyId\":\"JG0100000500000000\"}}";
        URNFlag flag = new URNFlag();
        flag.setProto(MeshFlag.GRPC.getCode());
        flag.setCodec(MeshFlag.JSON.getCode());
        Dispatcher dispatcher = ServiceLoader.load(Dispatcher.class).getDefault();
        URN urn = new URN();
        urn.setDomain(URN.MESH_DOMAIN);
        urn.setNodeId("JG0000000000000000");
        urn.setFlag(flag);
        urn.setName("edge.serving.model.hostModelDeployCallback");
        log.info("{}", dispatcher.invoke(urn.toString(), ServiceLoader.load(Codec.class).getDefault().decodeString(input, Types.MapObject)));
    }
}
