/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Context, Network} from "@/prsim";
import {Environ, Institution, Page, Paging, Route, Versions} from "@/kinds";
import {ServiceProxy} from "@/mpc";

export class MeshNetwork extends Network {

    private readonly proxy: Network;

    constructor() {
        super();
        this.proxy = ServiceProxy.proxy(Network);
    }


    getEnviron(ctx: Context): Promise<Environ> {
        return this.proxy.getEnviron(ctx);
    }

    accessible(ctx: Context, route: Route): Promise<boolean> {
        return this.proxy.accessible(ctx, route);
    }

    refresh(ctx: Context, routes: Route[]): Promise<void> {
        return this.proxy.refresh(ctx, routes);
    }

    getRoute(ctx: Context, node_id: string): Promise<Route> {
        return this.proxy.getRoute(ctx, node_id);
    }

    getRoutes(ctx: Context): Promise<Route[]> {
        return this.proxy.getRoutes(ctx);
    }

    getDomains(ctx: Context): Promise<Route[]> {
        return this.proxy.getDomains(ctx);
    }

    putDomains(ctx: Context, domains: Route[]): Promise<void> {
        return this.proxy.putDomains(ctx, domains);
    }

    weave(ctx: Context, route: Route): Promise<void> {
        return this.proxy.weave(ctx, route);
    }

    ack(ctx: Context, route: Route): Promise<void> {
        return this.proxy.ack(ctx, route);
    }

    disable(ctx: Context, node_id: string): Promise<void> {
        return this.proxy.disable(ctx, node_id);
    }

    enable(ctx: Context, node_id: string): Promise<void> {
        return this.proxy.enable(ctx, node_id);
    }

    index(ctx: Context, index: Paging): Promise<Page<Route>> {
        return this.proxy.index(ctx, index);
    }

    version(ctx: Context, node_id: string): Promise<Versions> {
        return this.proxy.version(ctx, node_id);
    }

    instx(ctx: Context, index: Paging): Promise<Page<Institution>> {
        return this.proxy.instx(ctx, index);
    }

    instr(ctx: Context, institutions: Institution[]): Promise<void> {
        return this.proxy.instr(ctx, institutions);
    }

    ally(ctx: Context, node_ids: string[]): Promise<void> {
        return this.proxy.ally(ctx, node_ids);
    }

    disband(ctx: Context, node_ids: string[]): Promise<void> {
        return this.proxy.disband(ctx, node_ids);
    }
}