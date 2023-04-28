/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.tool.Tool;
import lombok.Data;

import java.io.Serializable;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Versions implements Serializable, Comparable<Versions> {

    private static final long serialVersionUID = -7164749633684964423L;
    /**
     * Platform version.
     * <p>
     * main.sub.feature.bugfix like 1.5.0.0
     */
    @Index(0)
    private String version;
    /**
     * Network modules info.
     */
    @Index(5)
    private Map<String, String> infos;

    /**
     * Is the version less than the appointed version.
     *
     * @param v compared version
     * @return true less than
     */
    public int compare(String v) {
        if (Tool.optional(this.version)) {
            return -1;
        }
        if (Tool.optional(v)) {
            return 1;
        }
        String[] rvs = Tool.split(this.version, "\\.");
        String[] vs = v.split("\\.");
        for (int index = 0; index < rvs.length; index++) {
            if (vs.length <= index) {
                return 0;
            }
            if (Tool.equals(vs[index], "*")) {
                continue;
            }
            int c = rvs[index].compareTo(vs[index]);
            if (0 != c) {
                return c;
            }
        }
        return 0;
    }

    @Override
    public int compareTo(Versions ver) {
        if (null == ver) {
            return 1;
        }
        return this.compare(ver.version);
    }

    public Versions anyVersions(String sets) {
        Map<String, String> anyInfos = Optional.ofNullable(infos).orElseGet(Collections::emptyMap);
        String vk = String.format("%s.version", sets);
        String ak = String.format("%s.arch", sets);
        String ok = String.format("%s.os", sets);
        String ck = String.format("%s.commit_id", sets);
        Map<String, String> myInfos = new HashMap<>(4);
        myInfos.put(vk, Tool.anyone(anyInfos.get(vk), ""));
        myInfos.put(ak, Tool.anyone(anyInfos.get(ak), ""));
        myInfos.put(ok, Tool.anyone(anyInfos.get(ok), ""));
        myInfos.put(ck, Tool.anyone(anyInfos.get(ck), ""));
        Versions versions = new Versions();
        versions.setVersion(Tool.anyone(anyInfos.get(vk), ""));
        versions.setInfos(myInfos);
        return versions;
    }


}
