/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.plugin;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.schema.macro.Digest;
import com.be.mesh.client.schema.context.CompileContext;
import com.be.mesh.client.schema.CompilePlugin;

/**
 * @author coyzeng@gmail.com
 */
@SPI(value = "digest", meta = Digest.class)
public class DigestPlugin implements CompilePlugin {

    @Override
    public void proceed(CompileContext context) {

    }
}
