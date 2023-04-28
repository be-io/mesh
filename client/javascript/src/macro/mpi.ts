/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import "reflect-metadata"
import {Ark} from "@/macro/ark";


// 'design:paramtypes'
// 'design:returntype'
export class MPI {

    public readonly name: string;
    public readonly version: string;
    public readonly proto: string;
    public readonly codec: string;
    public readonly flags: number;
    public readonly timeout: number;
    public readonly retries: number;
    public readonly node: string;
    public readonly inst: string;
    public readonly zone: string;
    public readonly cluster: string;
    public readonly cell: string;
    public readonly group: string;
    public readonly address: string;
    public kind: any = undefined;
    public retype: any = undefined;
    public attributes: Map<string, string> = new Map();

    constructor(name: string, version: string, proto: string, codec: string, flags: number, timeout: number,
                retries: number, node: string, inst: string, zone: string, cluster: string, cell: string, group: string,
                address: string, retype: any) {
        this.name = name;
        this.version = version;
        this.proto = proto;
        this.codec = codec;
        this.flags = flags;
        this.timeout = timeout;
        this.retries = retries;
        this.node = node;
        this.inst = inst;
        this.zone = zone;
        this.cluster = cluster;
        this.cell = cell;
        this.group = group;
        this.address = address;
        this.retype = retype;
    }

    decorate<T>(target: Object, propertyKey: string | symbol, descriptor?: TypedPropertyDescriptor<T>): TypedPropertyDescriptor<T> | void {
        if (typeof propertyKey === "string") {
            this.attributes.set('name', propertyKey);
        }
        if (!descriptor) {
            this.kind = Reflect.getMetadata('design:type', target, propertyKey);
            Ark.annotate(MPI, propertyKey, target, this);
            return;
        }
        this.kind = Reflect.getMetadata('design:returntype', target, propertyKey)
        Ark.annotate(MPI, descriptor.value, target.constructor, this);
        return descriptor;
    }

}

export function mpi(
    name = '',
    retype: any = Object,
    version = '1.0.0',
    proto = 'grpc',
    codec = 'json',
    flags = 0,
    timeout = 5000,
    retries = 3,
    node = '',
    inst = '',
    zone = '',
    cluster = '',
    cell = '',
    group = '',
    address = '',
): PropertyDecorator & MethodDecorator {
    const meta = new MPI(name, version, proto, codec, flags, timeout, retries, node, inst, zone, cluster, cell,
        group, address, retype);
    return meta.decorate.bind(meta);
}
