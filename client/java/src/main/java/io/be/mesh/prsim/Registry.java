/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.macro.Index;
import io.be.mesh.macro.MPI;
import io.be.mesh.macro.SPI;
import io.be.mesh.struct.Registration;

import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Registry {

    @MPI("mesh.registry.put")
    void register(@Index(0) Registration registration);

    @MPI("mesh.registry.puts")
    void registers(@Index(0) List<Registration> registrations);

    @MPI("mesh.registry.remove")
    void unregister(@Index(0) Registration registration);

    @MPI("mesh.registry.export")
    List<Registration> export(@Index(0) String kind);
}
