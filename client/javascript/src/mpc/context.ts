/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Context, Deque, Key} from "@/prsim";
import {Location, Principal} from '@/kinds'
import Tool from "@/tool/tool";

export function context(): Context {
    return MeshContext.create();
}

export class MeshContext extends Context {

    private traceId: string = "";
    private spanId: string = "";
    private timestamp: number = 0;
    private runMode: number = 0;
    private urn: string = "";
    private consumer: Location = new Location();
    private calls: number = 0;
    private attachments: Map<string, string> = new Map<string, string>();
    private attributes: Map<string, any> = new Map<string, any>();
    private principals: Deque<Principal> = new Deque<Principal>();

    public static create(): Context {
        const context = new MeshContext();
        context.traceId = Tool.newTraceId();
        context.spanId = Tool.newSpanId("", 0);
        context.timestamp = new Date().getTime();
        context.runMode = 0;
        context.urn = "";
        context.consumer = new Location();
        return context;
    }

    getAttachments(): Map<string, string> {
        return this.attachments;
    }

    getAttribute<T>(key: Key<T>): any {
        return this.attributes.get(key.name)
    }

    getAttributes(): Map<string, any> {
        return this.attributes;
    }

    getConsumer(ctc: Context): Location {
        return this.consumer;
    }

    getPrincipals(): Deque<Principal> {
        return this.principals;
    }

    getProvider(ctx: Context): Location {
        return new Location();
    }

    getRunMode(): number {
        return this.runMode;
    }

    getSpanId(): string {
        return this.spanId;
    }

    getTimestamp(): number {
        return this.timestamp;
    }

    getTraceId(): string {
        return this.traceId;
    }

    getUrn(): string {
        return this.urn;
    }

    setAttribute<T>(key: Key<T>, value: any): void {
        this.attributes.set(key.name, value);
    }

    resume(): Context {
        this.calls++;
        const context = new MeshContext();
        context.traceId = this.getTraceId();
        context.spanId = Tool.newSpanId(this.getSpanId(), this.calls);
        context.timestamp = this.getTimestamp();
        context.runMode = this.getRunMode();
        context.urn = this.getUrn();
        context.consumer = this.getConsumer(this);
        if (this.getAttachments()) {
            this.getAttachments().forEach((k, v) => context.attachments.set(k, v));
        }
        if (this.getAttributes()) {
            this.getAttributes().forEach((k, v) => context.attributes.set(k, v));
        }
        if (this.getPrincipals()) {
            this.getPrincipals().forEach((v) => context.principals.push(v));
        }
        return context;
    }

    rewriteContext(ctx: Context): void {
        if (ctx.getTraceId()) {
            this.traceId = ctx.getTraceId();
        }
        if (ctx.getSpanId()) {
            this.spanId = ctx.getSpanId();
        }
        if (ctx.getTimestamp() && ctx.getTimestamp() > 0) {
            this.timestamp = ctx.getTimestamp();
        }
        if (ctx.getRunMode() && ctx.getRunMode() > 0) {
            this.runMode = ctx.getRunMode();
        }
        if (ctx.getUrn()) {
            this.urn = ctx.getUrn();
        }
        if (ctx.getConsumer(this)) {
            this.consumer = ctx.getConsumer(this);
        }
        if (this.getAttachments()) {
            this.getAttachments().forEach((k, v) => this.attachments.set(k, v));
        }
        if (this.getAttributes()) {
            this.getAttributes().forEach((k, v) => this.attributes.set(k, v));
        }
        if (this.getPrincipals()) {
            this.getPrincipals().forEach((v) => this.principals.push(v));
        }
    }

    rewriteURN(urn: string): void {
        this.urn = urn;
    }

}