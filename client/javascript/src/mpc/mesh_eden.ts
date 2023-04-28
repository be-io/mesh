/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Eden} from "@/mpc/eden";
import {Context} from "@/prsim";
import {
    Ark,
    Class,
    Dict,
    dict,
    Inspector,
    MethodInspector,
    MPI,
    Parameter,
    Return,
    ServiceLoader,
    spi,
    Type
} from "@/macro";
import {Invoker} from "@/mpc/invoker";
import {Environ, Reference, Service} from "@/kinds";
import {ServiceProxy} from "@/mpc/proxy";
import {ServiceInvokeHandler} from "@/mpc/service";
import {Compiler} from "@/mpc/compiler";
import {URN, URNFlag} from "@/kinds/urn";
import {Execution, GenericExecution} from "@/mpc/execution";
import Tool from "@/tool/tool";
import {MeshFlags} from "@/mpc/mesh";

class Provider {

    constructor(making: boolean, kind: Type<any>, service: any) {
        this.making = making;
        this.kind = kind;
        this.service = service;
    }

    public making: boolean;
    public kind: Type<any>;
    public service: any;
}

class Consumer {

    public metadata: any;
    public reference: Type<any>;
    public proxy: any;

    constructor(metadata: any, reference: Type<any>) {
        this.metadata = metadata;
        this.reference = reference;
        this.proxy = ServiceProxy.doProxy(reference, metadata);
    }
}

class Instance<T> implements Execution<any> {

    private readonly urn: string;
    private readonly kind: Type<T>;
    private readonly method: any;
    private readonly target: any;
    private readonly resource: T;
    private readonly inspector: Inspector;
    private readonly invoker: Invoker<any>;
    private readonly __intype: Class<any>;
    private readonly __retype: Class<any>;

    constructor(urn: string, kind: Type<T>, method: MethodInspector, target: any, resource: T) {
        this.urn = urn;
        this.kind = kind;
        this.method = method;
        this.target = target;
        this.resource = resource;
        this.inspector = method;
        this.invoker = new ServiceInvokeHandler(target);
        this.__intype = ServiceLoader.load(Compiler).getDefault().intype(method);
        this.__retype = ServiceLoader.load(Compiler).getDefault().retype(method);
    }

    inflect(ctx: Context): Parameter {
        return new this.__intype();
    }

    inspect(ctx: Context): Inspector {
        return this.inspector;
    }

    intype(ctx: Context): Class<Parameter> {
        return this.__intype;
    }

    reflect(ctx: Context): Return {
        return new this.__retype();
    }

    retype(ctx: Context): Class<Return> {
        return this.__retype;
    }

    schema(ctx: Context): any {
        return this.resource;
    }
}

@spi("mesh")
export class MeshEden extends Eden {

    private providers: Dict<any, Provider> = new dict();
    private consumers: Dict<any, Dict<any, Consumer>> = new dict();
    private indies: Dict<any, Dict<any, Dict<any, Execution<Reference>>>> = new dict();
    private services: Dict<string, Execution<Service>> = new dict();
    private references: Dict<string, Execution<Reference>> = new dict();


    define(ctx: Context, metadata: any, reference: Type<any>): any {
        return this.makeConsumer(ctx, metadata, reference).proxy;
    }

    store(ctx: Context, kind: Type<any>, service: any): void {
        this.providers.set(kind, new Provider(false, kind, service));
    }

    infer(ctx: Context, urn: string): Execution<Service> {
        const domain = URN.from(urn);
        const executions = this.makeServiceExecution(this.getEnviron(ctx));
        const execution = executions?.get(domain.name);
        if (execution) {
            return execution;
        }
        for (let kv of executions.entries()) {
            if (domain.matchName(kv[0])) {
                return kv[1];
            }
        }
        return new GenericExecution(ctx, domain);
    }

    refer(ctx: Context, metadata: any, reference: Type<any>, method: Inspector): Execution<Reference> {
        const consumer = this.makeConsumer(ctx, metadata, reference);
        const env = this.getEnviron(ctx);
        const execution = this.indies.computeIfy(reference, k => {
            return new dict<any, Dict<any, Execution<Reference>>>();
        }).computeIfy(metadata, k => {
            const executions = new dict<any, Execution<Reference>>();
            this.getMethods(ctx, reference).forEach((methods, kind) => methods.forEach(m => {
                const refer = this.makeMethodAsReference(env, metadata, kind, m);
                const instance = new Instance(refer.urn, reference, m, consumer.proxy, refer);
                this.references.set(refer.urn, instance);
                executions.set(m.getSignature(), instance);
            }))
            return executions;
        }).get(method.getSignature());
        if (execution) {
            return execution;
        }
        throw new Error(`Method ${method.getName()} cant be compatible`);
    }

    inferTypes(ctx: Context): Type<any>[] {
        const types: Type<any>[] = [];
        for (let kv of this.providers) {
            types.push(kv[1].kind);
        }
        return types;
    }

    referTypes(ctx: Context): Type<any>[] {
        const types: Type<any>[] = [];
        for (let key of this.consumers.keys()) {
            types.push(key);
        }
        return types;
    }


    private makeConsumer(ctx: Context, metadata: any, reference: Type<any>): Consumer {
        return this.consumers.computeIfy(reference, k => {
            return new dict<any, Consumer>()
        }).computeIfy(metadata, k => {
            return new Consumer(metadata, reference);
        });
    }


    private makeMethodAsReference(env: Environ, metadata: any, kind: Type<any>, method: MethodInspector): Reference {
        const reference = new Reference();
        reference.namespace = kind.name;
        reference.name = method.getName();
        reference.version = method.getAnnotation(MPI)?.version || '';
        reference.proto = method.getAnnotation(MPI)?.proto || '';
        reference.codec = method.getAnnotation(MPI)?.codec || '';
        reference.flags = method.getAnnotation(MPI)?.flags || 0;
        reference.timeout = method.getAnnotation(MPI)?.timeout || 12000;
        reference.retries = method.getAnnotation(MPI)?.retries || 0;
        reference.node = method.getAnnotation(MPI)?.node || '';
        reference.inst = method.getAnnotation(MPI)?.inst || '';
        reference.zone = method.getAnnotation(MPI)?.zone || '';
        reference.cluster = method.getAnnotation(MPI)?.cluster || '';
        reference.cell = method.getAnnotation(MPI)?.cell || '';
        reference.group = method.getAnnotation(MPI)?.group || '';
        reference.address = method.getAnnotation(MPI)?.address || '';
        reference.urn = this.getURN(method.getAnnotation(MPI)?.name || '', Service, env);
        return reference;
    }

    private getURN(alias: string, definition: any, env: Environ): string {
        const urn = new URN();
        urn.domain = URN.MESH_DOMAIN;
        urn.nodeId = Tool.anyone(definition.node, definition.inst, env.node_id);
        urn.name = alias;
        urn.flag = this.getURNFlag(definition);
        return urn.toString();
    }

    private getURNFlag(definition: any): URNFlag {
        const authority = (definition.address || "").split(":");
        const flag = new URNFlag();
        flag.v = "00";
        flag.proto = MeshFlags.ofName(definition.proto).code;
        flag.codec = MeshFlags.ofName(definition.codec).code;
        flag.version = definition.version;
        flag.zone = definition.zone;
        flag.cluster = definition.cluster;
        flag.cell = definition.cell;
        flag.group = definition.group;
        flag.address = authority.length > 0 ? authority[0] : "";
        flag.port = authority.length > 1 ? authority[1] : "";
        return flag;
    }

    private getMethods(ctx: Context, reference: Type<any>): Map<any, MethodInspector[]> {
        const methods = Ark.inspect(reference).map((k, v) => new MethodInspector(reference, k, v));
        const inspectors: Map<any, MethodInspector[]> = new dict();
        inspectors.set(reference, methods);
        return inspectors;
    }

    private getEnviron(ctx: Context): Environ {
        const environ = new Environ();
        environ.node_id = URN.LocalNodeId;
        environ.inst_id = URN.LocalInstId;
        return environ;
    }

    private makeServiceExecution(env: Environ): Map<string, Execution<Service>> {
        for (let kv of this.providers) {
            if (kv[1].making) {
                continue;
            }
        }
        return this.services;
    }

}