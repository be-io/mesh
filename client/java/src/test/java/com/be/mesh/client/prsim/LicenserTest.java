/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.ServiceProxy;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class LicenserTest {

    @Test
    public void testVerify() {
        Licenser licenser = ServiceProxy.proxy(Licenser.class);
        log.info("{}", licenser.verify());
    }

    @Test
    public void testExplain() {
        Licenser licenser = ServiceProxy.proxy(Licenser.class);
        log.info("{}", licenser.explain());
    }
}
