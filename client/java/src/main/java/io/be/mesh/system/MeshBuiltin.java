/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.Index;
import io.be.mesh.macro.MPS;
import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.Codec;
import io.be.mesh.mpc.ServiceLoader;
import io.be.mesh.mpc.Types;
import io.be.mesh.prsim.Builtin;
import io.be.mesh.prsim.Hodor;
import io.be.mesh.struct.Versions;
import io.be.mesh.tool.Tool;

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
