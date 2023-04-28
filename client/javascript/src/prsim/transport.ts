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

/**
 * Remote queue in async and blocking mode.
 */
@spi("mesh")
export abstract class Session {

    /**
     * Retrieves, but does not remove, the head of this queue, or returns None if this queue is empty.
     */
    @mpi("mesh.chan.peek", Uint8Array)
    peek(ctx: Context, @index(0, 'topic') topic: string): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Retrieves and removes the head of this queue, or returns None if this queue is empty.
     */
    @mpi("mesh.chan.pop", Uint8Array, "1.0.0", "grpc", "json", 0, 120 * 1000)
    pop(ctx: Context, @index(0, 'timeout') timeout: number, @index(1, 'topic') topic: string): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Inserts the specified element into this queue if it is possible to do
     *         so immediately without violating capacity restrictions.
     *         When using a capacity-restricted queue, this method is generally
     *         preferable to add, which can fail to insert an element only
     *         by throwing an exception.
     */
    @mpi("mesh.chan.push")
    push(ctx: Context, @index(0, 'payload') payload: Uint8Array, @index(1, 'metadata') metadata: Map<string, string>, @index(2, 'topic') topic: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Close the channel session.
     */
    @mpi("mesh.chan.release")
    release(ctx: Context, @index(0, 'timeout') timeout: number, @index(1, 'topic') topic: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }


}

/**
 * Private compute data channel in async and blocking mode.
 */
@spi("mesh")
export abstract class Transport {

    /**
     * Open a channel session.
     */
    @mpi("mesh.chan.open", Session) open(ctx: Context, @index(0, 'session_id') session_id: string, @index(1, 'metadata') metadata: Map<string, string>): Promise<Session> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Close the channel.
     */
    @mpi("mesh.chan.close")
    close(ctx: Context, @index(0, 'timeout') timeout: number): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Roundtrip with the channel.
     */
    @mpi("mesh.chan.roundtrip", Uint8Array)
    roundtrip(ctx: Context, @index(0, 'payload') payload: Uint8Array, @index(1, 'metadata') metadata: Map<string, string>): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }
}
