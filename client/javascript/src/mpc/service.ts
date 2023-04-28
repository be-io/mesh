/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Invocation, Invoker} from "@/mpc/invoker";
import {Context} from "@/prsim";
import {composite, PROVIDER} from "@/mpc/filter";

class ServiceInvoker implements Invoker<any> {

    private readonly service: any

    constructor(service: any) {
        this.service = service;
    }

    run(ctx: Context, invocation: Invocation): Promise<any> {
        return new Promise<any>((resolve, reject) => {
            try {
                const args: any[] = [ctx, ...invocation.getArguments(ctx)];
                const ret = invocation.getInspector(ctx).invoke(this.service, args)
                if (ret instanceof Promise) {
                    ret.then(r => resolve(r)).catch(e => reject(e));
                    return;
                }
                resolve(ret);
            } catch (e) {
                reject(e)
            }
        })
    }

}

export class ServiceInvokeHandler implements Invoker<any> {

    private invoker: Invoker<any>;

    constructor(service: any) {
        this.invoker = composite(new ServiceInvoker(service), PROVIDER);
    }

    run(ctx: Context, invocation: Invocation): Promise<any> {
        return this.invoker.run(ctx, invocation)
    }

}