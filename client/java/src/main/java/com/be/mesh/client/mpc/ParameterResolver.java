/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.tool.UUID;
import javassist.ClassClassPath;
import javassist.ClassPool;
import javassist.CtClass;
import javassist.CtField;
import lombok.extern.slf4j.Slf4j;

import java.io.Serializable;
import java.lang.reflect.Field;
import java.lang.reflect.Parameter;
import java.nio.ByteBuffer;
import java.util.Arrays;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * 解析byte[]参数报文匹配为接口参数列表.
 *
 * @author coyzeng@gmail.com
 */
@Slf4j(topic = "rpc-param-digest")
public class ParameterResolver {

    /**
     * 参数列表类缓存
     */
    private final Map<String, Class<?>> PCC = new ConcurrentHashMap<>();

    /**
     * 参数属性缓存
     */
    private final Map<String, Field> PFC = new ConcurrentHashMap<>();

    private static class Header implements Serializable {

        private static final long serialVersionUID = -4198557639963697677L;

        private String method;

        private String version;

        private String address;
    }

    public Object[] resolveArguments(byte[] body, Parameter[] parameters) throws Exception {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        Header header = codec.decode(ByteBuffer.wrap(body), Types.of(Header.class));
        Class<?> clazz = PCC.computeIfAbsent(header.method, key -> {
            try {
                return defineClass("Decorate" + UUID.getInstance().shortUUID(), parameters);
            } catch (RuntimeException | Error e) {
                log.warn("{} build parameter decoration failed. cause={}", header.method, e.getMessage());
                throw e;
            } catch (Throwable e) {
                log.error("{} build parameter decoration failed. cause={}", header.method, e.getMessage());
                throw new MeshException(e);
            }
        });
        Object payload = codec.decode(ByteBuffer.wrap(body), Types.of(clazz));
        Field param = clazz.getDeclaredField("param");
        Object value = param.get(payload);
        return Arrays.stream(parameters).map(parameter -> {
            Field field = PFC.computeIfAbsent((header.method + parameter.getName()).intern(), key -> {
                try {
                    return value.getClass().getField(parameter.getName());
                } catch (Exception e) {
                    log.error(
                            "{} parameter {} resolve failed. cause={}",
                            header.method,
                            parameter.getName(),
                            e.getMessage());
                    throw new MeshException(e, "%s illegal", parameter.getName());
                }
            });
            try {
                return field.get(value);
            } catch (IllegalAccessException e) {
                log.error("{} parameter {} read failed. cause={}", header.method, parameter.getName(), e.getMessage());
                throw new MeshException(e, "%s illegal", parameter.getName());
            }
        }).toArray();
    }

    private Class<?> defineClass(String className, Parameter[] parameters) throws Exception {
        if (null == parameters || parameters.length < 1) {
            return Header.class;
        }

        ClassPool pool = ClassPool.getDefault();
        pool.insertClassPath(new ClassClassPath(this.getClass()));
        CtClass paramClass = pool.makeClass(className);
        for (Parameter parameter : parameters) {
            String field = String.format("public %s %s;", parameter.getType().getName(), parameter.getName());
            CtField cf = CtField.make(field, paramClass);
            paramClass.addField(cf);
        }
        CtClass wrapClass = pool.makeClass(className + "DTO");
        wrapClass.addField(CtField.make("public String method;", wrapClass));
        wrapClass.addField(CtField.make("public String version;", wrapClass));
        wrapClass.addField(CtField.make("public String address;", wrapClass));
        wrapClass.addField(CtField.make(String.format("public %s param;", paramClass.toClass().getName()), wrapClass));
        return wrapClass.toClass();
    }
}
