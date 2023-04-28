/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

// Invoke the next invoker.
import {Context} from "@/prsim";
import {Inspector, Parameter} from "@/macro";
import {Execution} from "@/mpc/execution";

export interface Invoker<T> {

    run(ctx: Context, invocation: Invocation): Promise<any>;

}

export interface Invocation {

    /**
     * Get the delegate target object.
     * @param ctx
     */
    getProxy(ctx: Context): Invoker<any>;

    /**
     * Get the invocation inspector.
     * @param ctx
     */
    getInspector(ctx: Context): Inspector;

    /**
     * Invoke parameters. include arguments and parameters.
     * @param ctx
     */
    getParameters(ctx: Context): Parameter;

    /**
     * Invoke parameters.
     * @param ctx
     */
    getArguments(ctx: Context): any[];

    /**
     * Get the attachments. The attributes will be serialized.
     * @param ctx
     */
    getAttachments(ctx: Context): Map<string, string>;

    /**
     * Get the invocation execution.
     * @param ctx
     */
    getExecution<V>(ctx: Context): Execution<V>;

    /**
     * Is the method return future.
     * @param ctx
     */
    isAsync(ctx: Context): boolean
}

export class ServiceInvocation<T> implements Invocation {

    private readonly proxy: Invoker<T>;
    private readonly inspector: Inspector;
    private readonly parameters: Parameter;
    private readonly execution: Execution<any>;

    constructor(proxy: Invoker<T>, inspector: Inspector, parameters: Parameter, execution: Execution<any>) {
        this.proxy = proxy;
        this.inspector = inspector;
        this.parameters = parameters;
        this.execution = execution;
    }

    getArguments(ctx: Context): any[] {
        return this.parameters.getArguments();
    }

    getAttachments(ctx: Context): Map<string, string> {
        return this.parameters.getAttachments();
    }

    getExecution<V>(ctx: Context): Execution<V> {
        // @ts-ignore
        return this.execution;
    }

    getInspector(ctx: Context): Inspector {
        return this.inspector;
    }

    getParameters(ctx: Context): Parameter {
        return this.parameters;
    }

    getProxy(ctx: Context): Invoker<any> {
        return this.proxy;
    }

    isAsync(ctx: Context): boolean {
        return false;
    }
}