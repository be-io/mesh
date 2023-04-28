/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {mpi} from "@/macro/mpi"
import {spi} from "@/macro/spi"
import {Environ} from "@/kinds/environ";
import {index} from "@/macro/idx";
import {Route} from "@/kinds/route";
import {Page, Paging} from "@/kinds/page";
import {Versions} from "@/kinds/versions";
import {Institution} from "@/kinds/institution";
import type {Context} from "@/prsim/context";
import {Status} from "@/cause";

/**
 * @author coyzeng@gmail.com
 */
@spi("mesh")
export abstract class Network {

    /**
     * Get the meth network environment fixed information.
     *
     * @return Fixed system information.
     */
    @mpi("mesh.net.environ", Environ)
    getEnviron(ctx: Context): Promise<Environ> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Check the mesh network is accessible.
     *
     * @return true is accessible.
     */
    @mpi("mesh.net.accessible", Boolean)
    accessible(ctx: Context, @index(0, 'route') route: Route): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Refresh the routes to mesh network.
     */
    @mpi("mesh.net.refresh")
    refresh(ctx: Context, @index(0, 'routes') routes: Route[]): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * getRoute the network edge route.
     *
     * @return edge routes
     */
    @mpi("mesh.net.edge", Route)
    getRoute(ctx: Context, @index(0, 'node_id') node_id: string): Promise<Route> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * getRoutes the network edge routes.
     *
     * @return edge routes
     */
    @mpi("mesh.net.edges", [Array, Route])
    getRoutes(ctx: Context): Promise<Route[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * GetNetDomain the network domains.
     *
     * @return net domains
     */
    @mpi("mesh.net.domains", [Array, Route])
    getDomains(ctx: Context): Promise<Route[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Put the network domains.
     */
    @mpi("mesh.net.resolve")
    putDomains(ctx: Context, @index(0, 'domains') domains: Route[]): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Weave the network.
     */
    @mpi("mesh.net.weave")
    weave(ctx: Context, @index(0, 'route') route: Route): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Acknowledge the network.
     */
    @mpi("mesh.net.ack")
    ack(ctx: Context, @index(0, 'route') route: Route): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Disable the network.
     */
    @mpi("mesh.net.disable")
    disable(ctx: Context, @index(0, 'node_id') node_id: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Enable the network.
     */
    @mpi("mesh.net.enable")
    enable(ctx: Context, @index(0, 'node_id') node_id: string): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Index the network edges.
     */
    @mpi("mesh.net.index", [Page, Route])
    index(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<Route>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Network environment version.
     *
     * @return version.
     */
    @mpi("mesh.net.version", Versions)
    version(ctx: Context, @index(0, 'node_id') node_id: string): Promise<Versions> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Network institutions.
     */
    @mpi("mesh.net.instx", [Page, Institution])
    instx(ctx: Context, @index(0, 'index') index: Paging): Promise<Page<Institution>> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Network institutions.
     */
    @mpi("mesh.net.instr")
    instr(ctx: Context, @index(0, 'institutions') institutions: Institution[]): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Network form alliance.
     */
    @mpi("mesh.net.ally")
    ally(ctx: Context, @index(0, 'node_ids') node_ids: string[]): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Network quit alliance.
     */
    @mpi("mesh.net.disband")
    disband(ctx: Context, @index(0, 'node_ids') node_ids: string[]): Promise<void> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Network feature assert.
     */
    @mpi("mesh.net.assert", Boolean)
    assert(ctx: Context, @index(0, "feature") feature: string, @index(1, 'node_ids') node_ids: string[]): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}
