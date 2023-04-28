/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.*;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Network {

    /**
     * Get the meth network environment fixed information.
     *
     * @return Fixed system information.
     */
    @MPI("mesh.net.environ")
    Environ getEnviron();

    /**
     * Check the mesh network is accessible.
     *
     * @param route Network route.
     * @return true is accessible.
     */
    @MPI("mesh.net.accessible")
    boolean accessible(@Index(value = 0, name = "route") Route route);

    /**
     * Refresh the routes to mesh network.
     *
     * @param routes Network routes.
     */
    @MPI("mesh.net.refresh")
    void refresh(@Index(value = 0, name = "routes") List<Route> routes);

    /**
     * getRoute the network edge route.
     *
     * @return edge routes
     */
    @MPI("mesh.net.edge")
    Route getRoute(@Index(value = 0, name = "node_id") String nodeId);

    /**
     * getRoutes the network edge routes.
     *
     * @return edge routes
     */
    @MPI("mesh.net.edges")
    List<Route> getRoutes();

    /**
     * GetNetDomain the network domains.
     *
     * @return net domains
     */
    @MPI("mesh.net.domains")
    List<Route> getDomains();

    /**
     * Put the network domains.
     *
     * @param domains network domains for dns resovler
     */
    @MPI("mesh.net.resolve")
    void putDomains(List<Route> domains);

    /**
     * Weave the network.
     */
    @MPI("mesh.net.weave")
    void weave(@Index(value = 0, name = "route") Route route);

    /**
     * Acknowledge the network.
     */
    @MPI("mesh.net.ack")
    void ack(@Index(value = 0, name = "route") Route route);

    /**
     * Disable the network.
     */
    @MPI("mesh.net.disable")
    void disable(@Index(value = 0, name = "node_id") String nodeId);

    /**
     * Enable the network.
     */
    @MPI("mesh.net.enable")
    void enable(@Index(value = 0, name = "node_id") String nodeId);

    /**
     * Index the network edges.
     */
    @MPI("mesh.net.index")
    Page<List<Route>> index(Paging index);

    /**
     * Network environment version.
     *
     * @return version.
     */
    @MPI("mesh.net.version")
    Versions version(@Index(value = 0, name = "node_id") String nodeId);

    /**
     * Network institutions.
     */
    @MPI("mesh.net.instx")
    Page<List<Institution>> instx(Paging index);

    /**
     * Network institutions.
     */
    @MPI("mesh.net.instr")
    void instr(List<Institution> institutions);

    /**
     * Network form alliance.
     */
    @MPI("mesh.net.ally")
    void ally(@Index(value = 0, name = "node_ids") List<String> nodeIds);

    /**
     * Network quit alliance.
     */
    @MPI("mesh.net.disband")
    void disband(@Index(value = 0, name = "node_ids") List<String> nodeIds);

    /**
     * Network feature assert.
     */
    @MPI("mesh.net.assert")
    boolean asserts(@Index(0) String feature, @Index(name = "node_ids", value = 1) List<String> nodeIds);

}
