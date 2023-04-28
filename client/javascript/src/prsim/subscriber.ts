/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import type {Context} from "@/prsim/context";
import {Event} from "@/kinds/event";
import {mpi} from "@/macro/mpi";
import {spi} from "@/macro/spi";
import {index} from "@/macro";

@spi("mesh")
export abstract class Subscriber {

    /**
     * Subscribe the event with {@link com.be.mesh.client.annotate.Bindings} or {@link com.be.mesh.client.annotate.Binding}
     * @param ctx context
     * @param event event
     */
    @mpi("mesh.queue.subscribe")
    subscribe(ctx: Context, @index(0, 'event') event: Event): void {
        //
    }
}

@spi("mesh")
abstract class Listener {

    /**
     * Listen the event.
     * Listen function can't be blocked.
     *
     * @param ctx context
     * @param event event
     */
    listen(ctx: Context, @index(0, 'event') event: Event): void {
        //
    }

}
