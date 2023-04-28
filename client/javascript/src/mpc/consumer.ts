/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Reference} from "@/kinds";
import {Context} from "@/prsim";
import {spi} from "@/macro";
import {Execution} from "@/mpc";


@spi(Consumer.HTTP)
export abstract class Consumer {

    public static HTTP: string = "http";
    public static GRPC: string = "grpc";
    public static TCP: string = "tcp";
    public static MQTT: string = "mqtt";

    /**
     * Start the mesh broker.
     */
    abstract start(): void;

    /**
     * Stop the mesh broker.
     */
    abstract close(): void;

    /**
     * Consume the input payload.
     *
     * @param ctx       Call context.
     * @param address   Remote address.
     * @param urn       Actual uniform resource domain name.
     * @param execution Service reference.
     * @param inbound   Input arguments.
     * @return Output payload
     */
    abstract consume(ctx: Context, address: string, urn: string, execution: Execution<Reference>, inbound: Uint8Array): Promise<Uint8Array>;
}