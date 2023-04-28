/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.boost;

import ch.qos.logback.classic.PatternLayout;
import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.pattern.CompositeConverter;
import ch.qos.logback.core.spi.ContextAwareBase;
import com.be.mesh.client.tool.Envs;
import com.be.mesh.client.tool.Tool;
import lombok.AllArgsConstructor;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

/**
 * @author coyzeng@gmail.com
 */
public class RollingFileAppender<E> extends ch.qos.logback.core.rolling.RollingFileAppender<E> {

    private final Filer filer = new Filer(this);

    @Override
    public void setFile(String file) {
        String path = filer.getPath();
        addInfo(String.format("Rolling with %s", path));
        super.setFile(path);
    }

    public static final class SizeAndTimeBasedRollingPolicy<E> extends ch.qos.logback.core.rolling.SizeAndTimeBasedRollingPolicy<E> {

        private final Filer filer = new Filer(this);

        @Override
        public void setFileNamePattern(String fnp) {
            String path = filer.getTimingPathPattern();
            addInfo(String.format("Rolling with pattern %s", path));
            super.setFileNamePattern(path);
        }
    }


    @AllArgsConstructor
    private static final class Filer {
        private final ContextAwareBase base;

        private String getPath() {
            try {
                String pmp = getPromTailPath();
                if (writable(pmp)) {
                    return pmp;
                }
                return getRecommendPath();
            } finally {
                Envs.release();
            }
        }

        private String getTimingPathPattern() {
            try {
                String pmp = getPromTailPath();
                if (writable(pmp)) {
                    return String.join(File.separator, Paths.get(pmp).getParent().toString(), "app-%d{yyyy-MM}.%i.log");
                }
                return String.join(File.separator, Paths.get(getRecommendPath()).getParent().toString(), Tool.MESH_NAME.get() + "-%d{yyyy-MM}.%i.log");
            } finally {
                Envs.release();
            }
        }

        private boolean writable(String path) {
            try {
                Path p = Paths.get(path);
                if (Files.exists(p) && Files.isWritable(p)) {
                    return true;
                }
                if (!Files.exists(p.getParent())) {
                    Files.createDirectories(p.getParent());
                }
                if (!Files.exists(p)) {
                    Files.createFile(p);
                }
                return Files.isWritable(p);
            } catch (IOException e) {
                base.addWarn(String.format("Check %s, %s", path, e.getMessage()));
                return false;
            }
        }

        private String getPromTailPath() {
            String home = Tool.getProperty(String.join(File.separator, "", "var", "log", "be"), "LOG_HOME", "LOG.HOME", "log_home", "log.home");
            String name = Tool.getProperty(Tool.MESH_NAME.get(), "APP_NAME");
            return String.join(File.separator, home, name, "app.log");
        }

        private String getRecommendPath() {
            return String.join(File.separator, System.getProperty("user.home"), "logs", String.format("%s.log", Tool.MESH_NAME.get()));
        }
    }

    public static final class Escaper extends CompositeConverter<ILoggingEvent> {

        private static final String[] REPLACEMENT_CHARS = new String[128];

        static {
            for (int i = 0; i <= 31; ++i) {
                REPLACEMENT_CHARS[i] = String.format("\\u%04x", i);
            }
            REPLACEMENT_CHARS[34] = "\\\"";
            REPLACEMENT_CHARS[92] = "\\\\";
            REPLACEMENT_CHARS[9] = "\\t";
            REPLACEMENT_CHARS[8] = "\\b";
            REPLACEMENT_CHARS[10] = "\\n";
            REPLACEMENT_CHARS[13] = "\\r";
            REPLACEMENT_CHARS[12] = "\\f";
        }

        @Override
        protected String transform(ILoggingEvent event, String in) {
            StringBuilder buff = new StringBuilder(in.length());
            buff.append((char) 34);
            for (int index = 0; index < in.length(); ++index) {
                char character = in.charAt(index);
                if (character < 128) {
                    if (null == REPLACEMENT_CHARS[character]) {
                        buff.append(character);
                    } else {
                        buff.append(REPLACEMENT_CHARS[character]);
                    }
                    continue;
                }
                if (character == 8232) {
                    buff.append("\\u2028");
                    continue;
                }
                if (character == 8233) {
                    buff.append("\\u2029");
                    continue;
                }
                buff.append(character);
            }
            buff.append((char) 34);
            return buff.toString();
        }
    }

    static {
        PatternLayout.DEFAULT_CONVERTER_MAP.put("escape", Escaper.class.getName());
    }
}
