/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index, mpi, spi} from "@/macro";
import {Context} from "@/prsim/context";
import {Status} from "@/cause";
import {CommerceEnviron, CommerceLicense, License} from "@/kinds";

@spi("mesh")
export abstract class Commercialize {

    /**
     * Sign the license.
     */
    @mpi("mesh.license.sign", String)
    sign(ctx: Context, @index(0, 'lsr') lsr: License): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * History list the licenses.
     */
    @mpi("mesh.license.history", [Array, CommerceLicense])
    history(ctx: Context, @index(0, 'inst_id') inst_id: string): Promise<CommerceLicense[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Issued mesh node identity.
     */
    @mpi("mesh.net.issued", CommerceEnviron)
    issued(ctx: Context, @index(0, 'name') name: string, @index(1, 'kind') kind: string): Promise<CommerceEnviron> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Dump the node identity.
     */
    @mpi("mesh.net.dump", [Array, CommerceEnviron])
    dump(ctx: Context, @index(0, 'node_id') node_id: string): Promise<CommerceEnviron[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}