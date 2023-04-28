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
 * <pre>
 * &#064;MPS
 * &#064;MPI("${mesh.name}.workflow.callback")
 * public class WorkCallback implements Endpoint.Sticker&lt;WorkContext, String&gt; {
 *
 *     &#064;Override
 *     public String stick(WorkContext varg) {
 *         return "";
 *     }
 * }
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@Data
public class WorkContext implements Serializable {


    private static final long serialVersionUID = 2630932336754186614L;
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
    @Index(2)
    private String cno;
    /**
     * Workflow context
     */
    @Index(3)
    private Map<String, String> context;
    @Index(4)
    private WorkVertex vertex;
    @Index(5)
    private WorkTask task;
    @Index(6)
    private Worker applier;
    @Index(7)
    private Worker reviewer;
}
