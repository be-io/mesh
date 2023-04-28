/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.Captcha;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Graphics {

    /**
     * Apply a graphics captcha.
     */
    @MPI("mesh.graphics.captcha.apply")
    Captcha captcha(String kind, Map<String, String> features);

    /**
     * Verify a graphics captcha value.
     */
    @MPI("mesh.graphics.captcha.verify")
    boolean verify(String mno, String value);
}
