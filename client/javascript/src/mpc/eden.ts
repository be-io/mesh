/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Inspector, spi, Type} from "@/macro";
import {Context} from "@/prsim";
import {Reference, Service} from "@/kinds";
import {Execution} from "@/mpc/execution";

@spi("mesh")
export abstract class Eden {

    /**
     * Define the reference object.
     * @param ctx
     * @param metadata
     * @param reference
     */
    abstract define(ctx: Context, metadata: any, reference: Type<any>): any;

    /**
     * Refer the service reference by method.
     * @param ctx
     * @param metadata
     * @param reference
     * @param method
     */
    abstract refer(ctx: Context, metadata: any, reference: Type<any>, method: Inspector): Execution<Reference>;

    /**
     * Store the service object.
     * @param ctx
     * @param kind
     * @param service
     */
    abstract store(ctx: Context, kind: Type<any>, service: any): void;

    /**
     * Infer the reference service by domain.
     * @param ctx
     * @param urn
     */
    abstract infer(ctx: Context, urn: string): Execution<Service>;

    /**
     * Get all reference types.
     * @param ctx
     */
    abstract referTypes(ctx: Context): Type<any>[];

    /**
     * Get all service types.
     * @param ctx
     */
    abstract inferTypes(ctx: Context): Type<any>[];
}