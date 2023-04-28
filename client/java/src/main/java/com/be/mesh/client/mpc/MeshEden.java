/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.*;
import com.be.mesh.client.prsim.*;
import com.be.mesh.client.struct.*;
import com.be.mesh.client.tool.LDCable;
import com.be.mesh.client.tool.Mode;
import com.be.mesh.client.tool.Tool;
import lombok.Getter;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;

import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.lang.reflect.Parameter;
import java.lang.reflect.Type;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * Like: create.tenant.omega.json.http2.lx000001.mpi.trustbe.cn
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("mesh")
@Listener(topic = MeshEden.TOPIC, code = MeshEden.CODE)
public class MeshEden implements Eden, RuntimeHook, Subscriber {

    private String aliveId = "";
    public static final String TOPIC = "mesh.registry.event";
    public static final String CODE = "refresh";
    public static final Duration DEFAULT_TIMEOUT = Duration.ofSeconds(10);
    private final Map<Class<?>, Provider> providers = new ConcurrentHashMap<>();
    private final Map<Class<?>, Map<MPI, Consumer>> consumers = new ConcurrentHashMap<>();
    private final Map<Class<?>, Map<MPI, Map<String, Instance<Reference>>>> indies = new ConcurrentHashMap<>();
    private final Map<String, Instance<Service>> services = new ConcurrentHashMap<>();
    private final Map<String, Instance<Reference>> references = new ConcurrentHashMap<>();

    @Override
    public Object define(MPI mpi, Class<?> reference) {
        return makeConsumer(mpi, reference).getProxy();
    }

    @Override
    public Execution<Reference> refer(MPI mpi, Class<?> reference, Inspector inspector) {
        return makeReferenceExecution(mpi, reference, getEnviron()).get(inspector.getSignature());
    }

    @Override
    public void store(Class<?> type, Object service) {
        providers.computeIfAbsent(type, key -> new Provider(type, service));
    }

    @Override
    public Execution<Service> infer(String urn) {
        URN domain = URN.from(urn);
        Map<String, Instance<Service>> executions = makeServiceExecution(getEnviron());
        return Optional.ofNullable(executions.get(domain.getName())).orElseGet(() -> {
            for (Map.Entry<String, Instance<Service>> instance : executions.entrySet()) {
                if (domain.matchName(instance.getKey())) {
                    return instance.getValue();
                }
            }
            return null;
        });
    }

    @Override
    public List<Class<?>> referTypes() {
        return new ArrayList<>(consumers.keySet());
    }

    @Override
    public List<Class<?>> inferTypes() {
        return providers.values().stream().map(Provider::getKind).collect(Collectors.toList());
    }

    private Environ getEnviron() {
        Network network = ServiceLoader.load(Network.class).getDefault();
        return network.getEnviron();
    }

    private Consumer makeConsumer(MPI mpi, Class<?> reference) {
        return consumers.computeIfAbsent(reference, key -> new ConcurrentHashMap<>()).computeIfAbsent(mpi, key -> new Consumer(mpi, reference));
    }

    private Map<String, Instance<Reference>> makeReferenceExecution(MPI mpi, Class<?> reference, Environ environ) {
        Consumer consumer = makeConsumer(mpi, reference);
        return indies.computeIfAbsent(reference, key -> new ConcurrentHashMap<>()).computeIfAbsent(mpi, key -> {
            Map<String, Instance<Reference>> executions = new ConcurrentHashMap<>();
            getDeclaredMethods(reference).forEach((kind, methods) -> methods.forEach(method -> {
                StringBuilder signature = new StringBuilder(method.getName());
                for (Parameter parameter : method.getParameters()) {
                    signature.append(parameter.getType().getName());
                }
                Reference refer = makeMethodAsReference(environ, mpi, kind, method);
                Instance<Reference> instance = new Instance<>(refer.getUrn(), reference, method, consumer.getProxy(), refer);
                references.put(refer.getUrn(), instance);
                executions.put(signature.toString(), instance);
            }));
            return executions;
        });
    }

    private Map<String, Instance<Service>> makeServiceExecution(Environ environ) {
        providers.values().stream().filter(x -> !x.isMaking()).forEach(provider -> {
            Class<?> type = provider.getKind();
            Object service = provider.getService();
            getDeclaredMethods(type).forEach((kind, methods) -> methods.forEach(method -> {
                List<Service> schemas = makeSchema(environ, type, kind, method);
                for (Service schema : schemas) {
                    services.put(URN.from(schema.getUrn()).getName(), new Instance<>(schema.getUrn(), type, method, service, schema));
                }
            }));
            provider.setMaking(true);
        });
        return services;
    }

    /**
     * Build service schema.
     *
     * @param environ Environ
     * @param type    Class type.
     * @param kind    Interface type.
     * @param method  Invoked method.
     * @return Schemas
     */
    private List<Service> makeSchema(Environ environ, Class<?> type, Class<?> kind, Method method) {
        List<Service> schemas = new ArrayList<>();
        MPS mps = type.getAnnotation(MPS.class);
        if (null != mps) {
            Service service = makeMethodAsService(mps, kind, method, environ);
            if (!Tool.MESH_MODE.get().match(Mode.NoCommunal) || (service.getFlags() & 4) != 4) {
                schemas.add(service);
            }
        }
        if (Subscriber.class != kind) {
            return schemas;
        }
        Bindings bindings = type.getAnnotation(Bindings.class);
        if (null != bindings) {
            for (Binding binding : bindings.value()) {
                if (binding.meshable()) {
                    schemas.add(makeMethodAsService(binding, kind, method, environ));
                }
            }
        }
        Binding binding = type.getAnnotation(Binding.class);
        if (null != binding && binding.meshable()) {
            schemas.add(makeMethodAsService(binding, kind, method, environ));
        }
        return schemas;
    }

    private <T> T chooseIfPresent(Function<MPI, T> fn, T defaults, MPI... tags) {
        for (MPI tag : tags) {
            if (null != tag && Tool.required(fn.apply(tag)) && !Tool.equals(fn.apply(tag), defaults)) {
                return fn.apply(tag);
            }
        }
        return defaults;
    }

    private <T> Stream<T> choose(Function<MPI, T> fn, T defaults, MPI... tags) {
        List<T> vs = new ArrayList<>();
        for (MPI tag : tags) {
            if (null != tag && Tool.required(fn.apply(tag)) && !Tool.equals(fn.apply(tag), defaults)) {
                vs.add(fn.apply(tag));
            }
        }
        vs.add(defaults);
        return vs.stream();
    }

    private Reference makeMethodAsReference(Environ environ, MPI mpi, Class<?> type, Method method) {
        MPI om = method.getAnnotation(MPI.class);
        MPI ok = type.getAnnotation(MPI.class);
        String value = chooseIfPresent(MPI::value, "", om);
        String name = chooseIfPresent(MPI::name, "", om);
        String proto = chooseIfPresent(MPI::proto, MeshFlag.GRPC.getName(), mpi, om, ok);
        String codec = chooseIfPresent(MPI::codec, MeshFlag.JSON.getName(), mpi, om, ok);
        String alias = Tool.anyone(value, name, Tool.toLowerCase(String.format("%s.%s", type.getName(), method.getName())));
        Reference reference = new Reference();
        reference.setNamespace(method.getDeclaringClass().getName());
        reference.setName(method.getName());
        reference.setVersion(chooseIfPresent(MPI::version, "", mpi, om, ok));
        reference.setProto(Tool.anyone(proto, MeshFlag.GRPC.getName()));
        reference.setCodec(Tool.anyone(codec, MeshFlag.JSON.getName()));
        reference.setFlags(chooseIfPresent(MPI::flags, 0L, om, ok) | mpi.flags());
        reference.setTimeout(choose(MPI::timeout, DEFAULT_TIMEOUT.toMillis(), mpi, om, ok).max(Long::compare).orElse(DEFAULT_TIMEOUT.toMillis()));
        reference.setRetries(choose(MPI::retries, 3, mpi, om, ok).max(Integer::compare).orElse(3));
        reference.setNode(chooseIfPresent(MPI::node, "", mpi, om, ok));
        reference.setInst(chooseIfPresent(MPI::inst, "", mpi, om, ok));
        reference.setZone(chooseIfPresent(MPI::zone, Tool.MESH_ZONE.get(), mpi, om, ok));
        reference.setCluster(chooseIfPresent(MPI::cluster, Tool.MESH_CLUSTER.get(), mpi, om, ok));
        reference.setCell(chooseIfPresent(MPI::cell, Tool.MESH_CELL.get(), mpi, om, ok));
        reference.setGroup(chooseIfPresent(MPI::group, Tool.MESH_GROUP.get(), mpi, om, ok));
        reference.setAddress(chooseIfPresent(MPI::address, "", mpi, om, ok));
        reference.setUrn(getURN(alias, reference, environ));
        return reference;
    }

    private Service makeMethodAsService(MPS mps, Class<?> kind, Method method, Environ environ) {
        MPS md = Optional.ofNullable(method.getAnnotation(MPS.class)).orElse(mps);
        MPI om = method.getAnnotation(MPI.class);
        MPI ok = kind.getAnnotation(MPI.class);
        String value = chooseIfPresent(MPI::value, "", om).replace("${mesh.name}", Tool.MESH_NAME.get()).replace("${mesh.uname}", md.name());
        String name = chooseIfPresent(MPI::name, "", om).replace("${mesh.name}", Tool.MESH_NAME.get()).replace("${mesh.uname}", md.name());
        String proto = chooseIfPresent(MPI::proto, Tool.anyone(md.proto(), MeshFlag.GRPC.getName()), om, ok);
        String codec = chooseIfPresent(MPI::codec, Tool.anyone(md.codec(), MeshFlag.JSON.getName()), om, ok);
        String alias = Tool.anyone(value, name, Tool.toLowerCase(String.format("%s.%s", kind.getName(), method.getName())));
        Service service = initServiceEnviron(environ);
        service.setNamespace(method.getDeclaringClass().getName());
        service.setName(md.name());
        service.setVersion(md.version());
        service.setProto(proto);
        service.setCodec(codec);
        service.setFlags(chooseIfPresent(MPI::flags, 0L, om, ok) | md.flags());
        service.setTimeout(choose(MPI::timeout, md.timeout(), om, ok).max(Long::compare).orElse(DEFAULT_TIMEOUT.toMillis()));
        service.setKind("MPS");
        service.setUrn(getURN(alias, service, environ));
        return service;
    }

    private Service initServiceEnviron(Environ environ) {
        Service service = new Service();
        service.setRetries(0);
        service.setNode(environ.getNodeId());
        service.setInst(environ.getInstId());
        service.setZone(Optional.ofNullable(environ.getDistribution()).map(Distribution::getZone).orElse(Tool.MESH_ZONE.get()));
        service.setCluster(Optional.ofNullable(environ.getDistribution()).map(Distribution::getCluster).orElse(Tool.MESH_CLUSTER.get()));
        service.setCell(Optional.ofNullable(environ.getDistribution()).map(Distribution::getCell).orElse(Tool.MESH_CELL.get()));
        service.setGroup(Optional.ofNullable(environ.getDistribution()).map(Distribution::getGroup).orElse(Tool.MESH_GROUP.get()));
        service.setSets(Tool.MESH_NAME.get());
        service.setAddress(getListenedAddress());
        service.setLang("Java");
        service.setAttrs(new HashMap<>(0));
        return service;
    }

    private Service makeMethodAsService(Binding binding, Class<?> kind, Method method, Environ environ) {
        Binding metadata = Optional.ofNullable(method.getAnnotation(Binding.class)).orElse(binding);
        String proto = Tool.anyone(metadata.proto(), MeshFlag.GRPC.getName());
        String codec = Tool.anyone(metadata.codec(), MeshFlag.JSON.getName());
        String alias = Tool.anyone(Tool.toLowerCase(String.format("%s.%s", metadata.topic(), metadata.code())), Tool.toLowerCase(String.format("%s.%s", kind.getName(), method.getName())));
        Service service = initServiceEnviron(environ);
        service.setNamespace(metadata.topic());
        service.setName(metadata.code());
        service.setVersion(metadata.version());
        service.setProto(proto);
        service.setCodec(codec);
        service.setFlags(metadata.flags());
        service.setTimeout(metadata.timeout() < 10 ? DEFAULT_TIMEOUT.toMillis() : metadata.timeout());
        service.setKind("Binding");
        service.setUrn(getURN(alias, service, environ));
        return service;
    }

    private Map<Class<?>, List<Method>> getDeclaredMethods(Class<?> type) {
        List<Class<?>> kinds = new ArrayList<>(type.getInterfaces().length + 1);
        kinds.addAll(Arrays.asList(type.getInterfaces()));
        if (type.isInterface()) {
            kinds.add(type);
        }
        return kinds.stream().collect(Collectors.toMap(kind -> kind, kind -> Arrays.stream(kind.getDeclaredMethods()).filter(method -> method.getDeclaringClass() != Object.class && Modifier.isPublic(method.getModifiers()) && !Modifier.isStatic(method.getModifiers())).collect(Collectors.toList())));
    }

    /**
     * Get the definition or reference definition.
     *
     * @param alias      Service alias name.
     * @param definition URN flags.
     * @param environ    Net environ
     * @return Unique uniform resource domain name.
     */
    private String getURN(String alias, LDCable definition, Environ environ) {
        URN urn = new URN();
        urn.setDomain(URN.MESH_DOMAIN);
        urn.setNodeId(Tool.toLowerCase(Tool.anyone(definition.getNode(), definition.getInst(), environ.getNodeId())));
        urn.setName(alias);
        urn.setFlag(getURNFlag(definition));
        return urn.toString();
    }

    /**
     * Get the reference or reference definition.
     *
     * @param definition Service definition.
     * @return Unique uniform resource domain name.
     */
    private URNFlag getURNFlag(LDCable definition) {
        String[] authority = Tool.split(definition.getAddress(), ":");
        URNFlag flag = new URNFlag();
        flag.setV("00");
        flag.setProto(MeshFlag.ofName(definition.getProto()).getCode());
        flag.setCodec(MeshFlag.ofName(definition.getCodec()).getCode());
        flag.setVersion(definition.getVersion());
        flag.setZone(definition.getZone());
        flag.setCluster(definition.getCluster());
        flag.setCell(definition.getCell());
        flag.setGroup(definition.getGroup());
        flag.setAddress(authority.length > 0 ? authority[0] : "");
        flag.setPort(authority.length > 1 ? authority[1] : "");
        return flag;
    }

    private String getListenedAddress() {
        // For ICBC PaaS
        String xh = Tool.MESH_RUNTIME_IP.get();
        String xp = Tool.MESH_RUNTIME_PORT.get();
        if (Tool.required(xh) && Tool.required(xp)) {
            return String.format("%s:%s", xh, xp);
        }
        return String.format("%s:%d", Tool.MESH_RUNTIME.get().getHost(), Tool.MESH_RUNTIME.get().getPort());
    }

    private Registration newRegistration(Metadata metadata, Map<String, String> attachments) {
        Registration registration = new Registration();
        registration.setInstanceId(getListenedAddress());
        registration.setContent(metadata);
        registration.setKind(Registration.METADATA);
        registration.setAddress(getListenedAddress());
        registration.setName(Tool.MESH_NAME.get());
        registration.setTimestamp(System.currentTimeMillis());
        registration.setAttachments(attachments);
        return registration;
    }

    @Override
    public void stop() throws Throwable {
        Scheduler scheduler = ServiceLoader.load(Scheduler.class).getDefault();
        scheduler.cancel(aliveId);
        scheduler.shutdown(Duration.ofSeconds(3));
        Registry registry = ServiceLoader.load(Registry.class).getDefault();
        Metadata metadata = new Metadata();
        metadata.setReferences(Collections.emptyList());
        metadata.setServices(Collections.emptyList());
        Tool.MESH_ADDRESS.get().getServers().parallelStream().forEach(addr -> Mesh.contextSafeUncheck(() -> {
            try {
                Mesh.context().setAttribute(Mesh.REMOTE, addr.getAddress());
                registry.unregister(newRegistration(metadata, new HashMap<>(0)));
            } catch (Exception e) {
                log.warn("Register to {} with {}:{}", addr.getAddress(), e.getClass().getName(), e.getMessage());
            }
        }));
    }

    @Override
    public void refresh() throws Throwable {
        start();
        Scheduler scheduler = ServiceLoader.load(Scheduler.class).getDefault();
        aliveId = scheduler.period(Duration.ofSeconds(30), new Topic(MeshEden.TOPIC, MeshEden.CODE));
    }

    @Override
    public void subscribe(Event event) {
        Builtin builtin = ServiceLoader.load(Builtin.class).getDefault();
        Versions versions = Optional.ofNullable(builtin.version()).orElseGet(Versions::new);
        Map<String, String> attachments = new HashMap<>();
        attachments.put("_PAAS_PORT_7700", Tool.getProperty("", "_PAAS_PORT_7700"));
        attachments.putAll(Optional.ofNullable(versions.getInfos()).orElseGet(Collections::emptyMap));
        Context.Metadata.MESH_SUBSET.set(attachments, Tool.MESH_SUBSET.get());
        Environ environ = getEnviron();
        makeServiceExecution(environ);
        consumers.values().stream().flatMap(x -> x.values().stream()).forEach(consumer -> makeReferenceExecution(consumer.getMpi(), consumer.getReference(), environ));
        Registry registry = ServiceLoader.load(Registry.class).getDefault();
        Metadata metadata = new Metadata();
        metadata.setReferences(references.values().stream().map(Instance::getResource).collect(Collectors.toList()));
        metadata.setServices(services.values().stream().map(Instance::getResource).collect(Collectors.toList()));

        Tool.MESH_ADDRESS.get().getServers().parallelStream().forEach(addr -> Mesh.contextSafeUncheck(() -> {
            try {
                Mesh.context().setAttribute(Mesh.REMOTE, addr.getAddress());
                registry.register(newRegistration(metadata, attachments));
            } catch (Exception e) {
                log.warn("Register to {} with {}:{}", addr.getAddress(), e.getClass().getName(), e.getMessage());
            }
        }));
    }

    @Getter
    private static final class Instance<T> implements Execution<T> {

        private final String urn;
        private final Class<?> type;
        private final Method method;
        private final Object target;
        private final T resource;
        private final Inspector inspector;
        private final Invoker<?> invoker;
        private final Class<? extends Parameters> intype;
        private final Class<? extends Returns> retype;

        private Instance(String urn, Class<?> type, Method method, Object target, T resource) {
            this.urn = urn;
            this.type = type;
            this.method = method;
            this.target = target;
            this.resource = resource;
            this.inspector = new MethodInspector(type, method);
            this.invoker = new ServiceInvokeHandler(target);
            this.intype = ServiceLoader.load(JCompiler.class).get(JCompiler.JAVASSIST).intype(this.method);
            this.retype = ServiceLoader.load(JCompiler.class).get(JCompiler.JAVASSIST).retype(this.method);
        }

        @Override
        public T schema() {
            return this.resource;
        }

        @Override
        public Inspector inspect() {
            return this.inspector;
        }

        @Override
        public <I extends Parameters> Types<I> intype() {
            return Types.of((Type) this.intype);
        }

        @Override
        public <O extends Returns> Types<O> retype() {
            return Types.of((Type) this.retype);
        }

        @Override
        public Parameters inflect() {
            return Tool.newInstance(this.intype);
        }

        @Override
        public Returns reflect() {
            return Tool.newInstance(this.retype);
        }

        @Override
        public Object invoke(Invocation invocation) throws Throwable {
            return this.invoker.invoke(invocation);
        }
    }

    @Getter
    private static final class Consumer {
        private final MPI mpi;
        private final Class<?> reference;
        private final Object proxy;

        public Consumer(MPI mpi, Class<?> reference) {
            this.mpi = mpi;
            this.reference = reference;
            this.proxy = ServiceProxy.proxy(mpi, reference);
        }
    }

    @Getter
    private static final class Provider {
        @Setter
        private boolean making;
        private final Class<?> kind;
        private final Object service;

        public Provider(Class<?> kind, Object service) {
            this.kind = kind;
            this.service = service;
        }
    }
}
