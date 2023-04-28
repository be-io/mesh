/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.reflect;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.schema.compiler.JdkCompiler;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.tool.Tool;
import com.google.common.collect.ImmutableMap;
import javassist.ClassPool;
import javassist.CtClass;
import javassist.CtField;
import javassist.CtMethod;
import javassist.bytecode.ClassFile;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.lang.invoke.MethodHandles;
import java.lang.reflect.Method;
import java.lang.reflect.Type;
import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class ReflectTest {

    @SuppressWarnings("unchecked")
    @Test
    public void defineClassTest() throws Exception {
        try {
            if (ClassFile.MAJOR_VERSION <= ClassFile.JAVA_9) {
                return;
            }
            byte[] clazz = ("Êþº¾\0\0\0005\0\26\11\0\11\0\12\10\0\13\12\0\14\0" + "\15\7\0\16\7\0\17\1\0\3foo\1\0\3()V\1\0\4Code\7\0\20\14\0\21\0\22\1\0" + "\30hello from dynamic class\7\0\23\14\0\24\0\25\1\0\4Lazy\1\0\20java/" + "lang/Object\1\0\20java/lang/System\1\0\3out\1\0\25Ljava/io/PrintStream;" + "\1\0\23java/io/PrintStream\1\0\7println\1\0\25(Ljava/lang/String;)V\6\0" + "\0\4\0\5\0\0\0\0\0\1\0\11\0\6\0\7\0\1\0\10\0\0\0\25\0\2\0\0\0\0\0\11²\0" + "\1\22\2¶\0\3±\0\0\0\0\0\0").getBytes(StandardCharsets.ISO_8859_1);
            Class<?> lookupType = MethodHandles.Lookup.class;
            Method privateLookupIn = MethodHandles.class.getDeclaredMethod("privateLookupIn", Class.class, lookupType);
            Object lookup = privateLookupIn.invoke(null, Lazy.class, MethodHandles.lookup());
            Class<Lazy> lazyClass = (Class<Lazy>) lookupType.getDeclaredMethod("defineClass", byte[].class).invoke(lookup, clazz);
            Lazy.foo();
        } catch (Exception e) {
            log.error("", e);
        }
    }

    interface Lazy {
        static void foo() {
        }
    }

    @Test
    public void serializeCauseTest() throws Exception {
        byte[] buffer = Tool.serializeCause(new MeshException("Test Is nice!"));
        Throwable e = Tool.deserializeCause(buffer);
        e.printStackTrace();
        log.info("codec");
        Cause cause = new Cause();
        cause.setBuff(buffer);
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        String text = codec.encodeString(cause);
        Cause next = codec.decodeString(text, Types.of(Cause.class));
        Throwable ne = Tool.deserializeCause(next.getBuff());
        ne.printStackTrace();
    }

    @Test
    public void jacksonTest() throws Exception {
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JACKSON);
        ResponseWrap wrap = codec.decode(codec.encode(new ResponseWrap(new Response())), Types.of(ResponseWrap.class));
        log.info("{}", wrap);
    }

    @Test
    public void javacTest() throws Throwable {
        CtClass ctClass = ClassPool.getDefault().makeClass("Testing");
        ctClass.addField(CtField.make("private boolean x;", ctClass));
        ctClass.addField(CtField.make("private long y;", ctClass));
        ctClass.addMethod(CtMethod.make("public void arguments(Object[] args){if(null != args){if(args.length > 0){this.x=(boolean)args[0];}if(args.length > 1){this.y=(long)args[1];}}}", ctClass));
        ctClass.addMethod(CtMethod.make("public Object[] arguments(){return new Object[0];}", ctClass));

        new JdkCompiler().doCompile("com.be.mesh.client.reflect.ReflectTest$ITest$X3Input", Tool.read(getClass().getResourceAsStream("/Test.javac")));
        new JdkCompiler().doCompile("Test", "public class Test<V,B>{private V v;private B b;}");
    }

    @Test
    public void complexTypeTest() throws Exception {
        Map<String, Object[]> cases = new HashMap<>();
        cases.put("x", new Object[]{new Object[]{new Response()}, new Response()});
        cases.put("x0", new Object[]{new Object[]{null}, 0});
        cases.put("x1", new Object[]{new Object[]{new byte[][]{{1}, {2}}}, new byte[][]{{3}, {4}}});
        cases.put("x2", new Object[]{new Object[]{0}, 1});
        cases.put("x3", new Object[]{new Object[]{true, 1L}, false});
        cases.put("x4", new Object[]{new Object[]{0}, 1});
        cases.put("x5", new Object[]{new Object[]{"1"}, new long[][]{{1}, {2}}});
        cases.put("x6", new Object[]{new Object[]{(byte) 1}, new String[]{"1", "2"}});
        cases.put("x7", new Object[]{new Object[]{Collections.singletonList(ImmutableMap.of("1", ImmutableMap.of("1", new Response())))}, Collections.singletonList(ImmutableMap.of("1", Collections.singletonList(new Response())))});
        cases.put("x8", new Object[]{new Object[]{1, (short) 1, (long) 1, (float) 1, (double) 1, true, '1', (byte) 1, new byte[]{1}, new byte[][]{{1}}}, new String[]{"1", "2"}});
        cases.put("x9", new Object[]{new Object[]{LocalDateTime.now(), LocalDate.now(), LocalTime.now(), new Date()}, LocalDateTime.now()});
        cases.put("x10", new Object[]{new Object[]{new byte[][][]{{{1, 2, 3}}}}, new String[][][]{{{"1", "2"}}}});
        cases.put("x11", new Object[]{new Object[]{new Integer[][][]{{{1, 2, 3}}}}, new float[][][]{{{1.1f, 2.1f}}}});
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        JCompiler loader = ServiceLoader.load(JCompiler.class).get(JCompiler.JAVASSIST);
        for (Method method : ITest.class.getDeclaredMethods()) {
            try {
                log.info("Load {}", method.toGenericString());
                log.info("{}", loader.intype(method).getName());
                log.info("{}", loader.retype(method).getName());

                Parameters parameters = loader.intype(method).getConstructor().newInstance();
                parameters.arguments((Object[]) cases.get(method.getName())[0]);
                Parameters x = codec.decode(codec.encode(parameters), Types.of(loader.intype(method)));

                Returns returns = loader.retype(method).getConstructor().newInstance();
                returns.setContent(cases.get(method.getName())[1]);
                Returns y = codec.decode(codec.encode(returns), Types.of(loader.retype(method)));

                log.info("{} {}", method.getName(), codec.encodeString(x));
                log.info("{} {}", method.getName(), codec.encodeString(y));
                log.info("{} {}", method.getName(), x.arguments().length == parameters.argumentTypes().size());
                log.info("{} {}", method.getName(), parameters.argumentTypes().values().stream().map(Type::getTypeName).collect(Collectors.joining("||")));
            } catch (Throwable e) {
                log.error("{}", method.getName());
                log.error("", e);
            }
        }
    }

    private interface ITest {

        Response x(Response i);

        void x0(boolean x);

        byte[][] x1(byte[][] y);

        int x2();

        Boolean x3(boolean x, long y);

        Integer x4(Object x);

        long[][] x5(String x);

        String[] x6(byte x);

        List<Map<String, List<Response>>> x7(List<Map<String, Map<String, Response>>> x);

        String[] x8(int q, short w, long e, float r, double t, boolean y, char u, byte i, byte[] o, byte[][] p);

        LocalDateTime x9(LocalDateTime x, LocalDate y, LocalTime z, Date v);

        String[][][] x10(byte[][][] x);

        float[][][] x11(Integer[][][] x);
    }

    @Data
    public static final class Response {
        private String name = "!";
    }

    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static final class ResponseWrap {
        private Response response;
    }
}
