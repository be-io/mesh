/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import {
    CONSUMER,
    Context,
    DataHouse,
    Endpoint,
    Filter,
    Invocation,
    Invoker,
    KV,
    Network,
    Registry,
    ServiceProxy,
    spi
} from "@mesh/mesh";

@spi("catch", CONSUMER)
class CatchFilter implements Filter {

    invoke(ctx: Context, invoker: Invoker<any>, invocation: Invocation): Promise<any> {
        return invoker.run(ctx, invocation).catch(e => {
            console.log(e);
        });
    }

}

class Services {

    public network: Network = ServiceProxy.proxy(Network);
    public datahouse: DataHouse = ServiceProxy.proxy(DataHouse);
    public registry: Registry = ServiceProxy.proxy(Registry);
    public endpoint: Endpoint = ServiceProxy.proxy(Endpoint);
    public kv: KV = ServiceProxy.proxy(KV);

}

const services = new Services()

export default services;