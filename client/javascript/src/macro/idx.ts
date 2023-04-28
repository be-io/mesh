/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import "reflect-metadata"
import {Ark, Type} from "@/macro/ark";

export class Idx {

    public value: number;
    public name: string;
    public transparent: boolean;
    public kind: Type<any>;
    public index: number = 1000;

    constructor(value: number, name: string, transparent: boolean, kind: Type<any>) {
        this.value = value;
        this.name = name;
        this.transparent = transparent;
        this.kind = kind;
    }

    decorate(target: Object, propertyKey: string | symbol | undefined, parameterIndex?: number): void {
        if (this.kind == Object) {
            this.kind = Reflect.getMetadata('design:type', target, propertyKey as string | symbol);
        }
        this.index = parameterIndex ? parameterIndex : 1000;
        if (parameterIndex) {
            Ark.annotate(Idx, `${propertyKey as string}:${parameterIndex}`, target.constructor, this);
        } else {
            Ark.annotate(Idx, propertyKey, target.constructor, this);
        }
    }
}

export function index(
    value = -1,
    name = '',
    kind: Type<any> = Object,
    transparent = false,
): PropertyDecorator & ParameterDecorator {
    const meta = new Idx(value, name, transparent, kind);
    return meta.decorate.bind(meta);
}