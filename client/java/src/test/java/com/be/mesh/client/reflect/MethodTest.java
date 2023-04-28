/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.reflect;

import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.struct.Principal;
import lombok.extern.slf4j.Slf4j;
import org.junit.Assert;
import org.testng.annotations.Test;

import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.*;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class MethodTest {

    @SuppressWarnings("unchecked")
    @Test
    public void genericParametersTest() throws Exception {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        JCompiler classLoader = ServiceLoader.load(JCompiler.class).get(JCompiler.JAVASSIST);
        AmTesting testing = new AmTesting();
        for (Method method : testing.getClass().getDeclaredMethods()) {
            if (method.getName().equals("invoke0")) {
                Map<String, Object> parameters = new HashMap<>();
                parameters.put("input", true);
                parameters.put("attachments", new HashMap<>());
                ByteBuffer buffer = codec.encode(parameters);
                Class<?> type = classLoader.intype(method);
                Object t = codec.decode(buffer, Types.of(type));
                Assert.assertTrue(testing.invoke0((boolean) type.getDeclaredMethod("getInput").invoke(t)));
                continue;
            }
            if (method.getName().equals("invoke1")) {
                Map<String, Object> parameters = new HashMap<>();
                parameters.put("input", Boolean.TRUE);
                parameters.put("attachments", new HashMap<>());
                ByteBuffer buffer = codec.encode(parameters);
                Class<?> type = classLoader.intype(method);
                Object t = codec.decode(buffer, Types.of(type));
                Assert.assertTrue(testing.invoke1((Boolean) type.getDeclaredMethod("getInput").invoke(t)));
                continue;
            }
            if (method.getName().equals("invoke2")) {
                Map<String, Object> parameters = new HashMap<>();
                parameters.put("input", "X".getBytes(StandardCharsets.UTF_8));
                parameters.put("attachments", new HashMap<>());
                ByteBuffer buffer = codec.encode(parameters);
                Class<?> type = classLoader.intype(method);
                Object t = codec.decode(buffer, Types.of(type));
                Assert.assertTrue(testing.invoke2((byte[]) type.getDeclaredMethod("getInput").invoke(t)));
                continue;
            }
            for (Parameter parameter : method.getParameters()) {
                log.info("x-{}", parameter.getParameterizedType().getTypeName());
                log.info("z-{}", new ParameterizedTypes(parameter.getParameterizedType()).toSignature());
            }
            Map<String, Principal> map = new HashMap<>();
            map.put("y", new Principal("z", "c"));
            Map<String, List<Map<String, Principal>>> input = new HashMap<>();
            input.put("x", Collections.singletonList(map));
            Map<String, Object> parameters = new HashMap<>();
            parameters.put("input", input);
            parameters.put("attachments", new HashMap<>());
            ByteBuffer buffer = codec.encode(parameters);
            Class<?> type = classLoader.intype(method);
            Object t = codec.decode(buffer, Types.of(type));
            Assert.assertTrue(testing.invoke((Map<String, List<Map<String, Principal>>>) type.getDeclaredMethod("getInput").invoke(t), new HashSet<>()));
        }
    }

    private interface Testing<T> {
        boolean invoke(Map<String, List<T>> input, Set<String> v);

        boolean invoke0(boolean input);

        boolean invoke1(Boolean input);

        boolean invoke2(byte[] input);
    }

    private static class AmTesting implements Testing<Map<String, Principal>> {

        private boolean v;
        private byte[] x;

        @Override
        public boolean invoke(Map<String, List<Map<String, Principal>>> input, Set<String> v) {
            return true;
        }

        @Override
        public boolean invoke0(boolean input) {
            return true;
        }

        @Override
        public boolean invoke1(Boolean input) {
            return true;
        }

        @Override
        public boolean invoke2(byte[] input) {
            return input.length > 0;
        }
    }
}
