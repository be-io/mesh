/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import type {Context} from "./context";
import {mpi} from "@/macro/mpi";
import {spi} from "@/macro/spi";
import {index} from "@/macro/idx";
import {Status} from "@/cause";

@spi("mesh")
export abstract class Tokenizer {

    /**
     * Apply a node token.
     */
    @mpi("mesh.trust.apply", String)
    apply(ctx: Context, @index(0, 'kind') kind: string, @index(1, 'duration') duration: number): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Verify some token verifiable.
     */
    @mpi("mesh.trust.verify", Boolean)
    verify(ctx: Context, @index(0, 'token') token: string): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}
