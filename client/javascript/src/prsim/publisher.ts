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
import {Event} from "@/kinds";
import {index} from "@/macro";

@spi("mesh")
export abstract class Publisher {

    /**
     * Publish
     */
    @mpi("mesh.queue.publish", [Array, String])
    publish(ctx: Context, @index(0, 'events') events: Event[]): Promise<string[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Broadcast
     */
    @mpi("mesh.queue.multicast", [Array, String])
    multicast(ctx: Context, @index(0, 'events') events: Event[]): Promise<string[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}