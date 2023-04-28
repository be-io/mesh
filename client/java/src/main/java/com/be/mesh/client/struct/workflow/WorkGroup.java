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

/**
 * @author coyzeng@gmail.com
 */
@Data
public class WorkGroup implements Serializable {

    private static final long serialVersionUID = -3099050384540233480L;
    /**
     * Work group identity
     */
    @Index(0)
    private String no;
    /**
     * Work group name
     */
    @Index(1)
    private String name;
    /**
     * Workflow group status
     */
    @Index(2)
    private long status;
    /**
     * Work group workers
     */
    @Index(3)
    private List<Worker> workers;
}
