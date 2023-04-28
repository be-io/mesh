/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import type {Context} from "@/prsim/context";
import {Status} from "@/cause";
import {index} from "@/macro";

@spi("mesh")
export abstract class Cluster {

    /**
     * Election will election leader of instances.
     */
    @mpi("mesh.cluster.election", Uint8Array)
    election(ctx: Context, @index(0, 'buff') buff: Uint8Array): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * IsLeader if same level.
     */
    @mpi("mesh.cluster.leader", Boolean)
    leader(ctx: Context): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }
}