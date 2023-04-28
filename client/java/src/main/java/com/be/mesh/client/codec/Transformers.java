/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.codec;

import com.be.mesh.client.annotate.Format;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.tool.Tool;

import java.io.IOException;
import java.lang.reflect.Field;
import java.time.*;
import java.time.format.DateTimeFormatter;
import java.util.Date;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
public class Transformers {

    private Transformers() {

    }

    static final DateTimeFormatter DATE = DateTimeFormatter.ofPattern("yyyy-MM-dd");
    static final DateTimeFormatter TIME = DateTimeFormatter.ofPattern("HH:mm:ss");
    static final DateTimeFormatter DATE_TIME = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
    private static final Map<String, DateTimeFormatter> formatters = new ConcurrentHashMap<>(3);

    @SPI("instant")
    public static final class DateTransformer implements Transformer<Date> {

        @Override
        public void form(Writer writer, Field field, Date value) throws IOException {
            if (null == value) {
                writer.writeNull();
                return;
            }
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                writer.write(value.getTime());
                return;
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            writer.write(formatter.format(OffsetDateTime.ofInstant(value.toInstant(), OffsetDateTime.now().getOffset())));
        }

        @Override
        public Date from(Reader reader, Field field) throws IOException {
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                return LocalDateTimeTransformer.compatibleWithString(reader.readString()).map(Date::new).orElse(null);
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            return Optional.ofNullable(reader.readString()).map(x -> LocalDateTime.parse(x, formatter)).map(z -> new Date(z.toInstant(OffsetDateTime.now().getOffset()).toEpochMilli())).orElse(null);
        }

        @Override
        public boolean matches(Field field) {
            return field.getType() == Date.class;
        }
    }

    @SPI("datetime")
    public static final class LocalDateTimeTransformer implements Transformer<LocalDateTime> {

        @Override
        public void form(Writer writer, Field field, LocalDateTime value) throws IOException {
            if (null == value) {
                writer.writeNull();
                return;
            }
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                writer.write(value.toInstant(OffsetDateTime.now().getOffset()).toEpochMilli());
                return;
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            writer.write(formatter.format(value));
        }

        @Override
        public LocalDateTime from(Reader reader, Field field) throws IOException {
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                return from(reader.readString()).orElse(null);
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            return Optional.ofNullable(reader.readString()).map(x -> LocalDateTime.parse(x, formatter)).orElse(null);
        }

        @Override
        public boolean matches(Field field) {
            return field.getType() == LocalDateTime.class;
        }

        static Optional<LocalDateTime> from(String number) {
            return compatibleWithString(number).map(x -> LocalDateTime.ofInstant(Instant.ofEpochMilli(x), OffsetDateTime.now().getOffset()));
        }

        static Optional<Long> compatibleWithString(String value) {
            if (Tool.isNumeric(value)) {
                return Optional.of(Long.parseLong(value));
            }
            // yyyy-MM-dd
            if (Tool.required(value) && value.length() == 10) {
                return Optional.of(LocalDateTime.of(LocalDate.parse(value, DATE), LocalTime.MIN).toInstant(OffsetDateTime.now().getOffset()).toEpochMilli());
            }
            // HH:mm:ss
            if (Tool.required(value) && value.length() == 10) {
                return Optional.of(LocalDateTime.of(LocalDate.now(), LocalTime.parse(value, TIME)).toInstant(OffsetDateTime.now().getOffset()).toEpochMilli());
            }
            // yyyy-MM-dd HH:mm:ss
            if (Tool.required(value) && value.length() == 10) {
                return Optional.of(LocalDateTime.parse(value, DATE_TIME).toInstant(OffsetDateTime.now().getOffset()).toEpochMilli());
            }
            return Optional.empty();
        }
    }

    @SPI("date")
    public static final class LocalDateTransformer implements Transformer<LocalDate> {

        @Override
        public void form(Writer writer, Field field, LocalDate value) throws IOException {
            if (null == value) {
                writer.writeNull();
                return;
            }
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                writer.write(LocalDateTime.of(value, LocalTime.MIN).toInstant(OffsetDateTime.now().getOffset()).toEpochMilli());
                return;
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            writer.write(formatter.format(value));
        }

        @Override
        public LocalDate from(Reader reader, Field field) throws IOException {
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                return LocalDateTimeTransformer.from(reader.readString()).map(LocalDateTime::toLocalDate).orElse(null);
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            return Optional.ofNullable(reader.readString()).map(x -> LocalDate.parse(x, formatter)).orElse(null);
        }

        @Override
        public boolean matches(Field field) {
            return field.getType() == LocalDate.class;
        }
    }

    @SPI("time")
    public static final class LocalTimeTransformer implements Transformer<LocalTime> {

        @Override
        public void form(Writer writer, Field field, LocalTime value) throws IOException {
            if (null == value) {
                writer.writeNull();
                return;
            }
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                writer.write(LocalDateTime.of(LocalDate.now(), value).toInstant(OffsetDateTime.now().getOffset()).toEpochMilli());
                return;
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            writer.write(formatter.format(value));
        }

        @Override
        public LocalTime from(Reader reader, Field field) throws IOException {
            Format format = field.getAnnotation(Format.class);
            if (null == format || Tool.optional(format.pattern())) {
                return LocalDateTimeTransformer.from(reader.readString()).map(LocalDateTime::toLocalTime).orElse(null);
            }
            DateTimeFormatter formatter = formatters.computeIfAbsent(format.pattern(), DateTimeFormatter::ofPattern);
            return Optional.ofNullable(reader.readString()).map(x -> LocalTime.parse(x, formatter)).orElse(null);
        }

        @Override
        public boolean matches(Field field) {
            return field.getType() == LocalTime.class;
        }
    }

}
