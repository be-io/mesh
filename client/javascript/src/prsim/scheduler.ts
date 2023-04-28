/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import type {Context} from "@/prsim/context";
import {index} from "@/macro";
import {Status} from "@/cause";

@spi("mesh")
export abstract class Scheduler {

    /**
     * Timeout
     * Schedules the specified {@link Timeout} for one-time execution after the specified delay.
     */
    @mpi("mesh.schedule.timeout", String)
    timeout(ctx: Context, @index(0, 'timeout') timeout: string): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }


}