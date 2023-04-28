/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.ServiceProxy;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class DocumentTest {

    @Test
    public void testExport() {
        Mesh.contextSafeUncheck(() -> {
            Mesh.context().setAttribute(Mesh.REMOTE_NAME, "janus");
            Mesh.context().setAttribute(Mesh.REMOTE, "10.99.1.33:570");
            Builtin documentor = ServiceProxy.proxy(Builtin.class);
            String doc = documentor.doc("janus", "json");
            log.info(doc);
        });
    }
}
