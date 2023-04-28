/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Invocation, Invoker} from "@/mpc/invoker";
import {Context} from "@/prsim";
import {Ark, ServiceLoader, SPI, spi} from "@/macro";

export const PROVIDER: string = "PROVIDER";
export const CONSUMER: string = "CONSUMER";

@spi("*")
export abstract class Filter {

    /**
     * Invoke the next filter.
     * @param ctx
     * @param invoker
     * @param invocation
     */
    abstract invoke(ctx: Context, invoker: Invoker<any>, invocation: Invocation): Promise<any>
}

/**
 * Composite the filter spi providers as a invoker.
 * @param invoker
 * @param pattern
 */
export function composite(invoker: Invoker<any>, pattern: string): Invoker<any> {
    const filters = ServiceLoader.load(Filter).list().filter(filter => {
        const spi = Ark.metadata(SPI, filter?.constructor, filter?.constructor);
        return !spi || spi.pattern == "" || spi?.pattern == pattern;
    }).sort((p, n): number => {
        const pp = Ark.metadata(SPI, p?.constructor, p?.constructor)?.priority || 0;
        const np = Ark.metadata(SPI, n?.constructor, n?.constructor)?.priority || 0;
        return pp == np ? 0 : (pp < np ? -1 : 1);
    });
    return composite0(invoker, filters);
}

function composite0(invoker: Invoker<any>, filters: Filter[]): Invoker<any> {
    let last = invoker;
    for (let index = filters.length - 1; index >= 0; index--) {
        let filter = filters[index];
        let next = last;
        last = new class implements Invoker<any> {
            run(ctx: Context, invocation: Invocation): Promise<any> {
                return filter.invoke(ctx, next, invocation);
            }
        }
    }
    return last;
}