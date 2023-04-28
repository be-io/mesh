/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Ark, Idx, spi, Type, Types} from "@/macro";
import {Codec} from "@/codec/codec";
import {Buffer} from "buffer";


@spi(Codec.JSON)
export class JSONCodec extends Codec {

    private encoder: TextEncoder = new TextEncoder();
    private decoder: TextDecoder = new TextDecoder();

    decode<T>(buffer: Uint8Array, type: Type<T>): T {
        const x = this.decoder.decode(buffer);
        if (!type || !type.constructor) {
            return JSON.parse(x);
        }
        const dict = JSON.parse(x);
        if ((typeof dict !== 'object' && typeof dict !== 'function') || dict == null) {
            return dict;
        }
        switch (type.name) {
            case "Object":
                return dict;
        }
        return this.decorateAny(type, dict);
    }

    decorateAny(type: Type<any>, dict: any): any {
        const generics = type instanceof Types;
        const kind = generics ? type.raw : type;
        switch (kind.name) {
            case "Map":
                if (null == dict) {
                    return new Map();
                }
                if (!generics || type.parameters.length < 1) {
                    return new Map(Object.entries(dict));
                }
                const map = new Map();
                Object.keys(dict).forEach((key) => {
                    map.set(this.decorateAny(type.parameters[0], key), this.decorateAny(type.parameters[1], dict[key]))
                });
                return map;
            case "Array":
                if (null == dict) {
                    return [];
                }
                if (!generics || type.parameters.length < 1) {
                    return dict;
                }
                return dict.map((v: any) => this.decorateAny(type.parameters[0], v));
            case "Page":
                if (null == dict) {
                    return null;
                }
                // @ts-ignore
                const page = new kind();
                if (!generics || type.parameters.length < 1) {
                    return Object.assign(page, dict);
                }
                Object.assign(page, dict);
                page.data = this.decorateAny(type.parameters[0], page.data);
                return page;
            case "String":
            case "Number":
            case "Boolean":
            case "Object":
                return dict;
            case "Date":
                if (null == dict) {
                    return null;
                }
                return new Date(dict);
            case "Uint8Array":
                if (null == dict) {
                    return null;
                }
                if (dict instanceof Uint8Array) {
                    return dict;
                } else {
                    return Buffer.from(dict as string, 'base64');
                }
            default:
                if (null == dict) {
                    return null;
                }
                // @ts-ignore
                const instance = new kind();
                Object.keys(instance).forEach((key, index) => {
                    const value = dict[key];
                    const metadata = Ark.metadata(Idx, kind, key);
                    if (!metadata) {
                        Reflect.set(instance, key, value);
                        return
                    }
                    Reflect.set(instance, key, this.decorateAny(metadata.kind, value));
                });
                return instance;
        }
    }

    encode(value: any): Uint8Array {
        return this.encoder.encode(JSON.stringify(value, (key: string, value: any): any => {
            if (value instanceof Uint8Array) {
                return Buffer.from(value).toString('base64');
            }
            if (value instanceof Date) {
                return value.getTime();
            }
            if (value instanceof Map) {
                return Object.fromEntries(value);
            }
            return value;
        }));
    }

    uint8ify<T>(buffer: string): Uint8Array {
        return this.encoder.encode(buffer);
    }

    stringify<T>(buffer: Uint8Array): string {
        return this.decoder.decode(buffer);
    }

}