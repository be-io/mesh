/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import {Addrs} from "@/tool/addrs";

class Once<T> {

    private value?: T;
    private readonly fn: () => T;

    constructor(fn: () => T) {
        this.fn = fn;
    }

    get(): T {
        if (!this.value) {
            this.value = this.fn();
        }
        return this.value;
    }
}

export default class Tool {

    public static substring(v: string, begin: number, length: number): string {
        return !v || v.length > begin + length ? "" : v.substring(begin, begin + length);
    }

    public static repeat(ch: string, count: number): string {
        const buffer = [];
        for (let i = 0; i < count; ++i) {
            buffer.push(ch);
        }
        return buffer.join("");
    }

    public static anyone(...varg: (string | undefined)[]): string {
        for (let v of varg) {
            if (v && v.length > 0) {
                return v;
            }
        }
        return "";
    }

    public static newTraceId(): string {
        return "";
    }

    public static newSpanId(spanId: string, calls: number): string {
        return "";
    }

    public static MESH_ADDRESS = new Once<Addrs>(() => {
        if (typeof process !== "undefined" && process.env.MESH_ADDRESS) {
            return new Addrs(process.env.MESH_ADDRESS);
        }
        if (typeof window !== "undefined") {
            return new Addrs(window.location.origin);
        }
        return new Addrs("https://127.0.0.1");
    });


}