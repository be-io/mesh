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
public class CommerceEnviron implements Serializable {

    private static final long serialVersionUID = -1829967329680141914L;
    @Index(0)
    private String cipher;
    @Index(5)
    private Environ explain;
    @Index(value = 10, name = "node_key")
    private String nodeKey;
}
