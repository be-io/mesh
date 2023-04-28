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

import java.util.Map;

/**
 * Sample crypt interface in multi node network.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Cryptor {

    /**
     * Encrypt binary to encrypted binary.
     */
    @MPI("mesh.crypt.encrypt")
    byte[] encrypt(@Index(0) byte[] buff, @Index(1) Map<String, String> features);

    /**
     * Decrypt binary to decrypted binary.
     */
    @MPI("mesh.crypt.decrypt")
    byte[] decrypt(@Index(0) byte[] buff, @Index(1) Map<String, String> features);

    /**
     * Hash compute the hash value.
     */
    @MPI("mesh.crypt.hash")
    byte[] hash(@Index(0) byte[] buff, @Index(1) Map<String, String> features);

    /**
     * Sign compute the signature value.
     */
    @MPI("mesh.crypt.sign")
    String sign(@Index(0) byte[] buff, @Index(1) Map<String, String> features);

    /**
     * Verify the signature value.
     */
    @MPI("mesh.crypt.verify")
    boolean verify(@Index(0) byte[] buff, @Index(1) Map<String, String> features);
}
