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
export abstract class Cryptor {

    // Encrypt binary to encrypted binary.
    @mpi("mesh.crypt.encrypt", Uint8Array)
    encrypt(ctx: Context, @index(0, 'buff') buff: Uint8Array, @index(1, 'features') features: Map<string, Uint8Array>): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Decrypt binary to decrypted binary.
    @mpi("mesh.crypt.decrypt", Uint8Array)
    decrypt(ctx: Context, @index(0, 'buff') buff: Uint8Array, @index(1, 'features') features: Map<string, Uint8Array>): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Hash compute the hash value.
    @mpi("mesh.crypt.hash", Uint8Array)
    hash(ctx: Context, @index(0, 'buff') buff: Uint8Array, @index(1, 'features') features: Map<string, Uint8Array>): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Sign compute the signature value.
    @mpi("mesh.crypt.sign", Uint8Array)
    sign(ctx: Context, @index(0, 'buff') buff: Uint8Array, @index(1, 'features') features: Map<string, Uint8Array>): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    // Verify the signature value.
    @mpi("mesh.crypt.verify", Boolean)
    verify(ctx: Context, @index(0, 'buff') buff: Uint8Array, @index(1, 'features') features: Map<string, Uint8Array>): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }
}