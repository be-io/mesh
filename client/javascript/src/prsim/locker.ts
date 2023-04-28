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
export abstract class Locker {

    /**
     * Lock create write lock.
     */
    @mpi("mesh.locker.w.lock", Boolean)
    lock(ctx: Context, @index(0, 'rid') rid: string, @index(1, 'timeout') timeout: number): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Unlock release write lock.
     */
    @mpi("mesh.locker.w.unlock")
    unlock(ctx: Context, @index(0, 'rid') rid: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * ReadLock create read lock.
     */
    @mpi("mesh.locker.r.lock", Boolean)
    readLock(ctx: Context, @index(0, 'rid') rid: string, @index(1, 'timeout') timeout: number): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * ReadUnlock release read lock.
     */
    @mpi("mesh.locker.r.unlock")
    readUnlock(ctx: Context, @index(0, 'rid') rid: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }
}