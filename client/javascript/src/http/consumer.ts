/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Consumer, Execution} from "@/mpc";
import {ServiceLoader, spi} from "@/macro";
import {Context} from "@/prsim";
import {Reference} from "@/kinds";
import axios, {AxiosInstance} from "axios";
import {Codec} from "@/codec";
import Tool from "@/tool/tool";
import {errorf} from "@/cause";

@spi(Consumer.HTTP)
export class HTTPConsumer extends Consumer {

    private client: AxiosInstance;

    constructor() {
        super();
        this.client = axios.create({
            baseURL: Tool.MESH_ADDRESS.get().any(),
            timeout: 10000,
            headers: {
                "Content-Type": "application/json;charset=UTF-8"
            },
            // httpsAgent: new https.Agent({
            //     rejectUnauthorized: false
            // })
        });
        this.client.interceptors.request.use(config => {
            //config.headers.token = 'token'
            return config
        });
    }

    start(): void {
        //
    }

    close(): void {
        //
    }

    async consume(ctx: Context, address: string, urn: string, execution: Execution<Reference>, inbound: Uint8Array): Promise<Uint8Array> {
        const codec = ServiceLoader.load(Codec).getDefault();
        const json = codec.decode(inbound, Object);
        const outbound = await this.invoke(json, Object, urn);
        return codec.encode(outbound);
    }

    invoke<T>(param: any, type: T, urn: string): Promise<T> {
        return new Promise<T>((resolve, reject) => {
            this.client.post("/mesh/invoke", param, {
                headers: {
                    "mesh-urn": urn,
                },
            }).then(r => {
                if (r.status !== 200) {
                    reject(errorf(r))
                    return
                }
                resolve(r.data);
            }).catch(e => {
                reject(errorf(e));
            });
        });
    }

}