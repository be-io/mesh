/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import {Entity} from "@/kinds/entity";
import type {Context} from "@/prsim/context";
import {Status} from "@/cause";
import {index} from "@/macro";

@spi("mesh")
export abstract class KV {

    /**
     * Get the value from kv store.
     */
    @mpi("mesh.kv.get", Entity)
    get(ctx: Context, @index(0, 'key') key: string): Promise<Entity> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Put the value to kv store.
     */
    @mpi("mesh.kv.put")
    put(ctx: Context, @index(0, 'key') key: string, @index(1, 'value') value: Entity): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Remove the kv store.
     */
    @mpi("mesh.kv.remove")
    remove(ctx: Context, @index(0, 'key') key: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Inspect keys of the kv store.
     */
    @mpi("mesh.kv.keys", [Array, String])
    keys(ctx: Context, @index(0, 'pattern') pattern: string): Promise<string[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}