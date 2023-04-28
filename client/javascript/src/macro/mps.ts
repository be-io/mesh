/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import "reflect-metadata"
import {Ark, Type} from "@/macro/ark";

export class MPS {

    private readonly name: string;
    private readonly version: string;
    private readonly proto: string;
    private readonly codec: string;
    private readonly flags: number;
    private readonly timeout: number;
    private traits?: Type<any>;

    constructor(name: string, version: string, proto: string, codec: string, flags: number, timeout: number) {
        this.name = name;
        this.version = version;
        this.proto = proto;
        this.codec = codec;
        this.flags = flags;
        this.timeout = timeout;
    }

    decorate<TFunction extends Function>(target: TFunction): TFunction | void {
        this.traits = Object.getPrototypeOf(target.prototype).constructor;
        Ark.annotate(MPS, target, target, this);
        Ark.register(MPS, this.name, target, this);
    }
}

export function mps(
    name = '',
    version = '',
    proto = 'grpc',
    codec = 'json',
    flags = 0,
    timeout = 10000): ClassDecorator {
    const meta = new MPS(name, version, proto, codec, flags, timeout);
    return meta.decorate.bind(meta);
}