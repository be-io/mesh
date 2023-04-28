/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index, mpi, spi} from "@/macro";
import type {Context} from "@/prsim";
import {Status} from "@/cause";
import {Page, Paging, WorkChart, WorkGroup, WorkIntent, WorkRoutine, WorkVertex} from "@/kinds";

@spi("mesh")
export abstract class Workflow {

    /**
     * Mass workflow work in group.
     * Return workflow code
     */
    @mpi("mesh.workflow.mass", String)
    mass(ctx: Context, @index(0, 'group') group: WorkGroup): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Groups page workflow review groups.
     */
    @mpi("mesh.workflow.groups", [Page, WorkGroup])
    groups(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<WorkGroup>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Compile workflow in engine.
     * Return workflow code
     */
    @mpi("mesh.workflow.compile", String)
    compile(ctx: Context, @index(0, 'chart') chart: WorkChart): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Index workflows.
     */
    @mpi("mesh.workflow.index", [Page, WorkChart])
    index(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<WorkChart>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Submit workflow.
     * Return workflow instance code
     */
    @mpi("mesh.workflow.submit", String)
    submit(ctx: Context, @index(0, 'intent') intent: WorkIntent): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Take action on workflow instance.
     */
    @mpi("mesh.workflow.take")
    take(ctx: Context, @index(0, 'intent') intent: WorkVertex): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Routines infer workflow instance as routine.
     */
    @mpi("mesh.workflow.routines", [Page, WorkRoutine])
    routines(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<WorkRoutine>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}