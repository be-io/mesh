/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import lombok.AllArgsConstructor;

import javax.annotation.processing.ProcessingEnvironment;
import javax.tools.Diagnostic;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CopyOnWriteArrayList;

/**
 * @author coyzeng@gmail.com
 */
public interface CompilePrinter {

    default void info(String format, Object... args) {
        Printers.write("INFO", format, args);
    }

    default void info(Throwable e, String format, Object... args) {
        Printers.write("INFO", e, format, args);
    }

    default void warn(String format, Object... args) {
        Printers.write("WARN", format, args);
    }

    default void warn(Throwable e, String format, Object... args) {
        Printers.write("WARN", e, format, args);
    }

    default void error(String format, Object... args) {
        Printers.write("ERROR", format, args);
    }

    default void error(Throwable e, String format, Object... args) {
        Printers.write("ERROR", e, format, args);
    }

    default void debug(String format, Object... args) {
        Printers.write("DEBUG", format, args);
    }

    default void debug(Throwable e, String format, Object... args) {
        Printers.write("DEBUG", e, format, args);
    }

    interface Printer {
        void write(String level, String msg, Throwable e);
    }

    class STDPrinter implements Printer {

        @Override
        public void write(String level, String msg, Throwable e) {
            System.out.printf("[%s] %s%n", level, msg);
            if (null != e) {
                e.printStackTrace();
            }
        }
    }

    @AllArgsConstructor
    class JavacPrinter implements Printer {
        private final ProcessingEnvironment environment;

        @Override
        public void write(String level, String msg, Throwable e) {
            switch (level) {
                case "ERROR":
                    this.environment.getMessager().printMessage(Diagnostic.Kind.ERROR, msg);
                    break;
                case "WARN":
                    this.environment.getMessager().printMessage(Diagnostic.Kind.WARNING, msg);
                    break;
                case "INFO":
                case "DEBUG":
                default:
                    this.environment.getMessager().printMessage(Diagnostic.Kind.NOTE, msg);
            }
        }
    }

    class Printers {
        public static final List<Printer> WRITERS = new CopyOnWriteArrayList<>();

        static {
            WRITERS.add(new STDPrinter());
        }

        public static void write(String level, Throwable e, String format, Object... args) {
            WRITERS.forEach(w -> w.write(level, String.format(Optional.ofNullable(format).orElse(""), args), e));
        }

        public static void write(String level, String format, Object... args) {
            WRITERS.forEach(w -> w.write(level, String.format(Optional.ofNullable(format).orElse(""), args), null));
        }
    }


}
