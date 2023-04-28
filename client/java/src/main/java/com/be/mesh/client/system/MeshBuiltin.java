/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPS;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.prsim.Builtin;
import com.be.mesh.client.prsim.Hodor;
import com.be.mesh.client.struct.Versions;
import com.be.mesh.client.tool.Tool;

import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@MPS
@SPI("mesh")
public class MeshBuiltin implements Builtin {

    @Override
    public String doc(String name, String formatter) {
        ByteBuffer buffer = ServiceLoader.resource(String.format("%s.schema", Tool.MESH_NAME.get()));
        Codec json = ServiceLoader.load(Codec.class).get(Codec.JSON);
        Codec codec = ServiceLoader.load(Codec.class).get(formatter);
        return codec.encodeString(json.decode(buffer, Types.MapObject));
    }

    @Override
    public Versions version() {
        ByteBuffer buffer = ServiceLoader.resource(String.format("%s.version", Tool.MESH_NAME.get()));
        Codec json = ServiceLoader.load(Codec.class).get(Codec.JSON);
        Versions versions = json.decode(buffer, Types.of(Versions.class));
        String name = Tool.toUpperCase(String.format("%s_version", Tool.MESH_NAME.get()));
        String mv = Tool.getProperty("", name);
        if (null != versions && Tool.required(mv)) {
            versions.setVersion(mv);
            Optional.ofNullable(versions.getInfos()).ifPresent(x -> x.put(name, mv));
        }
        return versions;
    }

    @Override
    public void debug(List<String> features) {
        List<Hodor> doors = ServiceLoader.load(Hodor.class).list();
        if (Tool.optional(doors)) {
            return;
        }
        for (Hodor hodor : doors) {
            hodor.debug(features);
        }
    }

    @Override
    public Map<String, String> stats(@Index(0) List<String> features) {
        List<Hodor> doors = ServiceLoader.load(Hodor.class).list();
        if (Tool.optional(doors)) {
            return new HashMap<>(0);
        }
        Map<String, String> indies = new HashMap<>();
        for (Hodor hodor : doors) {
            Map<String, String> stats = hodor.stats(features);
            if (Tool.required(stats)) {
                indies.putAll(stats);
            }
        }
        return indies;
    }

}
