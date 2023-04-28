/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro";
import {Consumer} from "@/mpc/consumer";


@spi(Consumer.GRPC)
export abstract class Provider {

    /**
     * Start the mesh broker.
     */
    abstract start(): void;

    /**
     * Stop the mesh broker.
     */
    abstract close(): void;

}