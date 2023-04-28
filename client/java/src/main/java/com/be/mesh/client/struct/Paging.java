/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import lombok.Data;

import java.io.Serializable;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Paging implements Serializable {

    private static final long serialVersionUID = -440196101885847262L;
    @Index(0)
    private String sid;
    @Index(5)
    private long index;
    @Index(10)
    private long limit;
    @Index(15)
    private Map<String, Object> factor;

    public <T> T tryReadFactor(Types<T> type) {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        return codec.decode(codec.encode(this.factor), type);
    }

}
