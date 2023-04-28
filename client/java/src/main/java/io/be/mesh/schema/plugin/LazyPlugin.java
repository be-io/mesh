/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.schema.plugin;

import io.be.mesh.macro.SPI;
import io.be.mesh.schema.macro.Lazy;
import io.be.mesh.schema.context.CompileContext;
import io.be.mesh.schema.CompilePlugin;

/**
 * @author coyzeng@gmail.com
 */
@SPI(value = "lazy", meta = Lazy.class)
public class LazyPlugin implements CompilePlugin {

    @Override
    public void proceed(CompileContext context) {
        //
    }
}
