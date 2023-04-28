/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import type {Context} from "@/prsim/context";
import {Status} from "@/cause";

@spi("mesh")
export abstract class Dispatcher {

    /**
     * Invoke with map param
     * In multi returns, it's an array.
     */
    invoke(ctx: Context, urn: string, param: Map<string, any>): Promise<any[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Invoke0 with generic param
     * In multi returns, it's an array.
     */
    invoke0(ctx: Context, urn: string, param: any): Promise<any[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * InvokeLR with fewer returns
     * In multi returns, it will discard multi returns
     */
    invokeLR(ctx: Context, urn: string, param: Map<string, any>): Promise<any> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * InvokeLRG with fewer returns in generic mode
     * In multi returns, it will discard multi returns
     */
    invokeLRG(ctx: Context, urn: string, param: any): Promise<any> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }
}