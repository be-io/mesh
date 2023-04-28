/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.mpc.ServiceProxy;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class SequenceTest {

    @Test
    public void nextTest() {
        log.info("{}", ServiceProxy.proxy(Sequence.class).next("INIT", 6));
    }
}
