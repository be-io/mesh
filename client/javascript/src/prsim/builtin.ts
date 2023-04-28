/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import type {Context} from "@/prsim/context";
import {Versions} from "@/kinds/versions";
import {Status} from "@/cause";
import {index} from "@/macro";

@spi("mesh")
export abstract class Builtin {

    /**
     * Doc export the documents.
     */
    @mpi("${mesh.name}.builtin.doc", String)
    doc(ctx: Context, @index(0, 'name') name: string, @index(1, 'formatter') formatter: string): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Version will get the builtin application version.
     */
    @mpi("${mesh.name}.builtin.version", Versions)
    version(ctx: Context): Promise<Versions> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Debug set the application log level.
     */
    @mpi("${mesh.name}.builtin.debug")
    debug(ctx: Context, @index(0, 'features') features: Map<string, string>): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Stats will collect health check stats.
     */
    @mpi("${mesh.name}.builtin.stats", [Map, String, String])
    stats(ctx: Context, @index(0, 'features') features: string[]): Promise<Map<string, string>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Fallback is fallback service
     */
    @mpi("${mesh.name}.builtin.fallback")
    fallback(ctx: Context): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }
}