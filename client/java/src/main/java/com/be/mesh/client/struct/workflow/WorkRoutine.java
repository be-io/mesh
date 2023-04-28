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
public class WorkRoutine implements Serializable {

    private static final long serialVersionUID = -6466458428365373814L;
    /**
     * Workflow routine code
     */
    @Index(0)
    private String rno;
    /**
     * Business code
     */
    @Index(1)
    private String bno;
    /**
     * Workflow context
     */
    @Index(2)
    private Map<String, String> context;
    /**
     * Workflow status
     */
    @Index(3)
    private long status;
    /**
     * Workflow chart
     */
    @Index(4)
    private WorkChart chart;
    /**
     * Workflow tasks
     */
    @Index(5)
    private List<WorkTask> tasks;
}
