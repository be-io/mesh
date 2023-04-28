/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.struct.Service;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Eden {

    /**
     * Define the reference object.
     *
     * @param reference Object reference type.
     * @param mpi       Meta custom annotation.
     * @return Object
     */
    Object define(MPI mpi, Class<?> reference);

    /**
     * Refer the service reference by method.
     *
     * @param mpi       Meta custom annotation.
     * @param reference Service reference type.
     * @param inspector Service method.
     * @return Service reference.
     */
    Execution<Reference> refer(MPI mpi, Class<?> reference, Inspector inspector);

    /**
     * Store the service object.
     *
     * @param type    Object type.
     * @param service Service object.
     */
    void store(Class<?> type, Object service);

    /**
     * Infer the reference service by domain.
     *
     * @param urn Service urn.
     * @return Service reference.
     */
    Execution<Service> infer(String urn);

    /**
     * Get all reference types.
     *
     * @return All reference types.
     */
    List<Class<?>> referTypes();

    /**
     * Get all service types.
     *
     * @return All service types.
     */
    List<Class<?>> inferTypes();
}
