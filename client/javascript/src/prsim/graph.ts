/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import {Entity} from "@/kinds/entity";
import type {Context} from "@/prsim/context";

@spi("mesh")
export abstract class Graph {

    /**
     * Get the value from kv store.
     * @param key
     */
    @mpi("mesh.kv.get")
    get(ctx: Context, key: string): Entity {
        return new Entity()
    }
}