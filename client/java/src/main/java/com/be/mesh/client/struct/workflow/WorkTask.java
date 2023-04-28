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
import java.util.List;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class WorkTask implements Serializable {

    private static final long serialVersionUID = -5187050021488117211L;
    /**
     * Workflow vertex
     */
    @Index(0)
    private WorkVertex vertex;
    /**
     * Workflow vertex reviewers
     */
    @Index(1)
    private List<Worker> workers;
    /**
     * Workflow vertex status
     */
    @Index(2)
    private long status;
    /**
     * Workflow context
     */
    @Index(3)
    private Map<String, String> context;
}
