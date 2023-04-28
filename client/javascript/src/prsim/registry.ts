/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import type {Context} from "@/prsim/context";
import {Registration} from "@/kinds/registration";
import {Status} from "@/cause";
import {index} from "@/macro";

@spi("mesh")
export abstract class Registry {

    /**
     * Register
     */
    @mpi("mesh.registry.put")
    register(ctx: Context, @index(0, 'registration') registration: Registration): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }


    /**
     * Register
     */
    @mpi("mesh.registry.puts")
    registers(ctx: Context, @index(0, 'registrations') registrations: Registration[]): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }


    /**
     * Unregister
     */
    @mpi("mesh.registry.remove")
    unregister(ctx: Context, @index(0, 'registration') registration: Registration): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Export
     */
    @mpi("mesh.registry.export", [Array, Registration])
    export(ctx: Context, @index(0, 'kind') kind: string): Promise<Registration[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}