/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro";

export enum VertexKind {
    WorkStart = 1,
    WorkFinish = 2,
    WorkJob = 4,
    WorkMan = 8,
    WorkWay = 16,
    WorkTimer = 32,
}

export class Worker {
    /**
     * Workflow worker code
     */
    @index(0)
    public no: string = "";
    /**
     * Workflow worker name
     */
    @index(1)
    public name: string = "";
    /**
     * Workflow worker alias
     */
    @index(2)
    public alias: string = "";
}

export class WorkIntent {
    /**
     * Business code
     */
    @index(0)
    public bno: string = "";
    /**
     * Workflow chart code
     */
    @index(1)
    public cno: string = "";
    /**
     * Workflow context
     */
    @index(2)
    public context: Map<string, string> = new Map();
    /**
     * Workflow applier
     */
    @index(3)
    public applier: Worker = new Worker();
}

export class WorkVertex {
    /**
     * Workflow name
     */
    @index(0)
    public name: string = "";
    /**
     * Workflow alias
     */
    @index(1)
    public alias: string = "";
    /**
     * Workflow context
     */
    @index(2)
    public attrs: Map<string, string> = new Map();
    /**
     * Workflow kind
     */
    @index(3)
    public kind: number = 0;
    /**
     * Workflow review group
     */
    @index(4)
    public group: string = "";
}

export class WorkSide {
    /**
     * Workflow side src name
     */
    @index(0)
    public src: string = "";
    /**
     * Workflow side dst name
     */
    @index(1)
    public dst: string = "";
    /**
     * Workflow side condition
     */
    @index(2)
    public condition: string = "";
}

export class WorkChart {
    /**
     * Workflow chart code
     */
    @index(0)
    public cno: string = "";
    /**
     * Workflow chart name
     */
    @index(1)
    public name: string = "";
    /**
     * Workflow vertices
     */
    @index(2)
    public vertices: WorkVertex[] = [];
    /**
     * Workflow sides
     */
    @index(3)
    public sides: WorkSide[] = [];
    /**
     * Workflow status
     */
    @index(3)
    public status: number = 0;
    /**
     * Workflow maintainer
     */
    @index(4)
    public maintainer: Worker = new Worker();

    public vertex(vertex: WorkVertex): void {
        this.vertices.push(vertex);
    }

    public link(side: WorkSide): void {
        this.sides.push(side);
    }

}


export class WorkRoutine {
    /**
     * Workflow routine code
     */
    @index(0)
    public rno: string = "";
    /**
     * Workflow business code
     */
    @index(1)
    public bno: string = "";
    /**
     * Workflow context
     */
    @index(2)
    public context: Map<string, string> = new Map();
    /**
     * Workflow status
     */
    @index(3)
    public status: number = 0;
    /**
     * Workflow chart
     */
    @index(4)
    public chart: WorkChart = new WorkChart();
    /**
     * Workflow tasks
     */
    @index(5)
    public tasks: WorkChart[] = [];
}

export class WorkTask {
    /**
     * Workflow vertex
     */
    @index(0)
    public vertex: WorkVertex = new WorkVertex();
    /**
     * Workflow reviewers
     */
    @index(1)
    public reviewers: Worker[] = [];
    /**
     * Workflow status
     */
    @index(2)
    public status: number = 0;
    /**
     * Workflow context
     */
    @index(3)
    public context: Map<string, string> = new Map();
}

export class WorkGroup {
    /**
     * Workflow group code
     */
    @index(0)
    public no: string = "";
    /**
     * Workflow group name
     */
    @index(1)
    public name: string = "";
    /**
     * Workflow status
     */
    @index(2)
    public status: number = 0;
    /**
     * Workflow reviewers
     */
    @index(3)
    public reviewers: Worker[] = [];
}

export class WorkContext {
    /**
     * Workflow routine code
     */
    @index(0)
    public rno: string = "";
    /**
     * Workflow business code
     */
    @index(1)
    public bno: string = "";
    /**
     * Workflow chart code
     */
    @index(2)
    public cno: string = "";
    /**
     * Workflow context
     */
    @index(3)
    public context: Map<string, string> = new Map();
    /**
     * Workflow vertex
     */
    @index(4)
    public vertex: WorkVertex = new WorkVertex();
    /**
     * Workflow task
     */
    @index(5)
    public task: WorkTask = new WorkTask();
    /**
     * Workflow applier
     */
    @index(6)
    public applier: Worker = new Worker();
    /**
     * Workflow reviewer
     */
    @index(7)
    public reviewer: Worker = new Worker();
}
