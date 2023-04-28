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
import axios, {AxiosInstance, CreateAxiosDefaults} from "axios";
import {Codec} from "@/codec";
import Tool from "@/tool/tool";
import {errorf} from "@/cause";

@spi(Consumer.HTTP)
export class HTTPConsumer extends Consumer {

    private c: AxiosInstance | undefined;

    async getClient(): Promise<AxiosInstance> {
        if (this.c) {
            return this.c
        }
        const option: CreateAxiosDefaults = {
            baseURL: Tool.MESH_ADDRESS.get().any(),
            timeout: 10000,
            headers: {
                "Content-Type": "application/json;charset=UTF-8"
            },
        };
        if (typeof window === 'undefined') {
            const https = await import("https")
            option.httpsAgent = new https.Agent({
                rejectUnauthorized: false
            })
        }
        this.c = axios.create(option);
        this.c.interceptors.request.use(config => {
            //config.headers.token = 'token'
            return config
        });
        return this.c;
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

    async invoke<T>(param: any, type: T, urn: string): Promise<T> {
        const c = await this.getClient();
        const r = await c.post("/mesh/invoke", param, {
            headers: {
                "mesh-urn": urn,
            },
        })
        if (r.status !== 200) {
            throw errorf(r);
        }
        return r.data;
    }

}