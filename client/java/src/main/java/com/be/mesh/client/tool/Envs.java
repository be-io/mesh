/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.Arrays;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Envs implements Serializable {

    private static final long serialVersionUID = 7143192905849312569L;
    @Index(value = 0, alias = {"mesh.address", "mesh_address", "MESH_ADDRESS"})
    private String address;
    @Index(value = 5, alias = {"mesh.runtime", "mesh.runtime", "MESH_RUNTIME"})
    private String runtime;
    @Index(value = 10, alias = {"mesh.mode", "mesh_mode", "MESH_MODE"})
    private Mode mode;
    @Index(value = 15, alias = {"mesh.name", "mesh_name", "MESH_NAME", "spring.application.name"})
    private String name;
    @Index(value = 20, alias = {"mesh.direct", "mesh_direct", "MESH_DIRECT"})
    private String direct;
    @Index(value = 25, alias = {"mesh.subset", "mesh_subset", "MESH_SUBSET"})
    private String subset;
    @Index(value = 30, alias = {"mesh.zone", "mesh_zone", "MESH_ZONE"})
    private String zone;
    @Index(value = 35, alias = {"mesh.cluster", "mesh_cluster", "MESH_CLUSTER"})
    private String cluster;
    @Index(value = 40, alias = {"mesh.cell", "mesh_cell", "MESH_CELL"})
    private String cell;
    @Index(value = 45, alias = {"mesh.group", "mesh_group", "MESH_GROUP"})
    private String group;

    static final Once<String> MESH_NAME = Once.with(() -> Tool.getProperty("unknown", "MESH_NAME", "mesh_name", "mesh.name", "spring.application.name"));
    static final Once<Addrs> MESH_ADDRESS = Once.with(() -> new Addrs(Tool.getProperty(String.format("127.0.0.1:%d", Tool.DEFAULT_MESH_PORT), "MESH_ADDRESS", "mesh_address", "mesh.address")));
    static final Once<Mode> MESH_MODE = Once.with(() -> Mode.from(Tool.getProperty(Mode.FAILFAST.toString(), "MESH_MODE", "mesh_mode", "mesh.mode")));
    static final Once<String> MESH_DIRECT = Once.with(() -> Tool.getProperty("", "MESH_DIRECT", "mesh_direct", "mesh.direct"));
    static final Once<String> MESH_SUBSET = Once.with(() -> Tool.getProperty("", "MESH_SUBSET", "mesh_subset", "mesh.subset"));
    static final Once<Features> MESH_FEATURE = Once.with(() -> new Features(Tool.getProperty("", true, "MESH_FEATURE", "mesh_feature", "mesh.feature")));
    static final Once<String> MESH_RUNTIME_IP = Once.with(() -> Tool.getProperty("", "MESH_RUNTIME_IP", "mesh_runtime_ip", "_PAAS_NODE_NAME", "mesh.runtime.ip", "MESH.RUNTIME.IP"));
    static final Once<String> MESH_RUNTIME_PORT = Once.with(() -> Tool.getProperty("", "MESH_RUNTIME_PORT", "mesh_runtime_port", "_PAAS_PORT_7706", "mesh.runtime.port", "MESH.RUNTIME.PORT"));
    static final Once<String> MESH_ZONE = Once.with(() -> Tool.getProperty("", "MESH_ZONE", "mesh_zone", "mesh.zone"));
    static final Once<String> MESH_CLUSTER = Once.with(() -> Tool.getProperty("", "MESH_CLUSTER", "mesh_cluster", "mesh.cluster"));
    static final Once<String> MESH_CELL = Once.with(() -> Tool.getProperty("", "MESH_CELL", "mesh_cell", "mesh.cell"));
    static final Once<String> MESH_GROUP = Once.with(() -> Tool.getProperty("", "MESH_GROUP", "mesh_group", "mesh.group"));
    static final Once<String> IP_HEX = Once.with(() -> Arrays.stream(Tool.IP.get().split("\\.")).reduce("", (pre, cur) -> pre + Tool.padding(Tool.toUpperCase(Integer.toHexString(Integer.parseInt(cur))), 2, '0')));
    static final Once<Boolean> MESH_ENABLE = Once.with(() -> Tool.equals("true", Tool.getProperty("true", "mesh.enable", "mesh_enable", "MESH_ENABLE")));

    public static void release() {
        MESH_NAME.release();
        MESH_ADDRESS.release();
        MESH_MODE.release();
        MESH_DIRECT.release();
        MESH_SUBSET.release();
        MESH_FEATURE.release();
        MESH_RUNTIME_IP.release();
        MESH_RUNTIME_PORT.release();
        MESH_ZONE.release();
        MESH_CLUSTER.release();
        MESH_CELL.release();
        MESH_GROUP.release();
    }
}
