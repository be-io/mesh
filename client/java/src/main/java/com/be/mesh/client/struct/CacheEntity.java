/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class CacheEntity implements Serializable {

    private static final long serialVersionUID = -6797593694315549700L;
    @Index(0)
    private String version;
    @Index(5)
    private Entity entity;
    @Index(10)
    private long timestamp;
    @Index(15)
    private long duration;
    @Index(20)
    private String key;
}
