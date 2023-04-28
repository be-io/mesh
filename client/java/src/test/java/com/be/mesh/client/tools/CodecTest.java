/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tools;

import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import com.google.common.collect.ImmutableMap;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class CodecTest {

    private String ID = "1";

    private String I = "1";

    private String IDsDF = "1";

    private String sABC = "1";

    private Map<String, String> x = ImmutableMap.of("sABC", "2", "IDsDF", "2", "ID", "2");

    @Test
    public void testCodec() {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        String text = codec.encodeString(new CodecTest());
        log.info(text);
        log.info(codec.decodeString(text, Types.of(CodecTest.class)).ID);
        log.info(codec.decodeString(text, Types.of(CodecTest.class)).I);
        log.info(codec.decodeString(text, Types.of(CodecTest.class)).IDsDF);
        log.info(codec.decodeString(text, Types.of(CodecTest.class)).sABC);
        log.info("{}", codec.decodeString(text, Types.of(CodecTest.class)).x);
    }
}
