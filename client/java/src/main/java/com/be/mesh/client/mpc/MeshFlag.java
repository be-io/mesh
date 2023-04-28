/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.Arrays;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Getter
@AllArgsConstructor
public enum MeshFlag {

    HTTP("00", Consumer.HTTP),
    GRPC("01", Consumer.GRPC),
    MQTT("02", Consumer.MQTT),
    TCP("03", Consumer.TCP),

    JSON("00", Codec.JSON),
    PROTOBUF("01", Codec.PROTOBUF),
    XML("02", Codec.XML),
    THRIFT("03", Codec.THRIFT),
    YAML("04", Codec.YAML),

    ;

    public static final List<MeshFlag> PROTO = Arrays.asList(HTTP, GRPC, MQTT, TCP);
    public static final List<MeshFlag> CODEC = Arrays.asList(JSON, PROTOBUF, XML, THRIFT, YAML);
    private final String code;
    private final String name;

    public static MeshFlag ofProto(String code) {
        return PROTO.stream().filter(x -> x.code.equals(code)).findFirst().orElse(MeshFlag.HTTP);
    }

    public static MeshFlag ofCodec(String code) {
        return CODEC.stream().filter(x -> x.code.equals(code)).findFirst().orElse(MeshFlag.JSON);
    }

    public static MeshFlag ofName(String name) {
        return Arrays.stream(MeshFlag.values()).filter(x -> x.name.equals(name)).findFirst().orElse(MeshFlag.HTTP);
    }
}
