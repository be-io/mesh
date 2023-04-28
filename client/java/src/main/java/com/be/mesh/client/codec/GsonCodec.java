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
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.tool.Once;
import com.be.mesh.client.tool.Tool;
import com.google.gson.*;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.internal.*;
import com.google.gson.internal.bind.ReflectiveTypeAdapterFactory;
import com.google.gson.internal.bind.TreeTypeAdapter;
import com.google.gson.reflect.TypeToken;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonToken;
import com.google.gson.stream.JsonWriter;
import lombok.AllArgsConstructor;

import java.io.IOException;
import java.lang.reflect.Field;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.lang.reflect.TypeVariable;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;

/**
 * <pre>
 *         List<TypeAdapterFactory> factories = (List<TypeAdapterFactory>) Tool.getField(instance, "factories");
 *         ConstructorConstructor constructor = (ConstructorConstructor) Tool.getField(instance, "constructorConstructor");
 *         boolean complexMapKeySerialization = (boolean) Tool.getField(instance, "complexMapKeySerialization");
 *         List<TypeAdapterFactory> facts = new ArrayList<>();
 *         facts.add(new MapTypeAdapterFactory(constructor, complexMapKeySerialization));
 *         for (TypeAdapterFactory factory : factories) {
 *             if (factory.getClass() == com.google.gson.internal.bind.MapTypeAdapterFactory.class) {
 *                 continue;
 *             }
 *             facts.add(factory);
 *         }
 *         Tool.setField(instance, "factories", Collections.unmodifiableList(facts));
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@SPI(Codec.JSON)
@SuppressWarnings({"unchecked", "rawtypes", "deprecation"})
public class GsonCodec implements Codec, InstanceCreator<Object> {

    private static final Once<Gson> gson = Once.with(() -> Tool.uncheck(() -> {
        Gson instance = new Gson().newBuilder().
                registerTypeHierarchyAdapter(byte[].class, new ByteArrayAdapter()).
                registerTypeAdapter(Duration.class, new DurationAdapter()).
                registerTypeAdapterFactory(new TypesAdapter()).
                setFieldNamingStrategy(new FieldNamingPolicy()).
                setObjectToNumberStrategy(ToNumberPolicy.LONG_OR_DOUBLE).
                setNumberToNumberStrategy(ToNumberPolicy.LAZILY_PARSED_NUMBER).
                disableHtmlEscaping().
                create();
        List<TypeAdapterFactory> factories = (List<TypeAdapterFactory>) Tool.getField(instance, "factories");
        ConstructorConstructor constructor = (ConstructorConstructor) Tool.getField(instance, "constructorConstructor");
        List<TypeAdapterFactory> facts = new ArrayList<>();
        for (TypeAdapterFactory factory : factories) {
            if (factory.getClass() == ReflectiveTypeAdapterFactory.class) {
                continue;
            }
            facts.add(factory);
        }
        facts.add(new FieldAccessAdapterFactory(constructor, instance.fieldNamingStrategy(), instance.excluder()));
        Tool.setField(instance, "factories", Collections.unmodifiableList(facts));
        return instance;
    }));

    @Override
    public ByteBuffer encode(Object value) {
        return this.encode0(value, object -> ByteBuffer.wrap(gson.get().toJson(object).getBytes(StandardCharsets.UTF_8)));
    }

    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        return this.decode0(buffer, type, (x, y) -> gson.get().fromJson(new String(x.array(), StandardCharsets.UTF_8), y));
    }

    @Override
    public Object createInstance(Type type) {
        try {
            Class<?> t = resolveClass(type);
            if (Collection.class.isAssignableFrom(t)) {
                return new LinkedList<>();
            }
            if (Map.class.isAssignableFrom(t)) {
                return new LinkedHashMap<>();
            }
            return t.getDeclaredConstructor().newInstance();
        } catch (RuntimeException | Error e) {
            throw e;
        } catch (Exception e) {
            throw new MeshException(e);
        }
    }

    private Class<?> resolveClass(Type type) {
        if (type instanceof ParameterizedType) {
            return ((Class<?>) ((ParameterizedType) type).getRawType());
        }
        return ((Class<?>) type);
    }

    static class ByteArrayAdapter implements JsonSerializer<byte[]>, JsonDeserializer<byte[]> {

        @Override
        public byte[] deserialize(JsonElement value, Type type, JsonDeserializationContext context) throws JsonParseException {
            return Optional.ofNullable(value.getAsString()).map(x -> Base64.getDecoder().decode(x)).orElse(new byte[0]);
        }

        @Override
        public JsonElement serialize(byte[] value, Type type, JsonSerializationContext context) {
            return new JsonPrimitive(Base64.getEncoder().encodeToString(Optional.ofNullable(value).orElse(new byte[0])));
        }
    }

    static final class DurationAdapter extends TypeAdapter<Duration> {

        @Override
        public void write(JsonWriter out, Duration value) throws IOException {
            out.value(value.toMillis());
        }

        @Override
        public Duration read(JsonReader in) throws IOException {
            if (!in.hasNext()) {
                return Duration.ZERO;
            }
            return Duration.ofMillis(in.nextLong());
        }
    }

    static final class TypesAdapter extends TypeAdapter<Types<?>> implements TypeAdapterFactory {

        @Override
        public void write(JsonWriter out, Types<?> value) throws IOException {
            out.value(Tool.compress(gson.get().toJson(new ParameterizedTypes(value))));
        }

        @Override
        public Types<?> read(JsonReader in) throws IOException {
            if (!in.hasNext()) {
                return Types.MapObject;
            }
            String struct = Tool.decompress(in.nextString());
            ParameterizedTypes types = gson.get().fromJson(struct, ParameterizedTypes.class);
            return Types.of(PatternParameterizedType.make(types));
        }

        @SuppressWarnings("unchecked")
        @Override
        public <T> TypeAdapter<T> create(Gson gson, TypeToken<T> type) {
            if (type.getRawType() == Types.class) {
                return (TypeAdapter<T>) this;
            }
            return null;
        }
    }

    /**
     * Type adapter that reflects over the fields and methods of a class.
     */
    static final class FieldAccessAdapterFactory implements TypeAdapterFactory {
        private static final List<Class<?>> FS = Arrays.asList(LocalDateTime.class, LocalDate.class, LocalTime.class, Date.class, Object.class);
        private final ConstructorConstructor constructor;
        private final FieldNamingStrategy namingStrategy;
        private final Excluder excluder;

        public FieldAccessAdapterFactory(ConstructorConstructor constructor, FieldNamingStrategy namingStrategy, Excluder excluder) {
            this.constructor = constructor;
            this.namingStrategy = namingStrategy;
            this.excluder = excluder;
        }

        public boolean excludeField(Field f, boolean serialize) {
            return excludeField(f, serialize, excluder);
        }

        static boolean excludeField(Field f, boolean serialize, Excluder excluder) {
            return !excluder.excludeClass(f.getType(), serialize) && !excluder.excludeField(f, serialize);
        }

        /**
         * first element holds the default name
         */
        private List<String> getFieldNames(Field f) {
            SerializedName annotation = f.getAnnotation(SerializedName.class);
            if (annotation == null) {
                String name = namingStrategy.translateName(f);
                return Collections.singletonList(name);
            }
            String serializedName = annotation.value();
            String[] alternates = annotation.alternate();
            if (alternates.length == 0) {
                return Collections.singletonList(serializedName);
            }
            List<String> fieldNames = new ArrayList<>(alternates.length + 1);
            fieldNames.add(serializedName);
            fieldNames.addAll(Arrays.asList(alternates));
            return fieldNames;
        }

        @Override
        public <T> TypeAdapter<T> create(Gson gson, final TypeToken<T> type) {
            Class<? super T> raw = type.getRawType();
            if (!Object.class.isAssignableFrom(raw) || raw.isPrimitive()) {
                return null; // it's a primitive!
            }
            ObjectConstructor<T> objectConstructor = this.constructor.get(type);
            return new ObjectAccessAdapter<>(objectConstructor, getBoundFields(gson, type, raw));
        }

        private FieldAccessAdapter createBoundField(Gson context, Transformer transformer, Field field, String name, TypeToken<?> fieldType, boolean serialize, boolean deserialize) {
            boolean isPrimitive = Primitives.isPrimitive(fieldType.getRawType());
            // special casing primitives here saves ~5% on Android...
            JsonAdapter annotation = field.getAnnotation(JsonAdapter.class);
            TypeAdapter<?> mapped = null;
            if (annotation != null) {
                mapped = getTypeAdapter(constructor, context, fieldType, annotation);
            }
            boolean jsonAdapterPresent = mapped != null;
            if (mapped == null) mapped = context.getAdapter(fieldType);
            TypeAdapter<?> typeAdapter = mapped;
            return new FieldAccessAdapter(transformer, name, serialize, deserialize) {
                @Override
                void write(JsonWriter writer, Object value) throws IOException, IllegalAccessException {
                    if (null != transformer) {
                        transformer.form(new JsonTransformWriter(writer), field, field.get(value));
                        return;
                    }
                    Object fieldValue = field.get(value);
                    TypeAdapter adapter = jsonAdapterPresent ? typeAdapter : new TypeAdapterRuntimeTypeWrapper(context, typeAdapter, fieldType.getType());
                    adapter.write(writer, fieldValue);
                }

                @Override
                void read(JsonReader reader, Object value) throws IOException, IllegalAccessException {
                    if (null != transformer) {
                        Tool.setField(value, field, transformer.from(new JsonTransformReader(reader), field));
                        return;
                    }
                    Object fieldValue = typeAdapter.read(reader);
                    if (fieldValue != null || !isPrimitive) {
                        Tool.setField(value, field, fieldValue);
                    }
                }

                @Override
                public boolean writeField(Object value) throws IOException, IllegalAccessException {
                    if (!serialized) {
                        return false;
                    }
                    Object fieldValue = field.get(value);
                    return fieldValue != value; // avoid recursion for example for Throwable.cause
                }
            };
        }

        public TypeAdapter getTypeAdapter(ConstructorConstructor constructor, Gson gson, TypeToken<?> type, JsonAdapter annotation) {
            Object instance = constructor.get(TypeToken.get(annotation.value())).construct();
            TypeAdapter<?> typeAdapter;
            if (instance instanceof TypeAdapter) {
                typeAdapter = (TypeAdapter<?>) instance;
            } else if (instance instanceof TypeAdapterFactory) {
                typeAdapter = ((TypeAdapterFactory) instance).create(gson, type);
            } else if (instance instanceof JsonSerializer || instance instanceof JsonDeserializer) {
                JsonSerializer<?> serializer = instance instanceof JsonSerializer ? (JsonSerializer) instance : null;
                JsonDeserializer<?> deserializer = instance instanceof JsonDeserializer ? (JsonDeserializer) instance : null;
                typeAdapter = new TreeTypeAdapter(serializer, deserializer, gson, type, null);
            } else {
                throw new CompatibleException("Invalid attempt to bind an instance of " + instance.getClass().getName() + " as a @JsonAdapter for " + type.toString() + ". @JsonAdapter value must be a TypeAdapter, TypeAdapterFactory," + " JsonSerializer or JsonDeserializer.");
            }
            if (typeAdapter != null && annotation.nullSafe()) {
                typeAdapter = typeAdapter.nullSafe();
            }

            return typeAdapter;
        }

        private Map<String, FieldAccessAdapter> getBoundFields(Gson context, TypeToken<?> type, Class<?> raw) {
            Map<String, FieldAccessAdapter> result = new LinkedHashMap<>();
            if (raw.isInterface()) {
                return result;
            }
            Type declaredType = type.getType();
            while (!FS.contains(raw)) {
                Field[] fields = raw.getDeclaredFields();
                for (Field field : fields) {
                    boolean serialize = excludeField(field, true);
                    boolean deserialize = excludeField(field, false);
                    if (!serialize && !deserialize) {
                        continue;
                    }
                    Transformer transformer = ServiceLoader.load(Transformer.class).list().stream().filter(x -> x.matches(field)).findFirst().orElse(null);
                    Tool.makeAccessible(field);
                    Type fieldType = $Gson$Types.resolve(type.getType(), raw, Tool.getGenericType(field));
                    List<String> fieldNames = getFieldNames(field);
                    FieldAccessAdapter previous = null;
                    for (int i = 0, size = fieldNames.size(); i < size; ++i) {
                        String name = fieldNames.get(i);
                        if (i != 0) serialize = false; // only serialize the default name
                        FieldAccessAdapter boundField = createBoundField(context, transformer, field, name, TypeToken.get(fieldType), serialize, deserialize);
                        FieldAccessAdapter replaced = result.put(name, boundField);
                        if (previous == null) previous = replaced;
                    }
                    if (previous != null) {
                        throw new CompatibleException(declaredType + " declares multiple JSON fields named " + previous.name);
                    }
                }
                type = TypeToken.get($Gson$Types.resolve(type.getType(), raw, raw.getGenericSuperclass()));
                raw = type.getRawType();
            }
            return result;
        }
    }

    abstract static class FieldAccessAdapter {
        final String name;
        final boolean serialized;
        final boolean deserialized;
        final Transformer transformer;

        protected FieldAccessAdapter(Transformer transformer, String name, boolean serialized, boolean deserialized) {
            this.name = name;
            this.serialized = serialized;
            this.deserialized = deserialized;
            this.transformer = transformer;
        }

        abstract boolean writeField(Object value) throws IOException, IllegalAccessException;

        abstract void write(JsonWriter writer, Object value) throws IOException, IllegalAccessException;

        abstract void read(JsonReader reader, Object value) throws IOException, IllegalAccessException;
    }

    static final class ObjectAccessAdapter<T> extends TypeAdapter<T> {
        private final ObjectConstructor<T> constructor;
        private final Map<String, FieldAccessAdapter> boundFields;

        ObjectAccessAdapter(ObjectConstructor<T> constructor, Map<String, FieldAccessAdapter> boundFields) {
            this.constructor = constructor;
            this.boundFields = boundFields;
        }

        @Override
        public T read(JsonReader in) throws IOException {
            if (in.peek() == JsonToken.NULL) {
                in.nextNull();
                return null;
            }
            T instance = constructor.construct();
            try {
                in.beginObject();
                while (in.hasNext()) {
                    String name = in.nextName();
                    FieldAccessAdapter field = boundFields.get(name);
                    if (field == null || !field.deserialized) {
                        in.skipValue();
                    } else {
                        field.read(in, instance);
                    }
                }
            } catch (IllegalStateException e) {
                throw new CompatibleException(e);
            } catch (IllegalAccessException e) {
                throw new AssertionError(e);
            }
            in.endObject();
            return instance;
        }

        @Override
        public void write(JsonWriter out, T value) throws IOException {
            if (value == null) {
                out.nullValue();
                return;
            }
            out.beginObject();
            try {
                for (FieldAccessAdapter boundField : boundFields.values()) {
                    if (boundField.writeField(value)) {
                        out.name(boundField.name);
                        boundField.write(out, value);
                    }
                }
            } catch (IllegalAccessException e) {
                throw new AssertionError(e);
            }
            out.endObject();
        }
    }

    static final class TypeAdapterRuntimeTypeWrapper<T> extends TypeAdapter<T> {
        private final Gson context;
        private final TypeAdapter<T> delegate;
        private final Type type;

        TypeAdapterRuntimeTypeWrapper(Gson context, TypeAdapter<T> delegate, Type type) {
            this.context = context;
            this.delegate = delegate;
            this.type = type;
        }

        @Override
        public T read(JsonReader in) throws IOException {
            return delegate.read(in);
        }

        @Override
        public void write(JsonWriter out, T value) throws IOException {
            // Order of preference for choosing type adapters
            // First preference: a type adapter registered for the runtime type
            // Second preference: a type adapter registered for the declared type
            // Third preference: reflective type adapter for the runtime type (if it is a sub class of the declared type)
            // Fourth preference: reflective type adapter for the declared type
            getTypeAdapter(value).write(out, value);
        }

        @SuppressWarnings({"unchecked"})
        private TypeAdapter<T> getTypeAdapter(T value) {
            Type runtimeType = getRuntimeTypeIfMoreSpecific(type, value);
            if (runtimeType == type) {
                return delegate;
            }
            TypeAdapter<?> runtimeTypeAdapter = context.getAdapter(TypeToken.get(runtimeType));
            if (!(runtimeTypeAdapter instanceof ReflectiveTypeAdapterFactory.Adapter)) {
                // The user registered a type adapter for the runtime type, so we will use that
                return (TypeAdapter<T>) runtimeTypeAdapter;
            } else if (!(delegate instanceof ReflectiveTypeAdapterFactory.Adapter)) {
                // The user registered a type adapter for Base class, so we prefer it over the
                // reflective type adapter for the runtime type
                return delegate;
            } else {
                // Use the type adapter for runtime type
                return (TypeAdapter<T>) runtimeTypeAdapter;
            }
        }

        /**
         * Finds a compatible runtime type if it is more specific
         */
        private Type getRuntimeTypeIfMoreSpecific(Type type, Object value) {
            if (value != null && (type instanceof TypeVariable<?> || type instanceof Class<?>)) {
                type = value.getClass();
            }
            return type;
        }
    }

    static final class FieldNamingPolicy implements FieldNamingStrategy {

        @Override
        public String translateName(Field field) {
            return GsonCodec.translateName(field);
        }
    }

    @AllArgsConstructor
    static final class JsonTransformReader implements Transformer.Reader {

        private final JsonReader reader;

        @Override
        public Number readNumber() throws IOException {
            if (reader.peek() == JsonToken.NULL) {
                reader.nextNull();
                return null;
            }
            return reader.nextLong();
        }

        @Override
        public String readString() throws IOException {
            if (reader.peek() == JsonToken.NULL) {
                reader.nextNull();
                return null;
            }
            return reader.nextString();
        }

        @Override
        public Boolean readBoolean() throws IOException {
            if (reader.peek() == JsonToken.NULL) {
                reader.nextNull();
                return false;
            }
            return reader.nextBoolean();
        }
    }

    @AllArgsConstructor
    static final class JsonTransformWriter implements Transformer.Writer {

        private final JsonWriter writer;

        @Override
        public void write(Number value) throws IOException {
            writer.value(value);
        }

        @Override
        public void write(String value) throws IOException {
            writer.value(value);
        }

        @Override
        public void write(Boolean value) throws IOException {
            writer.value(value);
        }

        @Override
        public void writeNull() throws IOException {
            writer.nullValue();
        }
    }

    static String translateName(Field field) {
        String name = Tool.getAnnotationMeta(field, Index.class, Index::name);
        if (Tool.required(name)) {
            return name;
        }
        name = Tool.getAnnotationMeta(field, Format.class, Format::value);
        if (Tool.required(name)) {
            return name;
        }
        return field.getName();
    }


}
