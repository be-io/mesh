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
import java.util.ArrayList;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class WorkChart implements Serializable {

    private static final long serialVersionUID = -4401681747588197153L;
    /**
     * Workflow chart code
     */
    @Index(0)
    private String cno;
    /**
     * Workflow chart name.
     */
    @Index(1)
    private String name;
    /**
     * Workflow vertices
     */
    @Index(2)
    private List<WorkVertex> vertices;
    /**
     * Workflow sides
     */
    @Index(3)
    private List<WorkSide> sides;
    /**
     * Workflow status
     */
    @Index(4)
    private long status;
    /**
     * Workflow maintainer
     */
    @Index(5)
    private Worker maintainer;

    public void vertex(WorkVertex vertex) {
        if (null == this.vertices) {
            this.vertices = new ArrayList<>();
        }
        this.vertices.add(vertex);
    }

    public void link(WorkSide side) {
        if (null == this.sides) {
            this.sides = new ArrayList<>();
        }
        this.sides.add(side);
    }
}
