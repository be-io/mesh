/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import lombok.Data;
import lombok.Getter;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@Getter
public class Features implements Serializable {

    private static final long serialVersionUID = 2962337535863738187L;
    private static final Map<String, List<String>> EXCLUDES = new ConcurrentHashMap<>();
    private final List<Feature> features = new ArrayList<>();

    public Features(String feature) {
        if (Tool.optional(feature)) {
            return;
        }
        String[] fss = Tool.split(feature, " ");
        for (String fs : fss) {
            String[] pair = Tool.split(fs, "=");
            if (pair.length < 2 || Tool.optional(pair[0]) || Tool.optional(pair[1])) {
                continue;
            }
            String[] groups = Tool.split(pair[0], "@");
            String[] names = Tool.split(groups[0], ":");
            if (names.length < 2) {
                continue;
            }
            String[] group = groups.length > 1 ? Tool.split(groups[1], ",") : new String[]{""};
            for (String u : group) {
                Feature f = new Feature();
                f.setName(names[0]);
                f.setKind(names[1]);
                f.setNodeId(u);
                f.setValue(pair[1]);
                features.add(f);
            }
        }
    }

    public static List<String> getActive(Class<?> spi) {
        return getSPIFeatures(spi, "active");
    }

    public static List<String> getInactive(Class<?> spi) {
        return getSPIFeatures(spi, "inactive");
    }

    public static List<String> getSPIFeatures(Class<?> spi, String name) {
        return Tool.MESH_FEATURE.get().getFeatures().stream().
                filter(x -> Tool.equals(spi.getCanonicalName(), x.getName())).
                filter(x -> Tool.equals(name, x.getKind())).
                map(x -> Tool.split(x.getValue(), ",")).
                flatMap(Arrays::stream).
                filter(Tool::required).
                collect(Collectors.toList());
    }

    @Data
    public static class Feature implements Serializable {

        private static final long serialVersionUID = 3120594176983889307L;
        private String name;
        private String nodeId;
        private String kind;
        private String value;
    }
}
