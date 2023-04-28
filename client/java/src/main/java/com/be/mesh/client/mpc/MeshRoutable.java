/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.prsim.Routable;
import com.be.mesh.client.struct.Principal;
import com.be.mesh.client.tool.Tool;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;

import java.lang.reflect.*;
import java.util.*;
import java.util.stream.Stream;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@Setter
public class MeshRoutable<T> implements Routable<T> {

    private final T reference;
    private Map<String, String> attachments;
    private String address;

    public MeshRoutable(T reference) {
        this.reference = reference;
    }

    @Override
    public Routable<T> with(String key, String value) {
        Map<String, String> kv = new HashMap<>(1);
        kv.put(key, value);
        return with(kv);
    }

    @Override
    public Routable<T> with(Map<String, String> attachments) {
        Map<String, String> kvs = new HashMap<>();
        if (Tool.required(this.attachments)) {
            kvs.putAll(this.attachments);
        }
        if (Tool.required(attachments)) {
            kvs.putAll(attachments);
        }
        MeshRoutable<T> stream = this.cp();
        stream.setAttachments(kvs);
        return stream;
    }

    @Override
    public Routable<T> withAddress(String address) {
        MeshRoutable<T> stream = this.cp();
        stream.setAddress(address);
        return stream;
    }

    @Override
    public T local() {
        return reference;
    }

    @SuppressWarnings("unchecked")
    @Override
    public T any(Principal principal) {
        if (Tool.optional(principal.getNodeId()) && Tool.optional(principal.getInstId())) {
            throw new MeshException("Route key both cant be empty.");
        }
        Class<?> kind = reference.getClass();
        return (T) Proxy.newProxyInstance(kind.getClassLoader(), kind.getInterfaces(), new RouteInvocationHandler(reference, principal, attachments, address));
    }

    @Override
    public T any(String instId) {
        return this.any(Principal.ofInstId(instId));
    }

    @Override
    public Stream<T> many(Principal... principals) {
        if (Tool.optional(principals)) {
            return Stream.empty();
        }
        List<T> references = new ArrayList<>(principals.length);
        for (Principal principal : principals) {
            references.add(any(principal));
        }
        return references.stream();
    }

    @Override
    public Stream<T> many(String... instIds) {
        if (Tool.optional(instIds)) {
            return Stream.empty();
        }
        Principal[] principals = new Principal[instIds.length];
        for (int index = 0; index < instIds.length; index++) {
            principals[index] = Principal.ofInstId(instIds[index]);
        }
        return this.many(principals);
    }

    private MeshRoutable<T> cp() {
        MeshRoutable<T> stream = new MeshRoutable<>(this.reference);
        stream.setAttachments(this.attachments);
        stream.setAddress(this.address);
        return stream;
    }

    private static final class RouteInvocationHandler implements InvocationHandler {

        private final Object reference;
        private final Principal principal;
        private final Map<String, String> attachments;
        private final String address;

        private RouteInvocationHandler(Object reference, Principal principal, Map<String, String> attachments, String address) {
            this.reference = reference;
            this.principal = principal;
            if (null != attachments) {
                this.attachments = new HashMap<>(attachments);
            } else {
                this.attachments = Collections.emptyMap();
            }
            this.address = address;
        }

        @Override
        public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
            return Mesh.contextSafe(() -> {
                Map<String, String> overrides = getOverrideAttachments();
                try {
                    if (Tool.required(attachments)) {
                        Mesh.context().getAttachments().putAll(attachments);
                    }
                    if (Tool.required(address)) {
                        Mesh.context().setAttribute(Mesh.REMOTE, address);
                    }
                    Mesh.context().getPrincipals().push(principal);
                    return method.invoke(reference, args);
                } catch (InvocationTargetException | UndeclaredThrowableException e) {
                    throw Tool.destructor(e.getCause(), method.getExceptionTypes());
                } catch (Throwable e) {
                    throw Tool.destructor(e, method.getExceptionTypes());
                } finally {
                    Mesh.context().getAttachments().putAll(overrides);
                    Mesh.context().getPrincipals().poll();
                }
            });
        }

        private Map<String, String> getOverrideAttachments() {
            if (Tool.optional(attachments)) {
                return Collections.emptyMap();
            }
            Map<String, String> overrides = new HashMap<>(attachments.size());
            attachments.forEach((key, value) -> overrides.put(key, Optional.ofNullable(Mesh.context().getAttachments().get(key)).orElse("")));
            return overrides;
        }
    }
}
