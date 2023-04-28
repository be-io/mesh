/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct.workflow;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Worker implements Serializable {

    private static final long serialVersionUID = -6101116313726773839L;
    /**
     * Worker identity
     */
    @Index(0)
    private String no;
    /**
     * Worker name
     */
    @Index(1)
    private String name;
    /**
     * Worker alias
     */
    @Index(2)
    private String alias;

}
