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
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class WorkVertex implements Serializable {

    private static final long serialVersionUID = 68232640882741144L;
    /**
     * Workflow name
     */
    @Index(0)
    private String name;
    /**
     * Workflow alias
     */
    @Index(1)
    private String alias;
    /**
     * Workflow vertex attributes
     */
    @Index(2)
    private Map<String, String> attrs;
    /**
     * Workflow vertex kind
     */
    @Index(3)
    private long kind;
    /**
     * Workflow review group code
     */
    @Index(4)
    private String group;
}
