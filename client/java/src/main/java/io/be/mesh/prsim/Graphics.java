/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.macro.MPI;
import io.be.mesh.macro.SPI;
import io.be.mesh.struct.Captcha;

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
