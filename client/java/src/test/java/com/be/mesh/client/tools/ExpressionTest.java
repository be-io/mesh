/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tools;

import com.be.mesh.client.schema.plugin.VersionPlugin;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class ExpressionTest {

    @Test
    public void testVersionMatcher() {
        log.info("release/1.5.6.1.1".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("release-1.5.6.1.1".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("release/1.5.6.1".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("release-1.5.6.1".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("release/1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("release-1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("feature/1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("feature-1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("dev/1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("dev-1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("master".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("dev".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("release".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("dev-1.5.7".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("v1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
        log.info("1.5.6".replaceAll(VersionPlugin.PATTERN, "$1"));
    }
}
