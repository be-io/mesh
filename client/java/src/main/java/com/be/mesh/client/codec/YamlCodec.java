/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.codec;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.Types;
import org.yaml.snakeyaml.Yaml;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI(Codec.YAML)
public class YamlCodec implements Codec {

    private final Codec codec = new GsonCodec();

    @Override
    public ByteBuffer encode(Object value) {
        return this.encode0(value, object -> {
            Yaml yaml = new Yaml();
            return ByteBuffer.wrap(yaml.dump(object).getBytes(StandardCharsets.UTF_8));
        });
    }

    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        return this.decode0(buffer, type, (x, y) -> {
            Yaml yaml = new Yaml();
            Map<?, ?> map = yaml.loadAs(new String(x.array(), StandardCharsets.UTF_8), Map.class);
            return codec.decode(codec.encode(map), y);
        });
    }
}
