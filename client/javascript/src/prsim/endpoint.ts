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
export abstract class Endpoint {

    /**
     * Fuzzy call with generic param.
     */
    @mpi("${mesh.uname}", Uint8Array)
    fuzzy(ctx: Context, @index(0, 'buff') buff: Uint8Array): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}

/**
 * I is the input
 * O is the output
 */
interface EndpointSticker<I, O> {

    /**
     * Stick with generic param
     */
    stick(ctx: Context, varg: I): O;
}