/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import type {Context} from "@/prsim/context";
import {Page, Paging} from "@/kinds/page";
import {Doc} from "@/kinds/doc";
import {Status} from "@/cause";
import {index} from "@/macro";

@spi("mesh")
export abstract class DataHouse {

    // Writes
    @mpi("mesh.dh.writes")
    writes(ctx: Context, @index(0, 'docs') docs: Doc): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Write
    @mpi("mesh.dh.write")
    write(ctx: Context, @index(0, 'doc') doc: Doc): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Read
    @mpi("mesh.dh.read", [Page, Object])
    read(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<any>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Indies
    @mpi("mesh.dh.indies", [Page, Object])
    indies(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<any>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Tables
    @mpi("mesh.dh.tables", [Page, Object])
    tables(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<any>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}