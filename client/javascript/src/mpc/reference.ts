/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Inspector, InvocationHandler, MPI, Return, ServiceLoader, Type} from "@/macro";
import {Invocation, Invoker, ServiceInvocation} from "@/mpc/invoker";
import {Context} from "@/prsim";
import {composite, CONSUMER} from "@/mpc/filter";
import {Reference} from "@/kinds";
import {Eden} from "@/mpc/eden";
import {URN} from "@/kinds/urn";
import Tool from "@/tool/tool";
import {context} from "@/mpc/context";
import {Consumer} from "@/mpc/consumer";
import {Mesh, MeshFlags} from "@/mpc/mesh";
import {Codec} from "@/codec";
import {Execution} from "@/mpc/execution";
import {errorf, Status} from "@/cause";


export class ReferenceInvokeHandler<T extends object> implements Invoker<any>, InvocationHandler {

    private readonly macro: any;
    private readonly invoker: Invoker<any>;

    constructor(macro: any) {
        this.macro = macro;
        this.invoker = composite(this, CONSUMER);
    }

    context(args: any[]): Context {
        if (!args || args.length < 1) {
            return context()
        }
        for (let arg of args) {
            if (arg instanceof Context) {
                return arg
            }
        }
        return context()
    }

    invoke(proxy: any, method: Inspector, args: any[]): any {
        return new Promise<any>(async (resolve, reject) => {
            try {
                const ctx = this.context(args);
                const execution = this.referExecution(ctx, method);
                const urn = this.rewriteURN(ctx, execution);
                ctx.rewriteURN(urn);
                ctx.setAttribute(Mesh.REMOTE, this.rewriteAddress(urn));

                const parameters = execution.inflect(ctx);
                parameters.setArguments(args);
                parameters.setAttachments(new Map<string, string>());

                const invocation = new ServiceInvocation(this, method, parameters, execution);

                ctx.setAttribute(Mesh.INVOCATION, invocation);

                resolve(await this.invoker.run(ctx, invocation));
            } catch (e) {
                reject(errorf(e))
            }
        });
    }

    async run(ctx: Context, invocation: Invocation): Promise<any> {
        const execution = this.referExecution(ctx, invocation.getInspector(ctx));
        const consumers = ServiceLoader.load(Consumer);
        const consumer = consumers.getDefault();
        const address = Tool.anyone(ctx.getAttribute(Mesh.REMOTE), Tool.MESH_ADDRESS.get().any());
        const name = execution.schema(ctx).codec || MeshFlags.JSON.name;
        const codecs = ServiceLoader.load(Codec);
        const codec = codecs.get(name);
        const buffer = codec.encode(invocation.getParameters(ctx));
        const pc = MeshFlags.ofName(consumers.defaultName()).code;
        const cc = MeshFlags.ofName(codecs.defaultName()).code;
        const urn = URN.from(ctx.getUrn()).resetFlag(pc, cc).toString();
        const future = await consumer.consume(ctx, address, urn, execution, buffer);
        return await this.deserialize(ctx, execution, codec, future)
    }

    deserialize(ctx: Context, execution: Execution<Reference>, codec: Codec, future: Uint8Array): Promise<any> {
        return new Promise<any>((resolve, reject) => {
            try {
                const retype: Type<Return> = execution.retype(ctx);
                const returns = codec.decode(future, retype);
                if (returns.getCause()) {
                    reject(errorf(returns.getCause()));
                    return;
                }
                if (Status.SUCCESS.getCode() != returns.getCode()) {
                    reject(returns.getMessage());
                    return;
                }
                resolve(returns.getContent());
            } catch (e) {
                reject(errorf(e));
            }
        });
    }

    referExecution(ctx: Context, inspector: Inspector): Execution<Reference> {
        const eden = ServiceLoader.load(Eden).getDefault();
        const execution = eden.refer(ctx, inspector.getAnnotation(MPI), inspector.getType(), inspector);
        if (execution) {
            return execution;
        }
        throw new Error(`Method ${inspector.getName()} cant be compatible`)
    }

    /**
     * Rewrite the urn by execution context.
     */
    private rewriteURN(ctx: Context, execution: Execution<Reference>): string {
        const nodeId = ctx.getPrincipals().peek()?.node_id;
        const instId = ctx.getPrincipals().peek()?.inst_id;
        const uname = ctx.getAttribute(Mesh.UNAME);
        const name = ctx.getAttribute(Mesh.REMOTE_NAME);
        if (!nodeId && !instId && !uname && !name) {
            return execution.schema(ctx).urn;
        }
        const urn = URN.from(execution.schema(ctx).urn);
        if (nodeId) {
            urn.nodeId = nodeId;
        }
        if (instId) {
            urn.nodeId = instId;
        }
        if ("" != uname) {
            urn.name = uname;
        }
        if ("" != name) {
            urn.name = urn.name.replace("${mesh.name}", name);
        }
        return urn.toString();
    }

    private rewriteAddress(uname: string): string {
        return Tool.MESH_ADDRESS.get().any();
    }

}