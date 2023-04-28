/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import {expect, test} from 'vitest'
import {Context, Endpoint, KV, Network} from "@/prsim";
import {CONSUMER, context, Filter, Invocation, Invoker, Mesh, ServiceProxy} from "@/mpc";
import "@/http"
import "@/grpc"
import {Entity, Versions} from "@/kinds";
import {ServiceLoader, spi, Types} from "@/macro";
import {Codec} from "@/codec";

// Edit an assertion and save to see HMR in action

@spi("catch", CONSUMER, 1)
class CatchFilter extends Filter {

    invoke(ctx: Context, invoker: Invoker<any>, invocation: Invocation): Promise<any> {
        return invoker.run(ctx, invocation).catch(e => {
            console.log(e);
        });
    }

}

test('network.environ', async () => {
    const codec = ServiceLoader.load(Codec).getDefault();

    const version = new Versions();
    version.version = "1.5.0";
    version.infos = new Map();
    version.infos.set("1", "1");
    const vv = codec.decode(codec.encode(version), Versions);
    await expect(vv.infos.get("1")).eq("1");

    const endpoint = ServiceProxy.proxy(Endpoint);
    const ctx = context();
    ctx.setAttribute(Mesh.UNAME, "mesh.dot.exe");
    const output = await endpoint.fuzzy(ctx, codec.encode('mesh status'));
    await expect(output).instanceof(Uint8Array);

    const network = ServiceProxy.proxy(Network);
    const environ = await network.getEnviron(context());
    await expect(environ.node_id).toBe('LX0000010000110');

    const routes = await network.getRoutes(context());
    await expect(routes.length).gte(0);

    const kv = ServiceProxy.proxy(KV);
    await kv.put(context(), 'x', Entity.wrap(2));
    const entity = await kv.get(context(), 'x');
    const v = entity.readObject();
    await expect(v).eq(2);

    const dict: Map<string, Map<string, Map<string, string>>> = codec.decode(codec.encode({x: {x: {x: 'z'}}}), new Types([Map, String, [Map, String, [Map, String, String]]]));
    await expect(dict?.get("x")?.get("x")?.get("x")).eq('z');

    //const service = new Service();
});
//
// @spi("mesh")
// abstract class Interface {
//
//     @mpi("mesh.net.environ")
//     x(): Map<string, string> {
//         return new Map();
//     }
//
// }
//
//
// class Service {
//
//     @mpi("mesh")
//     public i!: Interface;
// }