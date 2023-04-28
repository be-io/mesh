/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.cause.TimeoutException;
import com.be.mesh.client.prsim.Network;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.tool.Mode;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;

import java.time.Duration;
import java.util.concurrent.locks.LockSupport;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI(value = "robust", pattern = Filter.CONSUMER, priority = Integer.MAX_VALUE)
public class RobustFilter implements Filter {

    private final Network network = ServiceLoader.load(Network.class).getDefault();

    @Override
    public Object invoke(Invoker<?> invoker, Invocation invocation) throws Throwable {
        Execution<Reference> execution = invocation.getExecution();
        int retries = Math.min(execution.schema().getRetries(), 3);
        for (int i = 0; i < retries - 1; i++) {
            try {
                Object ret = invoker.invoke(invocation);
                Tool.MESH_ADDRESS.get().available(Mesh.context().getAttribute(Mesh.REMOTE), true);
                return ret;
            } catch (MeshException e) {
                if (e.unavailable() && isUnavailableTrust(invocation.getUrn())) {
                    Tool.MESH_ADDRESS.get().available(Mesh.context().getAttribute(Mesh.REMOTE), false);
                    Mesh.context().setAttribute(Mesh.REMOTE, Tool.MESH_ADDRESS.get().any());
                }
                if (isHealthCheckURN(invocation.getUrn()) || cantRetry(e)) {
                    throw e;
                }
            } catch (Exception e) {
                if (isHealthCheckURN(invocation.getUrn()) || cantRetry(e)) {
                    throw e;
                }
            }
        }
        return invoker.invoke(invocation);
    }

    private boolean isUnavailableTrust(URN urn) {
        if (isHealthCheckURN(urn)) {
            return true;
        }
        if (!Mode.FAILFAST.match(Tool.MESH_MODE.get())) {
            return false;
        }
        return Tool.isInMyNet(network.getEnviron().getNodeId(), urn.getNodeId());
    }

    private boolean isHealthCheckURN(URN urn) {
        return Tool.contains(urn.getName(), "mesh");
    }

    private boolean cantRetry(Exception e) throws Exception {
        if (!shouldRetry(e)) {
            return true;
        }
        log.warn("Retry with {}:{}", e.getClass().getName(), e.getMessage());
        if (shouldDelay(e)) {
            LockSupport.parkNanos(Duration.ofMillis(30).toNanos());
        }
        return false;
    }

    public boolean shouldDelay(Throwable e) {
        boolean delay = instanceOf(e, java.net.SocketException.class);
        delay = delay || instanceOf(e, java.net.SocketTimeoutException.class);
        delay = delay || instanceOf(e, TimeoutException.class);
        return delay;
    }

    public boolean shouldRetry(Throwable e) {
        boolean can = instanceOf(e, java.net.SocketTimeoutException.class);
        can = can || instanceOf(e, java.net.SocketException.class);
        can = can || instanceOf(e, java.net.ConnectException.class);
        can = can || instanceOf(e, TimeoutException.class);
        can = can || instanceOf(e, java.sql.SQLException.class);
        return can;
    }

    private boolean instanceOf(Throwable e, Class<? extends Throwable> ex) {
        if (ex.isAssignableFrom(e.getClass())) {
            return true;
        }
        return null != e.getCause() && ex.isAssignableFrom(e.getCause().getClass());
    }
}