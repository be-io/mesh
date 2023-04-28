/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Dynamic, MethodProxy, mpi, spi, Type} from "@/macro";
import {ReferenceInvokeHandler} from "@/mpc/reference";

@spi("mesh")
export class ServiceProxy extends MethodProxy {

    proxy<T>(kind: Type<T>): T {
        return ServiceProxy.proxy(kind);
    }

    public static proxy<T>(kind: Type<T>): T {
        return this.doProxy(kind, mpi(''));
    }

    public static doProxy<T>(kind: Type<T>, metadata: any): any {
        return Dynamic.newProxyInstance(kind, new ReferenceInvokeHandler(metadata)) as T;
    }


}
