/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Consumer, Execution} from "@/mpc";
import {spi} from "@/macro";
import {Context} from "@/prsim";
import {Reference} from "@/kinds";
import {Status} from "@/cause";

@spi(Consumer.GRPC)
export class GRPCConsumer extends Consumer {

    start(): void {

    }

    close(): void {
    }

    consume(ctx: Context, address: string, urn: string, execution: Execution<Reference>, inbound: Uint8Array): Promise<Uint8Array> {
        return Promise.reject(Status.URN_NOT_PERMIT);
    }

}