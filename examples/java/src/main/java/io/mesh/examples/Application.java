/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package io.mesh.examples;

import com.jiumi.mesh.schema.macro.Version;
import com.jiumi.mesh.mpc.ServiceLoader;
import com.jiumi.mesh.schema.CompilePlugin;

/**
 * @author coyzeng@gmail.com
 */
@Version("mesh")
public class Application {

    public static void main(String[] args) {
        System.out.println(ServiceLoader.load(CompilePlugin.class).map().size());
    }
}
