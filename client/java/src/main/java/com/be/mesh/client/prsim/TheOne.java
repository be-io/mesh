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
import com.be.mesh.client.struct.Profile;

/**
 * Configurator use everywhere, keep the one.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface TheOne {

    /**
     * Push the profile to bind data id.
     *
     * @param profile data profile
     */
    @MPI("mesh.theone.pain")
    void pain(Profile profile);

    /**
     * Load the profiles by data id.
     *
     * @param dataId data id
     * @return data profiles
     */
    @MPI("mesh.theone.gain")
    Profile gain(@Index(value = 0, name = "data_id") String dataId);

    /**
     * Watch the profile. Trigger when change event subscribe from mesh node.
     *
     * @param dataId  data id
     * @param watcher data watcher
     */
    @MPI("mesh.theone.follow")
    void follow(@Index(value = 0, name = "data_id") String dataId, Watcher watcher);

    /**
     * Watcher function.
     */
    @FunctionalInterface
    interface Watcher {

        /**
         * Trigger when change event subscribe from mesh node.
         *
         * @param profile data profile
         */
        void receive(Profile profile);

    }
}
