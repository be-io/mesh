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
import {License} from "@/kinds/license";
import {index} from "@/macro";

@spi("mesh")
export abstract class Licensor {

    /**
     * Imports the licenses.
     */
    @mpi("mesh.license.imports", String)
    imports(ctx: Context, @index(0, 'license') license: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Exports the licenses.
     */
    @mpi("mesh.license.exports", String)
    exports(ctx: Context): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Explain the license.
     */
    @mpi("mesh.license.explain", License)
    explain(ctx: Context): Promise<License> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Verify the license.
     */
    @mpi("mesh.license.verify", Number)
    verify(ctx: Context): Promise<number> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Verify the license.
     */
    @mpi("mesh.license.features", [Map, String, String])
    features(ctx: Context): Promise<Map<string, string>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}