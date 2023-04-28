/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.tool.Tool;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Like:
 * <pre>
 *     create.tenant.omega . 0001 . json . http2 . lx000001 . trustbe.cn
 *     -------------------   ----   ----   -----   --------   -----------
 *           name            flag   codec  proto     node       domain
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@Data
@Slf4j
public class URN implements Serializable {

    private static final long serialVersionUID = -1427045979049617090L;
    public static final String MESH_DOMAIN = "trustbe.cn";
    private String domain;
    private String nodeId;
    private URNFlag flag;
    private String name;

    public static URN from(String urn) {
        URN name = new URN();
        if (Tool.optional(urn)) {
            log.warn("Unresolved urn {}", urn);
            return name;
        }
        List<String> names = name.asArray(urn);
        if (names.size() < 5) {
            log.warn("Unresolved urn {}", urn);
            return name;
        }
        Collections.reverse(names);
        name.domain = String.format("%s.%s", names.remove(1), names.remove(0)).intern();
        name.nodeId = Tool.toLowerCase(names.remove(0));
        name.flag = URNFlag.from(names.remove(0));
        name.name = String.join(".", names);
        return name;
    }

    public static URN local(String name) {
        return any(name, Tool.LOCAL_INST_ID);
    }

    public static URN any(String name, String nodeId) {
        URN urn = new URN();
        urn.setDomain(URN.MESH_DOMAIN);
        urn.setNodeId(nodeId);
        urn.setFlag(URNFlag.from("0001000000000000000000000000000000000"));
        urn.setName(name);
        return urn;
    }

    public static URN any(MeshFlag flag, String name, String nodeId) {
        URNFlag f = URNFlag.from("0001000000000000000000000000000000000");
        if (MeshFlag.PROTO.contains(flag)) {
            f.setProto(flag.getCode());
        }
        if (MeshFlag.CODEC.contains(flag)) {
            f.setCodec(flag.getCode());
        }
        URN urn = new URN();
        urn.setDomain(URN.MESH_DOMAIN);
        urn.setNodeId(nodeId);
        urn.setFlag(f);
        urn.setName(name);
        return urn;
    }

    @Override
    public String toString() {
        List<String> urn = new ArrayList<>();
        List<String> names = asArray(name);
        Collections.reverse(names);
        names.stream().filter(Tool::required).forEach(urn::add);
        urn.add(flag.toString());
        urn.add(Tool.toLowerCase(nodeId));
        urn.add(Tool.anyone(domain, MESH_DOMAIN));
        return String.join(".", urn);
    }

    private List<String> asArray(String text) {
        String[] pairs = text.split("\\.");
        List<String> names = new ArrayList<>();
        StringBuilder buffer = new StringBuilder();
        for (String pair : pairs) {
            if (Tool.startWith(pair, "${")) {
                buffer.append(pair).append('.');
                continue;
            }
            if (Tool.endsWith(pair, "}")) {
                buffer.append(pair);
                names.add(buffer.toString());
                buffer.delete(0, buffer.length());
                continue;
            }
            names.add(pair);
        }
        return names;
    }

    public boolean matchNode(String nodeId) {
        return Tool.equals(Tool.toLowerCase(this.nodeId), Tool.toLowerCase(nodeId));
    }

    public boolean matchName(String name) {
        return Tool.endsWith(name, ".*") && Tool.startWith(this.name, name.replace(".*", ""));
    }

}
