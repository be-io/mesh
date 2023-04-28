/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import {CacheEntity} from "@/kinds/entity";
import type {Context} from "@/prsim/context";
import {Status} from "@/cause";
import {index} from "@/macro";

@spi("mesh")
export abstract class Cache {

    /**
     * Get the value from cache.
     */
    @mpi("mesh.cache.get", CacheEntity)
    get(ctx: Context, @index(0, 'key') key: string): Promise<CacheEntity> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Put the value to cache.
     */
    @mpi("mesh.cache.put")
    put(ctx: Context, @index(0, 'cell') cell: CacheEntity): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Remove the cache value.
     */
    @mpi("mesh.cache.remove")
    remove(ctx: Context, @index(0, 'key') key: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Incr the cache of expire time.
     */
    @mpi("mesh.cache.incr", Number)
    incr(ctx: Context, @index(0, 'key') key: string, @index(1, 'value') value: number): Promise<number> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Decr the cache of expire time.
     */
    @mpi("mesh.cache.decr", Number)
    decr(ctx: Context, @index(0, 'key') key: string, @index(1, 'value') value: number): Promise<number> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * HGet get value in hash
     */
    @mpi("mesh.cache.hget", CacheEntity)
    hget(ctx: Context, @index(0, 'key') key: string, @index(1, 'name') name: string): Promise<CacheEntity> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * HSet put value in hash
     */
    @mpi("mesh.cache.hset")
    hset(ctx: Context, @index(0, 'key') key: string, @index(1, 'cell') cell: CacheEntity): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * HDel put value in hash
     */
    @mpi("mesh.cache.hdel")
    hdel(ctx: Context, @index(0, 'key') key: string, @index(1, 'name') name: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * HKeys get the hash keys
     */
    @mpi("mesh.cache.hkeys", [Array, String])
    hkeys(ctx: Context, @index(0, 'key') key: string): Promise<string[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }
}