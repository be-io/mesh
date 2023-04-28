/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import "reflect-metadata"
import {Ark, Type} from "@/macro/ark";

/**
 * Metadata annotation for Serial Peripheral Interface. Can be used with {@link ServiceLoader#load(Class)}
 * or dependency injection at compile time and runtime time.
 */
class Binding {
    private readonly topic: string;
    private readonly code: string;
    private readonly version: string;
    private readonly proto: string;
    private readonly codec: string;
    private readonly flags: number;
    private readonly timeout: number;
    private readonly meshable: boolean;
    private traits?: Type<any>;

    constructor(topic: string, code: string, version: string, proto: string, codec: string, flags: number,
                timeout: number, meshable: boolean) {
        this.topic = topic
        this.code = code
        this.version = version
        this.proto = proto
        this.codec = codec
        this.flags = flags
        this.timeout = timeout
        this.meshable = meshable
    }

    decorate<TFunction extends Function>(target: TFunction): TFunction | void {
        this.traits = Object.getPrototypeOf(target.prototype).constructor;
        Ark.annotate(Binding, target, target, this);
        Ark.register(Binding, `${this.topic}.${this.code}`, target, this);
    }
}

export function binding(
    topic = '',
    code = '',
    version = '1.0.0',
    proto = 'http2',
    codec = 'protobuf',
    flags = 0,
    timeout = 3000,
    meshable = true): ClassDecorator {
    const meta = new Binding(topic, code, version, proto, codec, flags, timeout, meshable);
    return meta.decorate.bind(meta);
}