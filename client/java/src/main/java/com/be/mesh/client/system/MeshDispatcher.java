/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.prsim.Dispatcher;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshDispatcher implements Dispatcher {

    private final Map<String, Object> refers = new ConcurrentHashMap<>();
    private final GenericInvokeHandler invoker = new GenericInvokeHandler();

    @SuppressWarnings("unchecked")
    @Override
    public <T> T reference(Class<T> mpi) {
        return (T) this.refers.computeIfAbsent(mpi.getName(), x -> ServiceProxy.proxy(mpi));
    }

    @Override
    public Object invoke(String urn, Map<String, Object> param) {
        try {
            return invoker.invoke(urn, param);
        } catch (RuntimeException | Error e) {
            throw e;
        } catch (Throwable e) {
            throw new MeshException(e);
        }
    }

    @Override
    public Object invoke(String urn, Object param) {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        return invoke(urn, codec.decode(codec.encode(param), Types.MapObject));
    }

}
