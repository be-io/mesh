/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.Mesh;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.mpc.Types;
import io.be.mesh.prsim.Context;
import io.be.mesh.prsim.Network;
import io.be.mesh.tool.Once;
import io.be.mesh.tool.Tool;
import io.be.mesh.struct.*;

import java.lang.reflect.UndeclaredThrowableException;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshNetwork implements Network {

    private final Network environs = ServiceProxy.proxy(Network.class);
    private final Once<Environ> cache = new Once<>();
    private final Context.Key<Environ> environKey = new Context.Key<>("mesh.environ", Types.of(Environ.class));

    @Override
    public Environ getEnviron() {
        try {
            return Mesh.contextSafe(() -> {
                if (cache.isPresentWithoutGet()) {
                    return cache.get();
                }
                if (Tool.required(Mesh.context().getAttribute(environKey))) {
                    return Mesh.context().getAttribute(environKey);
                }
                Environ environ = new Environ();
                environ.setNodeId("LX0000000000000");
                environ.setInstId("JG0000000000000000");
                Mesh.context().setAttribute(environKey, environ);
                return cache.get(environs::getEnviron);
            });
        } catch (RuntimeException | Error e) {
            throw e;
        } catch (Throwable e) {
            throw new UndeclaredThrowableException(e);
        }
    }

    @Override
    public boolean accessible(Route route) {
        return environs.accessible(route);
    }

    @Override
    public void refresh(List<Route> routes) {
        environs.refresh(routes);
    }

    @Override
    public Route getRoute(String nodeId) {
        return environs.getRoute(nodeId);
    }

    @Override
    public List<Route> getRoutes() {
        return environs.getRoutes();
    }

    @Override
    public List<Route> getDomains() {
        return environs.getDomains();
    }

    @Override
    public void putDomains(List<Route> domains) {
        environs.putDomains(domains);
    }

    @Override
    public void weave(Route route) {
        environs.weave(route);
    }

    @Override
    public void ack(Route route) {
        environs.ack(route);
    }

    @Override
    public void disable(String nodeId) {
        environs.disable(nodeId);
    }

    @Override
    public void enable(String nodeId) {
        environs.enable(nodeId);
    }

    @Override
    public Page<List<Route>> index(Paging index) {
        return environs.index(index);
    }

    @Override
    public Versions version(String nodeId) {
        return environs.version(nodeId);
    }

    @Override
    public Page<List<Institution>> instx(Paging index) {
        return environs.instx(index);
    }

    @Override
    public void instr(List<Institution> institutions) {
        environs.instr(institutions);
    }

    @Override
    public void ally(List<String> nodeIds) {
        environs.ally(nodeIds);
    }

    @Override
    public void disband(List<String> nodeIds) {
        environs.disband(nodeIds);
    }

    @Override
    public boolean asserts(String feature, List<String> nodeIds) {
        return environs.asserts(feature, nodeIds);
    }

}
