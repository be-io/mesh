/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import com.be.mesh.client.cause.CompatibleException;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.cause.ValidationException;
import com.be.mesh.client.mpc.Factory;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.Routable;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;

import java.io.*;
import java.lang.annotation.Annotation;
import java.lang.reflect.Proxy;
import java.lang.reflect.*;
import java.net.*;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.security.AccessController;
import java.security.PrivilegedAction;
import java.util.*;
import java.util.concurrent.ThreadLocalRandom;
import java.util.function.Function;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collector;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import java.util.zip.GZIPInputStream;
import java.util.zip.GZIPOutputStream;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public final class Tool {

    private Tool() {
    }

    public static final String LOCAL_NODE_ID = "LX0000000000000";
    public static final String LOCAL_INST_ID = "JG0000000000000000";
    public static final int DEFAULT_MESH_PORT = 570;
    // returned port range is [30000, 39999]
    private static final int RND_PORT_START = 30000;
    private static final int RND_PORT_RANGE = 10000;
    // valid port range is (0, 65535]
    private static final int MIN_PORT = 1;
    private static final int MAX_PORT = 65535;
    private static final BitSet USED_PORT = new BitSet(65536);
    //
    private static final Pattern PATTERN = Pattern.compile("([^&=]+)(=?)([^&]+)?");
    private static final Class<?>[] BASIC = new Class[]{String.class, Boolean.class, Character.class, Byte.class, Short.class, Integer.class, Long.class, Float.class, Double.class, Void.class};

    public static boolean isBasicType(Class<?> kind) {
        if (kind.isPrimitive()) {
            return true;
        }
        for (Class<?> type : BASIC) {
            if (type.isAssignableFrom(kind)) {
                return true;
            }
        }
        return false;
    }

    public static String boxType(String type) {
        switch (type) {
            case "boolean":
                return Boolean.class.getCanonicalName();
            case "int":
                return Integer.class.getCanonicalName();
            case "short":
                return Short.class.getCanonicalName();
            case "long":
                return Long.class.getCanonicalName();
            case "double":
                return Double.class.getCanonicalName();
            case "float":
                return Float.class.getCanonicalName();
            case "char":
                return Character.class.getCanonicalName();
            case "byte":
                return Byte.class.getCanonicalName();
            default:
                return type;
        }
    }

    public static boolean isNumeric(CharSequence chars) {
        if (!required(chars)) {
            return false;
        }
        for (int index = 0; index < chars.length(); index++) {
            if (!Character.isDigit(chars.charAt(index))) {
                return false;
            }
        }
        return true;
    }

    public static String read(InputStream input) throws IOException {
        return new String(readBytes(input), StandardCharsets.UTF_8);
    }

    public static byte[] readBytes(InputStream input) throws IOException {
        ByteArrayOutputStream output = new ByteArrayOutputStream();
        byte[] buffer = new byte[2048];
        for (int size = input.read(buffer); size > 0; size = input.read(buffer)) {
            output.write(buffer, 0, size);
        }
        return output.toByteArray();
    }

    // true if anyone is optional
    @SafeVarargs
    public static <T> boolean optional(T... inputs) {
        return !required(inputs);
    }

    // true if everyone is required
    @SafeVarargs
    public static <T> boolean required(T... inputs) {
        if (null == inputs || inputs.length < 1) {
            return false;
        }
        for (T input : inputs) {
            if (null == input) {
                return false;
            }
            if (input instanceof CharSequence) {
                return !((CharSequence) input).chars().allMatch(Character::isWhitespace);
            }
            if (input instanceof Collection) {
                return !((Collection<?>) input).isEmpty();
            }
            if (input instanceof Map) {
                return !((Map<?, ?>) input).isEmpty();
            }
            if (input.getClass().isArray()) {
                return Array.getLength(input) > 0;
            }
        }
        return true;
    }

    // return anyone not optional
    @SafeVarargs
    public static <T> T anyone(T... inputs) {
        for (T input : inputs) {
            if (required(input)) {
                return input;
            }
        }
        return inputs[0];
    }

    @SafeVarargs
    public static <T> T requiredOne(T... inputs) {
        for (T input : inputs) {
            if (required(input)) {
                return input;
            }
        }
        return inputs[0];
    }

    public static int indexOf(String v, String i, int start) {
        return null == v ? -1 : v.indexOf(i, start);
    }

    public static String[] split(String v, String e) {
        return null == v ? new String[0] : v.split(e);
    }

    public static String substring(String v, int begin, int length) {
        if (length < 0) {
            return null == v || begin >= v.length() + length ? "" : v.substring(begin, v.length() + length);
        }
        return null == v || v.length() < begin + length ? "" : v.substring(begin, begin + length);
    }

    public static boolean equals(Object v, Object c) {
        if (null == v && null == c) {
            return true;
        }
        if (null == v || null == c) {
            return false;
        }
        return v.equals(c);
    }

    public static boolean matches(String regex, String chars) {
        return Pattern.compile(regex).matcher(chars).matches();
    }

    public static String repeat(char ch, int count) {
        StringBuilder buffer = new StringBuilder();
        for (int i = 0; i < count; ++i) {
            buffer.append(ch);
        }
        return buffer.toString();
    }

    @SneakyThrows
    public static Map<String, String> parseQuery(String query) {
        Map<String, String> params = new LinkedHashMap<>();
        if (optional(query)) {
            return params;
        }
        Matcher matcher = PATTERN.matcher(query);
        while (matcher.find()) {
            String name = URLDecoder.decode(matcher.group(1), StandardCharsets.UTF_8.name());
            String equals = matcher.group(2);
            String value = matcher.group(3);
            String def = required(equals) ? "" : null;
            params.put(name, optional(value) ? def : URLDecoder.decode(value, StandardCharsets.UTF_8.name()));
        }
        return params;
    }


    public static String toCamel(String name) {
        return Optional.ofNullable(name).map(names -> Arrays.stream(names.split("\\."))).orElseGet(Stream::empty).flatMap(names -> Arrays.stream(names.split("_"))).filter(Objects::nonNull).map(names -> Character.toLowerCase(names.charAt(0)) + names.length() > 1 ? names.substring(1) : "").collect(Collectors.joining());
    }

    public static String toSnake(String str) {
        return str.replaceAll("[A-Z]", "_$0").toLowerCase();
    }

    public static <T> T newInstance(Class<T> clazz) {
        return uncheck(() -> clazz.getDeclaredConstructor().newInstance());
    }

    public static <T> T uncheck(LambdaSp<T> done) {
        try {
            return done.execute();
        } catch (UndeclaredThrowableException | InvocationTargetException e) {
            Throwable x = destructor(e.getCause());
            if (x instanceof Error) {
                throw (Error) x;
            }
            throw (RuntimeException) x;
        } catch (Error | RuntimeException e) {
            throw e;
        } catch (Throwable e) {
            throw new MeshException(e);
        }
    }

    public static void uncheck(LambdaRn done) {
        try {
            done.execute();
        } catch (UndeclaredThrowableException | InvocationTargetException e) {
            Throwable x = destructor(e.getCause());
            if (x instanceof Error) {
                throw (Error) x;
            }
            throw (RuntimeException) x;
        } catch (Error | RuntimeException e) {
            throw e;
        } catch (Throwable e) {
            throw new MeshException(e);
        }
    }

    public static boolean notSystemType(Class<?> type) {
        return null != type && type != Object.class && !isBasicType(type) && !isProxy(type) && !type.isInterface() && !type.isEnum() && !type.isAnnotation() && null != type.getPackage() && !type.getPackage().getName().startsWith("java");
    }

    public static <T, A extends Annotation> T getAnnotationMeta(AnnotatedElement element, Class<A> macro, Function<A, T> fn) {
        return Optional.ofNullable(element).map(x -> x.getAnnotation(macro)).map(fn).filter(Tool::required).orElse(null);
    }

    public static Method getMethod(Object target, String method, Class<?>... parameters) {
        boolean isClass = target instanceof Class<?>;
        Class<?> type = (isClass ? (Class<?>) target : target.getClass());
        try {
            Method invoke = type.getDeclaredMethod(method, parameters);
            makeAccessible(invoke);
            return invoke;
        } catch (NoSuchMethodException e) {
            try {
                Method invoke = type.getMethod(method, parameters);
                makeAccessible(invoke);
                return invoke;
            } catch (NoSuchMethodException ex) {
                throw new CompatibleException(e);
            }
        }
    }

    public static <T extends Annotation> Optional<T> getAnnotation(Class<T> annotation, AnnotatedElement... elements) {
        return Arrays.stream(elements).filter(Objects::nonNull).filter(e -> e.isAnnotationPresent(annotation)).map(e -> {
            if (e instanceof Class) {
                Deque<Class<?>> queue = new ArrayDeque<>();
                queue.push((Class<?>) e);
                while (null != queue.peek()) {
                    Class<?> type = queue.pop();
                    if (type.isAnnotationPresent(annotation)) {
                        return type.getAnnotation(annotation);
                    }
                    if (null != type.getSuperclass()) {
                        queue.push(type.getSuperclass());
                    }
                    for (Class<?> inter : type.getInterfaces()) {
                        queue.push(inter);
                    }
                }
            }
            return e.getAnnotation(annotation);
        }).findFirst();
    }

    public static boolean isProxy(Object target) {
        return null != target && (Proxy.isProxyClass(target.getClass()) || target.getClass().getName().contains("$$"));
    }

    public static Object exposeSpringProxy(Object proxy) {
        return uncheck(() -> {
            Deque<Object> objects = new ArrayDeque<>();
            objects.push(proxy);
            while (!objects.isEmpty()) {
                Object object = objects.pop();
                if (Proxy.isProxyClass(object.getClass())) {
                    exposeSpringProxy(object, "h").ifPresent(objects::push);
                    continue;
                }
                if (object.getClass().getName().contains("$$")) {
                    exposeSpringProxy(object, "CGLIB$CALLBACK_0").ifPresent(objects::push);
                    continue;
                }
                return object;
            }
            return proxy;
        });
    }

    public static Optional<Object> exposeSpringProxy(Object proxied, String name) {
        try {
            Field h = proxied.getClass().getSuperclass().getDeclaredField(name);
            makeAccessible(h);
            Object proxy = h.get(proxied);
            if (null == proxy) {
                return Optional.empty();
            }
            Field advised = proxy.getClass().getDeclaredField("advised");
            makeAccessible(advised);
            Method getTargetSource = advised.get(proxy).getClass().getMethod("getTargetSource");
            Object targetSource = getTargetSource.invoke(advised.get(proxy));
            return Optional.ofNullable(targetSource.getClass().getMethod("getTarget").invoke(targetSource));
        } catch (NoSuchFieldException | NoSuchMethodException | IllegalAccessException | InvocationTargetException e) {
            log.warn("Unboxed spring proxy bean with error, {}", e.getMessage());
            return Optional.empty();
        }
    }

    public static String compress(String data) {
        return uncheck(() -> {
            try (ByteArrayOutputStream bos = new ByteArrayOutputStream(data.length()); GZIPOutputStream gzip = new GZIPOutputStream(bos)) {
                gzip.write(data.getBytes(Charset.defaultCharset()));
                gzip.finish();
                return Base64.getEncoder().encodeToString(bos.toByteArray());
            }
        });
    }

    public static String decompress(String compressed) {
        return uncheck(() -> {
            try (ByteArrayInputStream bis = new ByteArrayInputStream(Base64.getDecoder().decode(compressed)); GZIPInputStream gis = new GZIPInputStream(bis)) {
                return new String(readBytes(gis), Charset.defaultCharset());
            }
        });
    }

    public static <K, V> Collector<V, Map<K, V>, Map<K, V>> map(Function<? super V, ? extends K> key) {
        return Collector.of(HashMap::new, (map, element) -> map.putIfAbsent(key.apply(element), element), (left, right) -> {
            left.putAll(right);
            return left;
        });
    }

    public static <K, V> Collector<V, List<V>, Stream<V>> distinct(Function<? super V, ? extends K> key) {
        return Collector.of(ArrayList::new, (list, element) -> {
            if (list.stream().noneMatch(node -> key.apply(node).equals(key.apply(element)))) {
                list.add(element);
            }
        }, (left, right) -> {
            left.addAll(right);
            return left;
        }, Collection::stream);
    }

    public static <K, V> Collector<V, List<V>, Stream<List<V>>> combine() {
        return Collector.of(ArrayList::new, List::add, (left, right) -> {
            left.addAll(right);
            return left;
        }, Stream::of);
    }

    public static String formatObjectName(Class<?> type) {
        String name = type.getSimpleName();
        if (name.length() > 2 && Character.isUpperCase(name.charAt(0)) && Character.isUpperCase(name.charAt(1))) {
            return Character.toLowerCase(name.charAt(1)) + name.substring(2);
        }
        return Character.toLowerCase(name.charAt(0)) + name.substring(1);
    }

    public static String firstLowerCase(String name) {
        return optional(name) || name.length() < 1 ? name : String.format("%s%s", Character.toLowerCase(name.charAt(0)), name.substring(1));
    }

    public static String firstUpperCase(String name) {
        return optional(name) || name.length() < 1 ? name : String.format("%s%s", Character.toUpperCase(name.charAt(0)), name.substring(1));
    }

    public static String toUpperCase(String name) {
        return Optional.ofNullable(name).map(x -> {
            StringBuilder bu = new StringBuilder(x.length());
            for (char c : x.toCharArray()) {
                bu.append(Character.toUpperCase(c));
            }
            return bu.toString();
        }).orElse("");
    }

    public static String toLowerCase(String name) {
        return Optional.ofNullable(name).map(x -> {
            StringBuilder bu = new StringBuilder(x.length());
            for (char c : x.toCharArray()) {
                bu.append(Character.toLowerCase(c));
            }
            return bu.toString();
        }).orElse("");
    }

    public static boolean contains(String s, String v) {
        return null != s && s.contains(v);
    }

    public static boolean startWith(String s, String v) {
        return null != s && s.startsWith(v);
    }

    public static boolean endsWith(String s, String v) {
        return null != s && s.endsWith(v);
    }

    public static boolean isClassPresent(String name) {
        try {
            Class.forName(name);
            return true;
        } catch (ClassNotFoundException e) {
            return false;
        }
    }

    public static boolean canService(Method method) {
        return method.getDeclaringClass() != Object.class && Modifier.isPublic(method.getModifiers()) && !Modifier.isStatic(method.getModifiers());
    }

    @SuppressWarnings({"unchecked"})
    public static void makeAccessible(AccessibleObject object) {
        if (object.isAccessible()) {
            return;
        }
        if (null == System.getSecurityManager()) {
            object.setAccessible(true);
            return;
        }
        AccessController.doPrivileged((PrivilegedAction) () -> {
            object.setAccessible(true);
            return object;
        });
    }

    public static Field getVariable(Object owner, String name) throws NoSuchFieldException {
        try {
            return owner.getClass().getField(name);
        } catch (NoSuchFieldException e) {
            return owner.getClass().getDeclaredField(name);
        }
    }

    public static void setField(Object owner, Field field, Object value) throws IllegalAccessException {
        makeAccessible(field);
        field.set(owner, value);
    }

    public static void setFieldInt(Object owner, Field field, int value) throws IllegalAccessException {
        makeAccessible(field);
        field.setInt(owner, value);
    }

    public static Object getField(Object owner, Field field) throws IllegalAccessException {
        makeAccessible(field);
        return field.get(owner);
    }

    public static void setField(Object owner, String name, Object value) throws NoSuchFieldException, IllegalAccessException {
        setField(owner, getVariable(owner, name), value);
    }

    public static Object getField(Object owner, String name) throws NoSuchFieldException, IllegalAccessException {
        return getField(owner, getVariable(owner, name));
    }

    public static Type getGenericType(Field field) {
        try {
            return field.getGenericType();
        } catch (MalformedParameterizedTypeException e) {
            log.warn("Field {} has invalid generic type, raw is {}.", field.getName(), field.getType().getCanonicalName());
            return field.getType();
        }
    }

    public static boolean isRoutable(Field variable) {
        if (variable.getGenericType() instanceof ParameterizedType) {
            return Routable.class.isAssignableFrom(variable.getType());
        }
        return false;
    }

    public static Class<?> detectReferenceType(Field variable) {
        if (!isRoutable(variable)) {
            return variable.getType();
        }
        Type[] types = ((ParameterizedType) variable.getGenericType()).getActualTypeArguments();
        if (types.length < 1 || types[0] == Object.class) {
            return variable.getType();
        }
        return (Class<?>) types[0];
    }

    public static synchronized int getAvailablePort() {
        int rpt = Math.max(RND_PORT_START + ThreadLocalRandom.current().nextInt(RND_PORT_RANGE), MIN_PORT);
        for (int index = rpt; index < MAX_PORT; index++) {
            if (USED_PORT.get(index)) {
                continue;
            }
            try (ServerSocket ignored = new ServerSocket(index)) {
                USED_PORT.set(index);
                return index;
            } catch (IOException e) {
                // continue
            }
        }
        return rpt;
    }

    /**
     * Destructor the exception.
     *
     * @param ecs ignore exceptions.
     * @param e   Throwable
     * @return Root cause
     */
    public static Throwable destructor(Throwable e, Class<?>... ecs) {
        if (e instanceof RuntimeException || e instanceof Error) {
            return e;
        }
        if (Tool.optional(ecs)) {
            return new MeshException(e);
        }
        if (Arrays.stream(ecs).anyMatch(ec -> ec.isAssignableFrom(e.getClass()))) {
            return e;
        }
        return new MeshException(e);
    }

    public static String getStackTrace(Throwable e) {
        StringWriter sw = new StringWriter();
        PrintWriter pw = new PrintWriter(sw, true);
        e.printStackTrace(pw);
        return sw.getBuffer().toString();
    }

    public static byte[] serializeCause(Throwable e) throws IOException {
        return serialize(e);
    }

    public static byte[] serialize(Object e) throws IOException {
        try (ByteArrayOutputStream buffer = new ByteArrayOutputStream(); ObjectOutputStream os = new ObjectOutputStream(buffer);) {
            os.writeObject(e);
            os.flush();
            return buffer.toByteArray();
        }
    }

    public static Throwable deserializeCause(byte[] bytes) throws IOException, ClassNotFoundException {
        return (Throwable) deserialize(bytes);
    }

    public static Object deserialize(byte[] bytes) throws IOException, ClassNotFoundException {
        try (ByteArrayInputStream buffer = new ByteArrayInputStream(bytes); ObjectInputStream os = new ObjectInputStream(buffer);) {
            return os.readObject();
        }
    }

    public static String replaceJSONValue(String json, List<String> keys, String value) {
        for (String key : keys) {
            if (null != json && equals(json, key)) {
                String pattern = String.format("(?<=\"%s\":\")[^\"]+(?=\")", key);
                json = json.replaceAll(pattern, value);
            }
        }
        return json;
    }

    public static StackTraceElement callerStack() {
        return Thread.currentThread().getStackTrace()[3];
    }

    public static Class<?> callerClass() {
        return uncheck(() -> Class.forName(Thread.currentThread().getStackTrace()[5].getClassName()));
    }

    public static Class<?> callerClass(int deep) {
        return uncheck(() -> Class.forName(Thread.currentThread().getStackTrace()[5 + deep].getClassName()));
    }

    public static String getProperty(String def, String... keys) {
        return getProperty(def, false, keys);
    }

    public static String getProperty(String def, boolean onlyEnv, String... keys) {
        return uncheck(() -> {
            for (String key : keys) {
                for (String v : getPropertyValues(onlyEnv, key)) {
                    if (required(v)) {
                        return v;
                    }
                }
            }
            return def;
        });
    }

    public static List<String> getPropertyValues(boolean onlyEnv, String key) {
        List<String> values = new ArrayList<>(3);
        values.add(System.getenv(key));
        values.add(System.getProperty(key));
        if (onlyEnv) {
            return values;
        }
        String cv = ServiceLoader.load(Factory.class).list().stream().map(f -> {
            try (InputStream stream = f.getResource(key)) {
                return read(stream);
            } catch (IOException e) {
                log.error(String.format("Read profile from %s", f.getClass().getCanonicalName()), e);
                return "";
            }
        }).filter(Tool::required).findFirst().orElse("");
        values.add(cv);
        return values;
    }

    public static int getProperty(int def, String... keys) {
        String prop = getProperty(String.valueOf(def), keys);
        return isNumeric(prop) ? Integer.parseInt(prop) : def;
    }

    public static final Once<String> IP = Once.with(() -> {
        String os = System.getProperty("os.name");
        if (os.startsWith("Mac OS") || os.startsWith("Windows")) {
            try {
                return InetAddress.getLocalHost().getHostAddress();
            } catch (Exception e) {
                log.error("Resolve localhost ip failed.", e);
                return "";
            }
        }
        try (DatagramSocket socket = new DatagramSocket()) {
            socket.connect(InetAddress.getByName("8.8.8.8"), 100);
            return socket.getLocalAddress().getHostAddress();
        } catch (Exception e) {
            log.error("Resolve localhost ip failed.", e);
            return "";
        }
    });

    public static final Once<String> HOST_NAME = Once.with(() -> {
        try {
            return InetAddress.getLocalHost().getHostName();
        } catch (Exception ignore) {
            return "";
        }
    });

    public static final Once<URI> MESH_RUNTIME = Once.with(() -> {
        String uri = getProperty(String.format("%s:%d", IP.get(), getAvailablePort()), "mesh.runtime", "mesh.runtime", "MESH_RUNTIME");
        if (uri.contains("://")) {
            return URI.create(uri);
        }
        String[] parts = uri.split(":");
        if (parts.length > 1) {
            return URI.create(String.format("https://%s:%s", parts[0], parts[1]));
        }
        return URI.create(String.format("https://%s:80", parts[0]));
    });

    private static String ipHex() {
        return IP_HEX.get();
    }

    public static String newTraceId() {
        return ipHex() + UUID.getInstance().shortUUID();
    }

    public static String newSpanId(String spanId, int index) {
        if (optional(spanId)) {
            return "0";
        }
        if (spanId.length() > 255) {
            return "0";
        }
        return String.format("%s.%d", spanId, index);
    }

    public static void must(boolean expression, String message) {
        if (!expression) {
            throw new ValidationException(message);
        }
    }

    public static boolean isInMyNet(String nodeId, String... ids) {
        if (optional(nodeId) || optional(ids)) {
            return false;
        }
        String seq = nodeId.length() == 18 ? Tool.substring(nodeId, 4, 6) : nodeId;
        seq = nodeId.length() == 15 ? Tool.substring(nodeId, 8, 6) : seq;
        for (String id : ids) {
            if (optional(id)) {
                continue;
            }
            if (LOCAL_NODE_ID.equalsIgnoreCase(id) || LOCAL_INST_ID.equalsIgnoreCase(id) || id.equalsIgnoreCase(nodeId)) {
                return true;
            }
            if (id.length() == 18 && Tool.equals(seq, Tool.substring(id, 4, 6))) {
                return true;
            }
            if (id.length() == 15 && Tool.equals(seq, Tool.substring(id, 8, 6))) {
                return true;
            }
        }
        return false;
    }

    public static String padding(String v, int length, char x) {
        if (optional(v)) {
            return Tool.repeat(x, length);
        }
        if (v.length() < length) {
            return Tool.repeat(x, length - v.length()) + v;
        }
        return v.substring(0, length);
    }

    public static final Once<String> MESH_NAME = Envs.MESH_NAME;
    public static final Once<Addrs> MESH_ADDRESS = Envs.MESH_ADDRESS;
    public static final Once<Mode> MESH_MODE = Envs.MESH_MODE;
    public static final Once<String> MESH_DIRECT = Envs.MESH_DIRECT;
    public static final Once<String> MESH_SUBSET = Envs.MESH_SUBSET;
    public static final Once<Features> MESH_FEATURE = Envs.MESH_FEATURE;
    public static final Once<String> MESH_RUNTIME_IP = Envs.MESH_RUNTIME_IP;
    public static final Once<String> MESH_RUNTIME_PORT = Envs.MESH_RUNTIME_PORT;
    public static final Once<String> MESH_ZONE = Envs.MESH_ZONE;
    public static final Once<String> MESH_CLUSTER = Envs.MESH_CLUSTER;
    public static final Once<String> MESH_CELL = Envs.MESH_CELL;
    public static final Once<String> MESH_GROUP = Envs.MESH_GROUP;
    public static final Once<String> IP_HEX = Envs.IP_HEX;
    public static final Once<Boolean> MESH_ENABLE = Envs.MESH_ENABLE;

}
