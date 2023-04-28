/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";
import {ServiceLoader} from "@/macro";
import {Codec} from "@/codec";

export class Entity {

    @index(0)
    public codec: string = "";
    @index(5)
    public schema: string = "";
    @index(10)
    public buffer: Uint8Array = new Uint8Array();

    public static empty(): Entity {
        return new Entity();
    }

    public static wrap(value: any): Entity {
        if (!value) {
            return Entity.empty();
        }
        const serviceLoader = ServiceLoader.load(Codec);
        const entity = new Entity();
        entity.codec = serviceLoader.defaultName();
        entity.schema = "";
        entity.buffer = serviceLoader.getDefault().encode(value);
        return entity;
    }

    public isPresent(): boolean {
        return null != this.buffer;
    }

    public readObject<T>(): T | null {
        if (!this.buffer || this.buffer.length < 1) {
            return null;
        }
        const codec = ServiceLoader.load(Codec).get(this.codec || Codec.JSON);
        // @ts-ignore
        return codec.decode(this.buffer)
    }

}


export const cacheVersion = "1.0.0";

export class CacheEntity {
    @index(0)
    public version: string = "";
    @index(5)
    public entity: Entity = new Entity();
    @index(10)
    public timestamp: number = 0;
    @index(15)
    public duration: number = 0;
    @index(20)
    public key: string = "";

    public wrap(key: string, value: any, duration: number): CacheEntity {
        const entity = Entity.wrap(value)
        const cn = new CacheEntity();
        cn.version = cacheVersion;
        cn.entity = entity;
        cn.timestamp = new Date().getTime();
        cn.duration = duration;
        cn.key = key;
        return cn;
    }

}