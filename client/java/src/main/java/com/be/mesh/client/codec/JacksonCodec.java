/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.codec;

import com.be.mesh.client.annotate.Format;
import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.CompatibleException;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ParameterizedTypes;
import com.be.mesh.client.mpc.PatternParameterizedType;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.tool.Once;
import com.be.mesh.client.tool.Tool;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.Version;
import com.fasterxml.jackson.databind.*;
import com.fasterxml.jackson.databind.cfg.MapperConfig;
import com.fasterxml.jackson.databind.introspect.AnnotatedField;
import com.fasterxml.jackson.databind.module.SimpleDeserializers;
import com.fasterxml.jackson.databind.module.SimpleSerializers;
import lombok.SneakyThrows;

import java.io.IOException;
import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.nio.ByteBuffer;
import java.time.*;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.regex.Pattern;

/**
 * @author coyzeng@gmail.com
 */
@SPI(Codec.JACKSON)
public class JacksonCodec implements Codec {

    private static final Map<String, JavaType> types = new ConcurrentHashMap<>();
    private static final Once<ObjectMapper> mapper = Once.with(() -> new ObjectMapper()
            .registerModule(new JacksonModule())
            .setSerializationInclusion(JsonInclude.Include.ALWAYS)
            .setDefaultPropertyInclusion(JsonInclude.Include.ALWAYS)
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false));
    private static final Pattern numberPattern = Pattern.compile("[0-9]*");

    @Override
    public ByteBuffer encode(Object value) {
        try {
            return ByteBuffer.wrap(mapper.get().writeValueAsBytes(value));
        } catch (JsonProcessingException e) {
            throw new CompatibleException(e);
        }
    }

    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        try {
            return mapper.get().readValue(buffer.array(), types.computeIfAbsent(type.toString(), key -> mapper.get().getTypeFactory().constructType(type)));
        } catch (IOException e) {
            throw new CompatibleException(e);
        }
    }

    static final class JacksonModule extends com.fasterxml.jackson.databind.Module {

        @Override
        public String getModuleName() {
            return "mesh-module";
        }

        @Override
        public Version version() {
            return Version.unknownVersion();
        }

        @Override
        public void setupModule(SetupContext context) {
            NamingStrategy namingStrategy = new NamingStrategy();
            SimpleDeserializers deserializers = new SimpleDeserializers();
            deserializers.addDeserializer(LocalDate.class, new JsonDeserializer<LocalDate>() {
                @Override
                public LocalDate deserialize(JsonParser reader, DeserializationContext ctx) throws IOException {
                    String text = reader.getText();
                    if (null == text) {
                        return null;
                    }
                    if (hasFormatAnnotation(reader) && isNumber(text)) {
                        long longTime = Long.parseLong(text);
                        return LocalDateTime.ofInstant(Instant.ofEpochMilli(longTime), ZoneId.systemDefault()).toLocalDate();
                    }
                    return LocalDate.parse(text, Transformers.DATE);
                }
            });
            deserializers.addDeserializer(LocalTime.class, new JsonDeserializer<LocalTime>() {
                @Override
                public LocalTime deserialize(JsonParser reader, DeserializationContext ctx) throws IOException {
                    String text = reader.getText();
                    if (null == text) {
                        return null;
                    }
                    if (hasFormatAnnotation(reader) && isNumber(text)) {
                        long longTime = Long.parseLong(text);
                        return LocalDateTime.ofInstant(Instant.ofEpochMilli(longTime), ZoneId.systemDefault()).toLocalTime();
                    }
                    return LocalTime.parse(text, Transformers.TIME);
                }
            });
            deserializers.addDeserializer(LocalDateTime.class, new JsonDeserializer<LocalDateTime>() {
                @Override
                public LocalDateTime deserialize(JsonParser reader, DeserializationContext ctx) throws IOException {
                    String text = reader.getText();
                    if (null == text) {
                        return null;
                    }
                    if (hasFormatAnnotation(reader) && isNumber(text)) {
                        long longTime = Long.parseLong(text);
                        return LocalDateTime.ofInstant(Instant.ofEpochMilli(longTime), ZoneId.systemDefault());
                    }
                    return LocalDateTime.parse(text, Transformers.DATE_TIME);
                }
            });
            deserializers.addDeserializer(Duration.class, new JsonDeserializer<Duration>() {
                @Override
                public Duration deserialize(JsonParser reader, DeserializationContext ctx) throws IOException {
                    return Duration.ofMillis(reader.nextIntValue(0));
                }
            });
            deserializers.addDeserializer(Types.class, new JsonDeserializer<Types<?>>() {
                @Override
                public Types<?> deserialize(JsonParser reader, DeserializationContext ctx) throws IOException {
                    String text = reader.getText();
                    if (null == text) {
                        return null;
                    }
                    String struct = Tool.decompress(text);
                    ParameterizedTypes types = mapper.get().readValue(struct, ParameterizedTypes.class);
                    return Types.of(PatternParameterizedType.make(types));
                }
            });
            SimpleSerializers serializers = new SimpleSerializers();
            serializers.addSerializer(LocalDate.class, new JsonSerializer<LocalDate>() {
                @Override
                public void serialize(LocalDate value, JsonGenerator writer, SerializerProvider provider) throws IOException {

                    if (null != value) {
                        if (hasFormatAnnotation(writer)) {
                            writer.writeNumber(value.toEpochDay());
                        } else {
                            writer.writeString(Transformers.DATE.format(value));
                        }
                    } else {
                        writer.writeNull();
                    }
                }
            });
            serializers.addSerializer(LocalTime.class, new JsonSerializer<LocalTime>() {
                @Override
                public void serialize(LocalTime value, JsonGenerator writer, SerializerProvider provider) throws IOException {
                    if (null != value) {
                        if (hasFormatAnnotation(writer)) {
                            writer.writeNumber(value.toNanoOfDay());
                        } else {
                            writer.writeString(Transformers.TIME.format(value));
                        }
                    } else {
                        writer.writeNull();
                    }
                }
            });
            serializers.addSerializer(LocalDateTime.class, new JsonSerializer<LocalDateTime>() {
                @Override
                public void serialize(LocalDateTime value, JsonGenerator writer, SerializerProvider provider) throws IOException {
                    if (null != value) {
                        if (hasFormatAnnotation(writer)) {
                            writer.writeNumber(value.atZone(ZoneId.systemDefault()).toInstant().toEpochMilli());
                        } else {
                            writer.writeString(Transformers.DATE_TIME.format(value));
                        }
                    } else {
                        writer.writeNull();
                    }
                }
            });
            serializers.addSerializer(Duration.class, new JsonSerializer<Duration>() {
                @Override
                public void serialize(Duration value, JsonGenerator writer, SerializerProvider provider) throws IOException {

                    if (null != value) {
                        writer.writeNumber(value.toMillis());
                    } else {
                        writer.writeNumber(0);
                    }
                }
            });
            serializers.addSerializer(Types.class, new JsonSerializer<Types>() {
                @Override
                public void serialize(Types value, JsonGenerator writer, SerializerProvider provider) throws IOException {
                    if (null != value) {
                        writer.writeString(Tool.compress(mapper.get().writeValueAsString(new ParameterizedTypes(value))));
                    } else {
                        writer.writeNull();
                    }
                }
            });
            context.setNamingStrategy(namingStrategy);
            context.addSerializers(serializers);
            context.addDeserializers(deserializers);
        }
    }

    static class NamingStrategy extends PropertyNamingStrategy {

        @Override
        public String nameForField(MapperConfig<?> config, AnnotatedField field, String defaultName) {
            Index index = field.getAnnotation(Index.class);
            if (null != index && Tool.required(index.name())) {
                return index.name();
            }
            Format format = field.getAnnotation(Format.class);
            if (null != format && Tool.required(format.value())) {
                return format.value();
            }
            return defaultName;
        }
    }


    @SneakyThrows
    static boolean hasFormatAnnotation(JsonParser jsonParser) {
        Field field = findField(jsonParser.getCurrentName(), jsonParser.getCurrentValue().getClass());
        if (field == null) {
            return false;
        }
        if (field.isAnnotationPresent(Format.class)) {
            return true;
        }
        return field.getAnnotation(Format.class) != null;
    }

    @SneakyThrows
    static boolean hasFormatAnnotation(JsonGenerator writer) {
        Field field = findField(writer.getOutputContext().getCurrentName(), writer.getCurrentValue().getClass());
        if (field == null) {
            return false;
        }
        if (field.isAnnotationPresent(Format.class)) {
            return true;
        }
        return field.getAnnotation(Format.class) != null;
    }

    static boolean isNumber(String str) {
        if (str.length() > 0) {
            return numberPattern.matcher(str).matches();
        }
        return false;
    }

    static Field findField(String name, Class<?> c) {
        for (; c != null; c = c.getSuperclass()) {
            for (Field field : c.getDeclaredFields()) {
                if (Modifier.isStatic(field.getModifiers())) {
                    continue;
                }
                if (field.getName().equals(name)) {
                    return field;
                }
            }
        }
        return null;
    }
}
