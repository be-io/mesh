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
import {Page, Paging, Script} from "@/kinds";
import {index} from "@/macro";


@spi("mesh")
export abstract class Evaluator {

    /**
     * Compile the named rule.
     */
    @mpi("mesh.eval.compile", String)
    compile(ctx: Context, @index(0, 'script') script: Script): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Exec the script with name.
     */
    @mpi("mesh.eval.exec", String)
    exec(ctx: Context, @index(0, 'code') code: string, @index(1, 'args') args: Map<string, string>, @index(2, 'dft') dft: string): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Dump the scripts.
     */
    @mpi("mesh.eval.dump", [Array, Script])
    dump(ctx: Context, @index(0, 'feature') feature: Map<string, string>): Promise<Script[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Index the scripts.
     */
    @mpi("mesh.eval.index", [Page, Script])
    index(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<Script>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}