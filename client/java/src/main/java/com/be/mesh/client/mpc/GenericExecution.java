/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.struct.Reference;

import java.lang.reflect.Type;

/**
 * @author coyzeng@gmail.com
 */
public class GenericExecution implements Execution<Reference> {

    private final Inspector inspector;
    private final Reference reference;

    public GenericExecution(URN urn) {
        this.reference = new Reference();
        this.reference.setUrn(urn.toString());
        this.reference.setNamespace("");
        this.reference.setName(urn.getName());
        this.reference.setVersion(urn.getFlag().getVersion());
        this.reference.setProto(MeshFlag.ofProto(urn.getFlag().getProto()).getName());
        this.reference.setCodec(MeshFlag.ofCodec(urn.getFlag().getCodec()).getName());
        this.reference.setFlags(0);
        this.reference.setTimeout(getTimeout());
        this.reference.setRetries(5);
        this.reference.setNode(urn.getNodeId());
        this.reference.setInst("");
        this.reference.setZone(urn.getFlag().getZone());
        this.reference.setCluster(urn.getFlag().getCluster());
        this.reference.setCell(urn.getFlag().getCell());
        this.reference.setGroup(urn.getFlag().getCell());
        this.reference.setAddress(urn.getFlag().getAddress());
        this.inspector = new GenericInspector(urn.toString());
    }

    @Override
    public Reference schema() {
        return this.reference;
    }

    @Override
    public Inspector inspect() {
        return this.inspector;
    }

    @Override
    public <I extends Parameters> Types<I> intype() {
        return Types.of((Type) GenericParameters.class);
    }

    @Override
    public <O extends Returns> Types<O> retype() {
        return Types.of((Type) GenericReturns.class);
    }

    @Override
    public Parameters inflect() {
        return new GenericParameters();
    }

    @Override
    public Returns reflect() {
        return new GenericReturns();
    }

    @Override
    public Object invoke(Invocation invocation) throws Throwable {
        return this.inspector.invoke(invocation.getProxy(), invocation.getParameters().arguments());
    }

    /**
     * set timeout, unit:ms
     * default: 10s
     *
     * @return
     */
    private int getTimeout() {
        String meshTimeoutStr = Mesh.context().getAttachments().get("Mesh-Timeout");
        if (meshTimeoutStr != null) {
            return Integer.parseInt(meshTimeoutStr);
        }
        return 10000;
    }

}
