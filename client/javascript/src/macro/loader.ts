/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {dict, Dict} from "@/macro/dsa";
import {SPI} from "@/macro/spi";
import {Ark, Class, Type} from "@/macro/ark";


class Anyone<T> {
    private static value: Dict<any, any> = new dict();
    public name: string;
    public kind: Class<T>;
    public prototype: boolean;
    public priority: number;

    constructor(name: string, kind: Class<T>) {
        this.name = name;
        this.kind = kind;
        this.prototype = false;
        this.priority = 1;
    }

    getIfAbsent(): T {
        if (this.prototype) {
            return this.create();
        }
        return Anyone.value.computeIfy(this.kind, key => this.create());
    }

    private create(): T {
        return new this.kind();
    }
}

class Instance<T> {
    public name: string = "";
    public providers: Anyone<T>[] = [];

    constructor(name: string, types: Class<T>[]) {
        this.name = name;
        this.providers = types.map(x => new Anyone<T>(name, x)).sort((p, n) => p.priority - n.priority);
    }

    getIfAbsent(): T {
        if (this.providers.length < 1) {
            throw new Error(`No such provider named ${this.name}`)
        }
        return this.providers[0].getIfAbsent()
    }
}

export class ServiceLoader<T> {

    private static loaders: Dict<Type<any>, ServiceLoader<any>> = new dict();

    public static load<T>(kind: Type<T>): ServiceLoader<T> {
        return this.loaders.computeIfy(kind, k => new ServiceLoader<any>(kind));
    }

    private providers: Dict<any, Dict<string, Instance<T>>> = new dict();
    private readonly spi: Type<T>;
    private readonly fist: string = "";

    constructor(spi: Type<T>) {
        this.spi = spi;
        this.fist = Ark.metadata(SPI, spi, spi)?.name || '';
    }

    defaultName(): string {
        return this.fist;
    }

    getDefault(): T {
        return this.get(this.fist);
    }

    get(name: string): T {
        const instance = this.getInstances().get(name);
        if (!instance) {
            throw new Error(`SPI ${this.spi.name} named ${name} not exist.`)
        }
        return instance.getIfAbsent()
    }

    list(): T[] {
        return this.getInstances().map((k, v) => v.getIfAbsent());
    }

    private getInstances(): Dict<string, Instance<T>> {
        return this.providers.computeIfy("$", k => {
            return Ark.providers(SPI, this.spi).groupBy(pc => {
                const spi = Ark.metadata(SPI, pc, pc);
                return (spi && spi.name) || ''
            }).groupKV(k => k, (name, pcs) => new Instance(name, pcs));
        })
    }
}