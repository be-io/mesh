/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.struct.Page;
import com.be.mesh.client.struct.Paging;
import com.be.mesh.client.struct.workflow.*;

/**
 * @author coyzeng@gmail.com
 */
public interface Workflow {

    /**
     * Mass workflow work in group.
     * Return workflow code
     */
    @MPI("mesh.workflow.mass")
    String mass(WorkGroup group);

    /**
     * Groups page workflow review groups.
     */
    @MPI("mesh.workflow.groups")
    Page<WorkGroup> groups(Paging index);

    /**
     * Compile workflow in engine.
     *
     * @return workflow code
     */
    @MPI("mesh.workflow.compile")
    String compile(WorkChart chart);

    /**
     * Index workflows.
     */
    @MPI("mesh.workflow.index")
    Page<WorkChart> index(Paging index);

    /**
     * Submit workflow.
     *
     * @return workflow instance code
     */
    @MPI("mesh.workflow.submit")
    String submit(@Index(0) WorkIntent intent);

    /**
     * Take action on workflow instance.
     */
    @MPI("mesh.workflow.take")
    void take(@Index(0) WorkVertex vertex);

    /**
     * Routines infer workflow instance as routine.
     */
    @MPI("mesh.workflow.routines")
    Page<WorkRoutine> routines(Paging index);
}
