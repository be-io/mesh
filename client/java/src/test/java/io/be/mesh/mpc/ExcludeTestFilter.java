/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.mpc;

import io.be.mesh.macro.SPI;

/**
 * @author coyzeng@gmail.com
 */
@SPI(value = "exclude-test-session", pattern = Filter.CONSUMER, priority = Integer.MIN_VALUE + 2, exclude = {"mesh", "mesh.net.edge"})
public class ExcludeTestFilter implements Filter {

    @Override
    public Object invoke(Invoker<?> invoker, Invocation invocation) throws Throwable {
        return invoker.invoke(invocation);
    }
}
