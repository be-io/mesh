/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import {Ark, Type} from "@/macro/ark";
import {spi} from "@/macro/spi";
import {MPI} from "@/macro/mpi";
import {dict, Dict} from "@/macro/dsa";
import {Idx} from "@/macro/idx";


const STRIP_COMMENTS = /((\/\/.*$)|(\/\*[\s\S]*?\*\/))/mg;
const ARGUMENT_NAMES = /([^\s,]+)/g;


export interface Inspector {
    /**
     * Get the inspector signature
     */
    getSignature(): string;

    /**
     * Get inspect parent type.
     */
    getType<T>(): Type<T>;

    /**
     * Get the name of inspector.
     */
    getName(): string;

    /**
     * Get the annotations of inspector.
     */
    getAnnotation<T>(kind: Type<T>): T | undefined;

    /**
     * Get the return type of inspector.
     */
    getParameters(): Argument[];

    /**
     * Get the origin method.
     */
    getReturnType(): Types;

    /**
     * Invoke the inspector object.
     */
    invoke(target: any, args: any[]): any;
}

export interface Accessor {

    /**
     * Get the executable accessors.
     */
    getMethods(): Inspector[];

}

@spi("mesh")
export abstract class MethodProxy {

    /**
     * Proxy interface.
     */
    abstract proxy<T>(kind: Type<T>): T;
}

/**
 * InvocationHandler is the interface implemented by the invocation handler of a proxy instance.
 * Each proxy instance has an associated invocation handler. When a method is invoked on a proxy instance,
 * the method invocation is encoded and dispatched to the invoke method of its invocation handler.
 */
export interface InvocationHandler {

    invoke(proxy: any, method: Inspector, args: any[]): any;
}

export interface Parameter {
    /**
     * Convert parameters to map.
     */
    map(): Map<string, any>;

    /**
     * Generic arguments array.
     */
    getArguments(): any[];

    /**
     * Generic arguments array.
     */
    setArguments(args: any[]): void;

    /**
     * Get the generic attachments. The attributes will be serialized. The attachments are mutable.
     */
    getAttachments(): Map<string, string>;

    /**
     * Attachment arguments.
     */
    setAttachments(attachments: Map<string, string>): void;
}

export interface Return {
    /**
     * Return code.
     */
    getCode(): string;

    /**
     * Return code.
     */
    setCode(code: string): void;

    /**
     * Return message.
     */
    getMessage(): string;

    /**
     * Return message.
     */
    setMessage(message: string): void;

    /**
     * Return cause.
     */
    getCause(): Cause | undefined;

    /**
     * Return cause.
     */
    setCause(cause: Cause): void;

    /**
     * Return content.
     */
    getContent(): any;

    /**
     * Return content.
     */
    setContent(content: any): void;
}

export class Argument {
    public index: number;
    public name: string;

    constructor(index: number, name: string) {
        this.index = index;
        this.name = name;
    }
}

export class Types implements Function {

    public readonly raw: Type<any>;
    public readonly parameters: Types[];
    prototype: any;
    length: number;
    arguments: any;
    caller: Function = () => {
    };
    name: string;

    constructor(types: any) {
        [this.raw, this.parameters] = this.parse(types);
        this.prototype = this.raw.prototype;
        this.length = this.raw.length;
        // this.arguments = this.raw.arguments;
        // this.caller = this.raw.caller;
        this.name = this.raw.name;
    }

    parse(types: any): any[] {
        if (!types) {
            return [Object, []];
        }
        if (!Array.isArray(types)) {
            return [types, []];
        }
        if (types.length < 2) {
            return [types[0], []];
        }
        return [types[0], types.slice(1).map(v => Array.isArray(v) ? new Types(v) : v)];
    }

    apply = (thisArg: any, argArray?: any) => {
        return this.raw.apply(thisArg, ...argArray);
    }

    call = (thisArg: any, ...argArray: any[]) => {
        return this.raw.call(thisArg, ...argArray);
    }

    bind = (thisArg: any, ...argArray: any[]) => {
        return this.raw.bind(thisArg, ...argArray);
    }

    toString(): string {
        return this.raw.toString();
    }

    [Symbol.hasInstance](value: any): boolean {
        return this.raw[Symbol.hasInstance](value);
    }

}

export class Cause {
    public name: string = "";
    public pos: string = "";
    public text: string = "";
    public buff: Uint8Array = new Uint8Array();

    public static of(err: Error): Cause {
        return new Cause();
    }

    public static ofCause(code: string, message: string, cause: Cause): Error {
        throw new Error()
    }
}

export class MethodInspector implements Inspector {

    private readonly kind: any;
    private readonly name: string;
    private readonly method: any;

    constructor(kind: any, name: string, method: any) {
        this.kind = kind;
        this.name = name;
        this.method = method;
    }

    getAnnotation<T>(macro: Type<T>): T | undefined {
        return Ark.metadata(macro, this.kind, this.method)
    }

    getName(): string {
        const metadata = Ark.metadata(MPI, this.kind, this.method);
        return metadata?.attributes?.get('name') || this.name;
    }

    getParameters(): Argument[] {
        const fns = this.method.toString().replace(STRIP_COMMENTS, '');
        const names = fns.slice(fns.indexOf('(') + 1, fns.indexOf(')')).match(ARGUMENT_NAMES);
        const args: Argument[] = [];
        for (let idx = 0; idx < names.length; idx++) {
            const metadata = Ark.metadata(Idx, this.kind, `${this.getName()}:${idx}`);
            if (metadata) {
                args.push(new Argument(metadata.index, metadata.name));
            }
        }
        return args;
    }

    getSignature(): string {
        return `${this.kind.constructor.name}{${this.getName()}(${this.getParameters().map(x => x.name).join(",")})}`;
    }

    getType<T>(): Type<T> {
        return this.kind;
    }

    getReturnType(): Types {
        return new Types(Ark.metadata(MPI, this.kind, this.method)?.retype)
    }

    invoke(target: any, args: any[]): any {
        return this.method(...args);
    }


}

class executable {

    private readonly inspector: Inspector;
    private readonly h: InvocationHandler;
    private readonly proxy: any;

    constructor(inspector: Inspector, h: InvocationHandler, proxy: any) {
        this.inspector = inspector;
        this.h = h;
        this.proxy = proxy;
    }

    apply<T>(method: T, target: any, args: any[]): any {
        return this.h.invoke(this.proxy, this.inspector, args);
    }
}

class proxy<T> implements Accessor {

    private readonly inspectors: Dict<string, Inspector> = new dict();
    private readonly executables: Dict<string, any> = new dict();

    constructor(kind: Type<T>, h: InvocationHandler) {
        Ark.inspect(kind).forEach((v, k) => {
            const inspector = new MethodInspector(kind, k, v);
            this.executables.set(k, new Proxy(v, new executable(inspector, h, this)));
            this.inspectors.set(k, inspector);
        })
    }

    getMethods(): Inspector[] {
        return this.inspectors.map((k, v) => v);
    }

    get(target: Function & { prototype: T }, p: string, receiver: any): any {
        return this.executables.get(p);
    }

}

export class Dynamic {

    private static proxies: Map<any, Map<any, any>> = new Map();

    public static newProxyInstance<T>(kind: Type<T>, h: InvocationHandler): T {
        if (!this.proxies.get(kind)) {
            this.proxies.set(kind, new Map());
        }
        if (!this.proxies.get(kind)?.get(h)) {
            this.proxies.get(kind)?.set(h, new Proxy(kind, new proxy(kind, h)));
        }
        return this.proxies.get(kind)?.get(h);
    }

}