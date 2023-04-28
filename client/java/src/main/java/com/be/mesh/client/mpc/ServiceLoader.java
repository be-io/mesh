/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.prsim.Routable;
import com.be.mesh.client.tool.Features;
import com.be.mesh.client.tool.Tool;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import java.io.IOException;
import java.io.InputStream;
import java.lang.reflect.Field;
import java.net.URL;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

/**
 * The type Extension loader.
 * This is done by loading the properties file.
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
public final class ServiceLoader<T> {

    private static final Map<Class<?>, ServiceLoader<?>> LOADERS = new ConcurrentHashMap<>();

    private static final Map<String, String> RESOURCES = new ConcurrentHashMap<>();

    private final Map<String, Map<String, Instance<T>>> providers = new ConcurrentHashMap<>(1);

    private final Class<T> spi;

    private final String first;

    /**
     * Instantiates a new Extension loader.
     *
     * @param spi the spi.
     */
    private ServiceLoader(Class<T> spi) {
        this.spi = spi;
        this.first = Optional.ofNullable(spi.getAnnotation(SPI.class)).map(SPI::value).orElse("");
    }

    /**
     * Gets extension loader.
     *
     * @param <T> the type parameter
     * @param spi the spi
     * @return the extension loader.
     */
    @SuppressWarnings("unchecked")
    public static <T> ServiceLoader<T> load(Class<T> spi) {
        return (ServiceLoader<T>) Optional.ofNullable(spi).map(c -> {
            if (!spi.isInterface()) {
                throw new MeshException("SPI %s is not interface!", spi.getName());
            }
            if (!spi.isAnnotationPresent(SPI.class)) {
                log.warn("SPI {} without @{} Annotation", spi.getName(), SPI.class.getSimpleName());
            }
            return LOADERS.computeIfAbsent(spi, key -> new ServiceLoader<>(spi));
        }).orElseThrow(() -> new MeshException("SPI interface required"));
    }

    /**
     * Get the resource from the plugin META-INF/janus/xxx.yaml or remote.
     *
     * @param name plugin name
     * @return profile object
     */
    public static ByteBuffer resource(String name) {
        String resource = RESOURCES.computeIfAbsent(name, key -> {
            String suffix = key.contains(".") ? "" : ".yaml";
            String mf = String.format("%s%s%s%s%s%s", "META-INF", "/", "mesh", "/", key, suffix);
            for (ClassLoader loader : getClassLoader()) {
                List<URL> urls = scanPath(mf, loader);
                if (Tool.optional(urls)) {
                    return "";
                }
                if (urls.size() > 1) {
                    log.warn("Resource {} has duplicate in {}", mf, urls.stream().map(URL::toString).collect(Collectors.joining(" or ")));
                }
                try (InputStream input = urls.get(0).openStream()) {
                    return Tool.read(input);
                } catch (IOException e) {
                    throw new MeshException(e, "Cant open resource %s at %s", mf, urls.get(0).toString());
                }
            }
            return "";
        });
        return ByteBuffer.wrap(resource.getBytes(StandardCharsets.UTF_8));
    }

    /**
     * Gets default join.
     *
     * @return the default join.
     */
    public T getDefault() {
        return get(first);
    }

    /**
     * Gets the default name.
     *
     * @return default spi provider name
     */
    public String defaultName() {
        return this.first;
    }

    /**
     * Gets join.
     *
     * @param name the name
     * @return the join.
     */
    public T get(String name) {
        return this.getOptional(name).orElseThrow(() -> new MeshException(String.format("SPI named %s not exist.", name)));
    }

    /**
     * Get the spi provider with special name.
     *
     * @param name spi provider name
     * @return spi provider
     */
    public Optional<T> getOptional(String name) {
        if (Tool.optional(name)) {
            throw new MeshException("Get join name is required");
        }
        return Optional.ofNullable(getInstances().get(name)).flatMap(Instance::get);
    }

    /**
     * Gets join.
     *
     * @param name the name
     * @return the join.
     */
    public T getIfAbsent(String name, String defaults) {
        return getOptional(name).orElseGet(() -> this.get(defaults));
    }

    /**
     * List all spi services.
     *
     * @return All spi services.
     */
    public List<T> list() {
        List<T> vs = new ArrayList<>();
        getInstances().forEach((name, instance) -> instance.getStrict().ifPresent(vs::add));
        return vs;
    }

    /**
     * Map all spi services.
     *
     * @return All spi services.
     */
    public Map<String, T> map() {
        Map<String, T> vs = new HashMap<>();
        getInstances().forEach((name, instance) -> instance.getStrict().ifPresent(x -> vs.put(name, x)));
        return vs;
    }

    /**
     * Get all instances.
     *
     * @return spi instances
     */
    private Map<String, Instance<T>> getInstances() {
        return providers.computeIfAbsent("$", key -> {
            List<String> excludes = Features.getInactive(this.spi);
            List<String> includes = Features.getActive(this.spi);
            Map<String, List<Class<T>>> types = this.loadTypes();
            Map<String, Instance<T>> instances = new HashMap<>(types.size());
            types.forEach((name, type) -> instances.put(name, new Instance<>(name, type, includes.contains(name), excludes.contains(name))));
            return instances;
        });
    }

    /**
     * Load files under META-INF/services.
     */
    private Map<String, List<Class<T>>> loadTypes() {
        String manifest = String.format("%s%s%s%s%s", "META-INF", "/", "services", "/", this.spi.getName());
        Map<String, List<Class<T>>> all = new HashMap<>();
        for (ClassLoader loader : getClassLoader()) {
            for (URL url : scanPath(manifest, loader)) {
                Map<String, List<Class<T>>> service = this.loadTypes(url, loader);
                service.forEach((key, types) -> all.computeIfAbsent(key, k -> new ArrayList<>()).addAll(types));
            }
        }
        return all.entrySet().stream().filter(x -> !x.getValue().isEmpty()).filter(entry -> {
            if (entry.getValue().size() > 1) {
                log.warn("SPI provider name duplicate with {}, consider disable provider {}", entry.getKey(), entry.getValue().stream().map(Class::getName).collect(Collectors.joining(" or ")));
            }
            return true;
        }).collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    }

    /**
     * Scan all the path in jar.
     *
     * @param path file path
     * @return all file url
     */
    private static List<URL> scanPath(String path, ClassLoader loader) {
        Enumeration<URL> urls = Optional.ofNullable(loader).map(x -> {
            try {
                return x.getResources(path);
            } catch (IOException e) {
                throw new MeshException(e, "Cant open resource %s", path);
            }
        }).orElseGet(() -> {
            try {
                return ClassLoader.getSystemResources(path);
            } catch (IOException e) {
                throw new MeshException(e, "Cant open system resource %s", path);
            }
        });
        if (null == urls) {
            return new ArrayList<>(0);
        }
        List<URL> all = new ArrayList<>();
        while (urls.hasMoreElements()) {
            all.add(urls.nextElement());
        }
        return all;
    }

    /**
     * Load the service type with the url.
     *
     * @param url resource url
     * @return loaded class types
     */
    private Map<String, List<Class<T>>> loadTypes(URL url, ClassLoader loader) {
        try (InputStream input = url.openStream()) {
            String services = Tool.read(input);
            return this.resolveTypes(services, loader);
        } catch (IOException e) {
            throw new MeshException(e, "Cant open resource %s", url.toString());
        }
    }

    /**
     * Resolve the services types.
     *
     * @param types provider types
     * @return loaded types
     */
    @SuppressWarnings({"unchecked", "rawTypes"})
    private Map<String, List<Class<T>>> resolveTypes(String types, ClassLoader loader) {
        return Arrays.stream(types.split("\n")).flatMap(name -> Arrays.stream(name.split(" "))).filter(Tool::required).map(String::trim).distinct().map(name -> {
            try {
                return Optional.of(Class.forName(name, true, loader));
            } catch (ClassNotFoundException | NoClassDefFoundError e) {
                log.error("SPI provider {} cant be initialed, {}", name, e.getMessage());
                return Optional.<Class<?>>empty();
            }
        }).filter(Optional::isPresent).map(Optional::get).filter(type -> {
            if (this.spi.isAssignableFrom(type)) {
                return true;
            }
            log.warn("SPI resources load error, {} subtype is not of {}", type.getName(), spi.getName());
            return false;
        }).map(x -> (Class<T>) x).collect(Collectors.groupingBy(type -> {
            SPI spi = type.getAnnotation(SPI.class);
            if (null != spi && Tool.required(spi.value())) {
                return spi.value();
            }
            return Tool.formatObjectName(type).replaceFirst(this.spi.getSimpleName(), "");
        }));
    }

    /**
     * Exchange classloader able.
     */
    private static List<ClassLoader> getClassLoader() {
        ClassLoader ml = ServiceLoader.class.getClassLoader();
        ClassLoader tl = Thread.currentThread().getContextClassLoader();
        if (ml == tl) {
            return Collections.singletonList(ml);
        }
        return Arrays.asList(ml, tl);
    }

    @Getter
    private static class Instance<T> {
        private final String name;
        private final List<Anyone<T>> providers;
        private final boolean include;
        private final boolean exclude;


        private Instance(String name, List<Class<T>> type, boolean include, boolean exclude) {
            this.name = name;
            this.providers = type.stream().map(x -> new Anyone<>(name, x)).collect(Collectors.toList());
            this.include = include;
            this.exclude = exclude;
        }

        private Optional<T> first(boolean strict) {
            return providers.stream().filter(x -> {
                if (!strict) {
                    return true;
                }
                if (this.include) {
                    return true;
                }
                if (this.exclude) {
                    return false;
                }
                return x.isEnable();
            }).min(Comparator.comparing(Anyone::getPriority)).map(Anyone::getIfAbsent);
        }

        public Optional<T> get() {
            Optional<T> v = getStrict();
            if (v.isPresent()) {
                return v;
            }
            return first(false);
        }

        public Optional<T> getStrict() {
            return first(true);
        }

    }

    @Getter
    private static class Anyone<T> {
        private static final Map<String, Object> value = new ConcurrentHashMap<>();
        private final String name;
        private final Class<T> type;
        private final boolean prototype;
        private final int priority;
        private final boolean enable;

        public Anyone(String name, Class<T> type) {
            this.name = name;
            this.type = type;
            this.prototype = Optional.ofNullable(type.getAnnotation(SPI.class)).map(SPI::prototype).orElse(false);
            this.priority = Optional.ofNullable(type.getAnnotation(SPI.class)).map(SPI::priority).orElse(0);
            this.enable = Optional.ofNullable(type.getAnnotation(SPI.class)).map(SPI::enable).orElse(true);
        }

        @SuppressWarnings("unchecked")
        public T getIfAbsent() {
            if (prototype) {
                return create();
            }
            return (T) value.computeIfAbsent(type.getName(), key -> this.create());
        }

        private T create() {
            try {
                T instance = type.getDeclaredConstructor().newInstance();
                for (Field variable : type.getDeclaredFields()) {
                    if (!variable.isAnnotationPresent(MPI.class)) {
                        continue;
                    }
                    Tool.makeAccessible(variable);
                    Class<?> referenceType = Tool.detectReferenceType(variable);
                    Object reference = ServiceProxy.proxy(referenceType);
                    if (Tool.isRoutable(variable)) {
                        Tool.setField(instance, variable, Routable.of(reference));
                    } else {
                        Tool.setField(instance, variable, reference);
                    }
                }
                return instance;
            } catch (Exception e) {
                throw new MeshException(e, "SPI instance(name: %s class: %s)  could not be instantiated: %s", name, type.getName(), e.getMessage());
            }
        }
    }

}
