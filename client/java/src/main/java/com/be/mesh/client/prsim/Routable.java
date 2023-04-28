/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.MeshRoutable;
import com.be.mesh.client.struct.Principal;

import java.util.Map;
import java.util.stream.Stream;

/**
 * @author coyzeng@gmail.com
 */
public interface Routable<T> {

    /**
     * Route with attachments.
     *
     * @param key   attachment key
     * @param value attachment value
     * @return self
     */
    Routable<T> with(String key, String value);

    /**
     * Route with attachments.
     *
     * @param attachments attachments
     * @return self
     */
    Routable<T> with(Map<String, String> attachments);

    /**
     * Route with mesh address.
     *
     * @param address mesh address
     * @return routable reference
     */
    Routable<T> withAddress(String address);

    /**
     * Invoke the service in local network.
     *
     * @return Local invoker.
     */
    T local();

    /**
     * Invoke the service in a network, it may be local or others.
     *
     * @param principal Network principal.
     * @return Service invoker.
     */
    T any(Principal principal);

    /**
     * Invoke the service in a network, it may be local or others.
     *
     * @param instId Network principal of instId.
     * @return Service invoker.
     */
    T any(String instId);

    /**
     * Invoke the service in many network, it may be local or others. Broadcast mode.
     *
     * @param principals Network principals.
     * @return Service invoker.
     */
    Stream<T> many(Principal... principals);

    /**
     * Invoke the service in many network, it may be local or others. Broadcast mode.
     *
     * @param instIds Network principals.
     * @return Service invoker.
     */
    Stream<T> many(String... instIds);

    /**
     * Wrap a service with streamable ability.
     *
     * @param reference Service reference.
     * @param <T>       Service type.
     * @return Streamable program interface.
     */
    static <T> Routable<T> of(T reference) {
        return new MeshRoutable<>(reference);
    }

}
