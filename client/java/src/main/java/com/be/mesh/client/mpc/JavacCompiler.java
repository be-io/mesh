/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.schema.compiler.JavaCompiler;
import com.be.mesh.client.schema.compiler.JdkCompiler;
import com.be.mesh.client.schema.runtime.TypeStruct;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.struct.Service;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;

import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.lang.reflect.Type;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Future;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI(JCompiler.JAVAC)
@SuppressWarnings("unchecked")
public class JavacCompiler implements JCompiler {

    private static final String PARAMETERS = "P";
    private static final String RETURNS = "R";
    private final Map<String, Class<?>> cache = new ConcurrentHashMap<>();
    private final JavaCompiler compiler = new JdkCompiler();

    @Override
    public <T extends Parameters> Class<T> intype(Method method) {
        if (method.getParameters().length < 1) {
            return (Class<T>) GenericParameters.class;
        }
        return (Class<T>) cache.computeIfAbsent(getSignature(method, PARAMETERS), key -> Tool.uncheck(() -> {
            log.info("Intype {}", key);
            String pn = method.getDeclaringClass().getPackage().getName();
            String name = getTypeName(method, PARAMETERS, false);
            String parameters = Tool.read(getClass().getResourceAsStream("/META-INF/mml/Parameters.mml"));
            StringBuilder vars = new StringBuilder();
            StringBuilder from = new StringBuilder();
            StringBuilder map = new StringBuilder();
            for (int index = 0; index < method.getParameters().length; index++) {
                Parameter parameter = method.getParameters()[index];
                int idx = Optional.ofNullable(parameter.getAnnotation(Index.class)).map(Index::value).orElse(index);
                String in = Optional.ofNullable(parameter.getAnnotation(Index.class)).map(Index::name).orElse("");
                String pt = new ParameterizedTypes(parameter.getParameterizedType()).toCanonicalName();
                String vn = parameter.getName();
                vars.append(String.format("@Index(value = %d, name = \"%s\")", idx, Tool.anyone(in, vn)));
                vars.append(String.format("private %s %s;", pt, vn));
                vars.append(String.format("public %s get%s(){return this.%s;}", pt, Tool.firstUpperCase(vn), vn));
                vars.append(String.format("public void set%s(%s %s){this.%s=%s;}", Tool.firstUpperCase(vn), pt, vn, vn, vn));
                from.append(String.format("if(args.length > %d){this.%s=(%s)args[%d];}", index, vn, pt, index));
                map.append(String.format("dict.put(\"%s\",this.%s)", vn, vn));
            }
            String to = Arrays.stream(method.getParameters()).map(Parameter::getName).collect(Collectors.joining(", "));
            return (Class<T>) compiler.compile(String.format(parameters, pn, name, vars, map, to, from), this.getClass().getClassLoader());
        }));
    }

    @Override
    public <T extends Parameters> Class<T> intype(Reference reference) {
        try {
            return Tool.uncheck(() -> intype(Class.forName(reference.getNamespace()).getMethod(reference.getName())));
        } catch (RuntimeException | Error e) {
            log.error("Compile with error for {}.{}", reference.getNamespace(), reference.getName());
            throw e;
        }
    }

    @Override
    public <T extends Returns> Class<T> retype(Method method) {
        return (Class<T>) cache.computeIfAbsent(getSignature(method, RETURNS), key -> Tool.uncheck(() -> {
            log.info("Retype {}", key);
            String pn = method.getDeclaringClass().getPackage().getName();
            String cn = getTypeName(method, RETURNS, false);
            Type returnType = Types.unbox(method, Future.class);
            boolean avoid = returnType == void.class || returnType == Void.class;
            String rt = avoid ? Object.class.getCanonicalName() : new ParameterizedTypes(returnType).toCanonicalName();
            String returns = Tool.read(getClass().getResourceAsStream("/META-INF/mml/Returns.mml"));
            return (Class<T>) compiler.compile(String.format(returns, pn, cn, rt, rt), this.getClass().getClassLoader());
        }));
    }

    @Override
    public <T extends Returns> Class<T> retype(Service service) {
        try {
            return Tool.uncheck(() -> retype(Class.forName(service.getNamespace()).getMethod(service.getName())));
        } catch (RuntimeException | Error e) {
            log.error("Compile with error for {}.{}", service.getNamespace(), service.getName());
            throw e;
        }
    }

    @Override
    public List<TypeStruct> documents() {
        return new ArrayList<>();
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T> Class<? extends T> compile(Class<T> interfaces, String implement) {
        return (Class<T>) cache.computeIfAbsent(implement, key -> Tool.uncheck(() -> {
            try {
                return compiler.compile(implement, this.getClass().getClassLoader());
            } catch (Throwable e) {
                log.error("Compile with error for {}.{}", interfaces.getCanonicalName(), implement);
                throw e;
            }
        }));
    }

}
