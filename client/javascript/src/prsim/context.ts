/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import {Location, Principal} from "@/kinds/principal";

export class Key<T> {

    public name: string;

    constructor(name: string) {
        this.name = name;
    }
}

export class Deque<T> extends Array {

    peek(): T | undefined {
        if (this.length < 1) {
            return undefined;
        }
        return this.at(this.length - 1);
    }
}

export abstract class Context {

    // GetTraceId the request trace id.
    abstract getTraceId(): string;

    // GetSpanId the request span id.
    abstract getSpanId(): string;

    // getTimestamp the request create time.
    abstract getTimestamp(): number

    // getRunMode the request run mode. {@link RunMode}
    abstract getRunMode(): number

    // getUrn mesh resource uniform name. Like: create.tenant.omega.json.http2.lx000001.mpi.trustbe.net
    abstract getUrn(): string

    // getConsumer the consumer network principal.
    abstract getConsumer(ctc: Context): Location;

    // getProvider the provider network principal.
    abstract getProvider(ctx: Context): Location;

    // getAttachments Dispatch attachments.
    abstract getAttachments(): Map<string, string>;

    // getPrincipals Get the mpc broadcast network principals.
    abstract getPrincipals(): Deque<Principal>;

    // getAttributes like getAttachments, but attribute wont be transfer in invoke chain.
    abstract getAttributes(): Map<string, any>;

    // getAttribute like getAttachments, but attribute wont be transfer in invoke chain.
    abstract getAttribute<T>(key: Key<T>): any;

    // setAttribute Like putAttachments, but attribute won't be transfer in invoke chain.
    abstract setAttribute<T>(key: Key<T>, value: any): void;

    // rewriteURN rewrite the urn.
    abstract rewriteURN(urn: string): void;

    // rewriteContext rewrite the context by another context.
    abstract rewriteContext(ctx: Context): void;

    // resume will open a new context.
    abstract resume(): Context;
}

/**
 * https://www.rfc-editor.org/rfc/rfc7540#section-8.1.2
 */
export class Metadata {
    public static MESH_TRACE_ID = new Metadata("mesh-trace-id");
    public static MESH_SPAN_ID = new Metadata("mesh-span-id");
    public static MESH_TIMESTAMP = new Metadata("mesh-timestamp");
    public static MESH_RUN_MODE = new Metadata("mesh-run-mode");
    public static MESH_CONSUMER = new Metadata("mesh-consumer");
    public static MESH_PROVIDER = new Metadata("mesh-provider");
    public static MESH_URN = new Metadata("mesh-urn");
    public static MESH_FROM_INST_ID = new Metadata("mesh-source-inst-id");
    public static MESH_FROM_NODE_ID = new Metadata("mesh-source-node-id");
    public static MESH_INCOMING_HOST = new Metadata("mesh-incoming-host");
    public static MESH_OUTGOING_HOST = new Metadata("mesh-outgoing-host");
    public static MESH_INCOMING_PROXY = new Metadata("mesh-incoming-proxy");
    public static MESH_OUTGOING_PROXY = new Metadata("mesh-outgoing-proxy");
    public static MESH_SUBSET = new Metadata("mesh-subset");
    public static MESH_SESSION_ID = new Metadata("mesh-session-id");
    public static MESH_VERSION = new Metadata("mesh-version");
    public static MESH_TIMEOUT = new Metadata("mesh-timeout");
    // PTP
    public static MESH_PTP_VERSION = new Metadata("x-ptp-version");
    public static MESH_PTP_TPC = new Metadata("x-ptp-tech-provider-code");
    public static MESH_PTP_TRACE_ID = new Metadata("x-ptp-trace-id");
    public static MESH_PTP_TOKEN = new Metadata("x-ptp-token");
    public static MESH_PTP_URI = new Metadata("x-ptp-uri");
    public static MESH_PTP_SOURCE_NODE_ID = new Metadata("x-ptp-source-node-id");
    public static MESH_PTP_TARGET_NODE_ID = new Metadata("x-ptp-target-node-id");
    public static MESH_PTP_SOURCE_INST_ID = new Metadata("x-ptp-source-inst-id");
    public static MESH_PTP_TARGET_INST_ID = new Metadata("x-ptp-target-inst-id");
    public static MESH_PTP_SESSION_ID = new Metadata("x-ptp-session-id");
    public static MESH_PTP_TOPIC = new Metadata("x-ptp-topic");
    public static MESH_PTP_TIMEOUT = new Metadata("x-ptp-timeout");

    private readonly key: string;

    constructor(key: string) {
        this.key = key;
    }

    public get(ctx: Context): string {
        const v = ctx.getAttachments()?.get(this.key);
        if (v) {
            return v;
        }
        const lv = ctx.getAttachments()?.get(this.key.replace("-", "_"));
        if (lv) {
            return lv;
        }
        for (let kv of ctx.getAttachments()?.entries()) {
            if (kv[1] && kv[0] && (kv[0].toLowerCase() == this.key || kv[0].replace("_", "-").toLowerCase() == this.key)) {
                return kv[1];
            }
        }
        return "";
    }

    public set(ctx: Context, value: string): void {
        ctx.getAttachments()?.set(this.key, value);
    }
}
