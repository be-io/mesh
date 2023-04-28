/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Workflow;
import io.be.mesh.struct.Page;
import io.be.mesh.struct.Paging;
import io.be.mesh.struct.workflow.*;
import lombok.extern.slf4j.Slf4j;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("mesh")
public class MeshWorkflow implements Workflow {

    private final Workflow workflow = ServiceProxy.proxy(Workflow.class);

    @Override
    public String mass(WorkGroup group) {
        return this.workflow.mass(group);
    }

    @Override
    public Page<WorkGroup> groups(Paging index) {
        return this.workflow.groups(index);
    }

    @Override
    public String compile(WorkChart chart) {
        return this.workflow.compile(chart);
    }

    @Override
    public Page<WorkChart> index(Paging index) {
        return this.workflow.index(index);
    }

    @Override
    public String submit(WorkIntent intent) {
        return this.workflow.submit(intent);
    }

    @Override
    public void take(WorkVertex vertex) {
        this.workflow.take(vertex);
    }

    @Override
    public Page<WorkRoutine> routines(Paging index) {
        return this.workflow.routines(index);
    }
}
