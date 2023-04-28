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
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

/**
 * 32 bits.
 *
 * @author coyzeng@gmail.com
 */
@Data
@Slf4j
public class URNFlag implements Serializable {
    private static final long serialVersionUID = 2044138788126771935L;
    private String v = "00";                          // 2
    private String proto = "00";                      // 2
    private String codec = "00";                      // 2
    private String version = "000000";                // 6
    private String zone = "00";                       // 2
    private String cluster = "00";                    // 2
    private String cell = "00";                       // 2
    private String group = "00";                      // 2
    private String address = "000000000000";          // 12
    private String port = "00080";                    // 5   Max port 65535

    public static URNFlag from(String value) {
        if (Tool.optional(value)) {
            log.error("Unresolved flag {}", value);
            return new URNFlag();
        }
        URNFlag flag = new URNFlag();
        flag.setV(Tool.substring(value, 0, 2));
        flag.setProto(Tool.substring(value, 2, 2));
        flag.setCodec(Tool.substring(value, 4, 2));
        flag.setVersion(reduce(Tool.substring(value, 6, 6), 2));
        flag.setZone(Tool.substring(value, 12, 2));
        flag.setCluster(Tool.substring(value, 14, 2));
        flag.setCell(Tool.substring(value, 16, 2));
        flag.setGroup(Tool.substring(value, 18, 2));
        flag.setAddress(reduce(Tool.substring(value, 20, 12), 3));
        flag.setPort(reduce(Tool.substring(value, 32, 5), 5));
        return flag;
    }

    @Override
    public String toString() {
        return String.format("%s%s%s%s%s%s%s%s%s%s",
                padding(Tool.anyone(v, ""), 2),
                padding(Tool.anyone(proto, ""), 2),
                padding(Tool.anyone(codec, ""), 2),
                paddingChain(Tool.anyone(version, ""), 2, 3),
                padding(Tool.anyone(zone, ""), 2),
                padding(Tool.anyone(cluster, ""), 2),
                padding(Tool.anyone(cell, ""), 2),
                padding(Tool.anyone(group, ""), 2),
                paddingChain(Tool.anyone(address, ""), 3, 4),
                padding(Tool.anyone(port, ""), 5));
    }

    private String padding(String v, int length) {
        String value = v.replace(".", "");
        if (value.length() == length) {
            return value;
        }
        if (value.length() < length) {
            return Tool.repeat('0', length - v.length()) + value;
        }
        return v.substring(0, length);
    }

    private String paddingChain(String v, int length, int size) {
        String[] array = v.split("[.*]");
        List<String> frags = Arrays.stream(array).map(x -> Tool.isNumeric(x) ? x : "").collect(Collectors.toList());
        if (frags.size() == size) {
            return frags.stream().map(x -> padding(x, length)).collect(Collectors.joining());
        }
        if (frags.size() < size) {
            for (int index = 0; index < size - array.length; index++) {
                frags.add("");
            }
            return frags.stream().map(x -> padding(x, length)).collect(Collectors.joining());
        }
        return frags.subList(0, size).stream().map(x -> padding(x, length)).collect(Collectors.joining());
    }

    private static String reduce(String v, int length) {
        if (Tool.optional(v)) {
            return v;
        }
        StringBuilder bu = new StringBuilder();
        for (int index = 0; index < v.length(); index += length) {
            boolean hasNoneZero = true;
            for (int offset = index; offset < index + length; offset++) {
                if (v.charAt(offset) != '0') {
                    hasNoneZero = false;
                }
                if (!hasNoneZero) {
                    bu.append(v.charAt(offset));
                }
            }
            if (hasNoneZero) {
                bu.append('0');
            }
            bu.append('.');
        }
        return bu.substring(0, bu.length() - 1);
    }
}
