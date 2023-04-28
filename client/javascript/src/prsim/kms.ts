/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import type {Context} from "@/prsim/context";
import {Environ} from "@/kinds";
import {Status} from "@/cause";
import {KeyCsr, Keys} from "@/kinds/keys";
import {index} from "@/macro";

@spi("mesh")
export abstract class KMS {

    /**
     * Environ will return the keystore environ.
     */
    @mpi("kms.store.environ", Environ)
    environ(ctx: Context): Promise<Environ> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * List will return the keystore environ.
     */
    @mpi("kms.crt.store.list", [Array, Keys])
    list(ctx: Context, @index(0, 'cno') cno: string): Promise<Keys[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * ApplyRoot will apply the root certification.
     */
    @mpi("kms.crt.apply.root", [Array, Keys])
    applyRoot(ctx: Context, @index(0, 'csr') csr: KeyCsr): Promise<Keys[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * ApplyIssue will apply the common certification.
     */
    @mpi("kms.crt.apply.issuøe", [Array, Keys])
    applyIssue(ctx: Context, @index(0, 'csr') csr: KeyCsr): Promise<Keys[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }


}