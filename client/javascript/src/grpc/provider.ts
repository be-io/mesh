/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Consumer, Provider} from "@/mpc";
import {spi} from "@/macro";

@spi(Consumer.GRPC)
export class GRPCProvider extends Provider {

    start(): void {
    }

    close(): void {
    }

}