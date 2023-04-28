/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import {Ark} from "@/macro/ark";

/**
 * Metadata annotation for Serial Peripheral Interface. Can be used with {@link ServiceLoader#load(Class)}
 * or dependency injection at compile time and runtime time.
 */
export class SPI {

    public readonly name: string;
    public readonly pattern: string;
    public readonly priority: number;
    public readonly prototype: boolean;

    constructor(name: string, pattern: string, priority: number, prototype: boolean) {
        this.name = name
        this.pattern = pattern
        this.priority = priority
        this.prototype = prototype
    }

    decorate<TFunction extends Function>(target: TFunction): TFunction | void {
        Ark.annotate(SPI, target, target, this);
        Ark.register(SPI, this.name, target, this);
    }

}

export function spi(
    name = '',
    pattern = '',
    priority = 0,
    prototype = false): ClassDecorator {
    const meta = new SPI(name, pattern, priority, prototype);
    return meta.decorate.bind(meta);
}