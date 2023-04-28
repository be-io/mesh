/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Context, Metadata} from "@/prsim";
import {Argument, Cause, Class, Inspector, Parameter, Return, ServiceLoader, Type, Types} from "@/macro";
import {Reference, Service} from "@/kinds";
import {URN} from "@/kinds/urn";
import {MeshFlags} from "@/mpc/mesh";
import {Codec} from "@/codec";

export interface Execution<T> {

    /**
     * Execution schema.
     * @param ctx
     */
    schema(ctx: Context): T;

    /**
     * Inspect execution.
     * @param ctx
     */
    inspect(ctx: Context): Inspector;

    /**
     * Execution input type.
     */
    intype(ctx: Context): Class<Parameter>;

    /**
     * Execution output return type.
     */
    retype(ctx: Context): Class<Return>;

    /**
     * Reflect input type.
     * @param ctx
     */
    inflect(ctx: Context): Parameter;

    /**
     * Reflect output return type.
     * @param ctx
     */
    reflect(ctx: Context): Return;
}

export class GenericExecution implements Execution<any> {

    private readonly inspector: Inspector;
    private readonly reference: Reference;
    private readonly service: Service;


    constructor(ctx: Context, urn: URN) {
        this.reference = new Reference();
        this.reference.urn = urn.toString();
        this.reference.namespace = "";
        this.reference.name = urn.name;
        this.reference.version = urn.flag.version;
        this.reference.proto = MeshFlags.ofProto(urn.flag.proto).name;
        this.reference.codec = MeshFlags.ofCodec(urn.flag.codec).name;
        this.reference.flags = 0;
        this.reference.timeout = parseInt(Metadata.MESH_TIMEOUT.get(ctx)) || 10000;
        this.reference.retries = 3;
        this.reference.node = urn.nodeId;
        this.reference.inst = "";
        this.reference.zone = urn.flag.zone;
        this.reference.cluster = urn.flag.cluster;
        this.reference.cell = urn.flag.cell;
        this.reference.group = urn.flag.group;
        this.reference.address = urn.flag.address;
        this.service = new Service();
        this.inspector = new GenericInspector(urn);
    }

    inflect(ctx: Context): Parameter {
        return new GenericParameters();
    }

    inspect(ctx: Context): Inspector {
        return this.inspector;
    }

    intype(ctx: Context): Class<Parameter> {
        return GenericParameters;
    }

    reflect(ctx: Context): Return {
        return new GenericReturns();
    }

    retype(ctx: Context): Class<Return> {
        return GenericReturns;
    }

    schema(ctx: Context): any {
        return this.reference;
    }

}

export class GenericInspector implements Inspector {

    private readonly urn: URN;

    constructor(urn: URN) {
        this.urn = urn;
    }

    getType(): Type<any> {
        return Object;
    }

    getAnnotation<T>(kind: Type<T>): T | undefined {
        return undefined;
    }

    getName(): string {
        return this.urn.name;
    }

    getParameters(): Argument[] {
        return [];
    }

    getSignature(): string {
        return this.urn.toString();
    }

    getReturnType(): Types {
        return new Types(Map)
    }

    invoke(target: any, args: any[]): any {
    }

}

export class GenericParameters extends Map<string, any> implements Parameter {

    map(): Map<string, any> {
        return this;
    }

    getArguments(): any[] {
        const args: any[] = [];
        super.forEach((v, k) => {
            if ("attachments" != k) {
                args.push(v);
            }
        })
        return args;
    }

    getAttachments(): Map<string, string> {
        const attachments = Object.getOwnPropertyDescriptor(this, "attachments")?.value || new Map();
        if (attachments instanceof Map) {
            return attachments;
        }
        const codec = ServiceLoader.load(Codec).get(Codec.JSON);
        const des = codec.decode(codec.encode(attachments), Map);
        this.setAttachments(des);
        return des;
    }

    setArguments(args: any[]): void {
        //
    }

    setAttachments(attachments: Map<string, string>): void {
        Object.assign(this, {attachments: Cause})
    }

}

export class GenericReturns extends Map<string, any> implements Return {

    getCause(): Cause | undefined {
        const cause = Object.getOwnPropertyDescriptor(this, "cause")?.value
        if (!cause || cause instanceof Cause) {
            return cause;
        }
        const codec = ServiceLoader.load(Codec).get(Codec.JSON);
        return codec.decode(codec.encode(cause), Cause);
    }

    getCode(): string {
        return Object.getOwnPropertyDescriptor(this, "code")?.value
    }

    getContent(): any {
        return Object.getOwnPropertyDescriptor(this, "content")?.value
    }

    getMessage(): string {
        return Object.getOwnPropertyDescriptor(this, "message")?.value || ''
    }

    setCause(cause: Cause): void {
        Object.assign(this, {cause: Cause})
    }

    setCode(code: string): void {
        Object.assign(this, {code: code})
    }

    setContent(content: any): void {
        Object.assign(this, {content: content})
    }

    setMessage(message: string): void {
        Object.assign(this, {message: message})
    }

}