/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.codec;

import com.be.mesh.client.annotate.Format;
import com.be.mesh.client.cause.CompatibleException;
import com.be.mesh.client.codec.proto.ProtoTestService;
import com.be.mesh.client.codec.proto.ProtobufTest;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.struct.Inbound;
import com.be.mesh.client.struct.Outbound;
import com.google.common.collect.ImmutableMap;
import com.google.protobuf.ByteString;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.testng.Assert;
import org.testng.annotations.Test;

import java.lang.reflect.Method;
import java.lang.reflect.Type;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
@Slf4j
public class CodecTest {

    private Date x = new Date();
    private LocalDate date = LocalDate.now();
    private LocalTime time = LocalTime.now();
    @Format(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime datetime = LocalDateTime.now();
    private LocalDateTime timestamp = LocalDateTime.now();
    private byte[] list = "mesh</>^;&#$!@*(&)[]{}=''``?/\\".getBytes(StandardCharsets.UTF_8);
    private byte[][] matrix = new byte[][]{{1, 2, 3}, "unicode \uD83D\uDE08".getBytes(StandardCharsets.UTF_8)};
    private byte[][][] z = new byte[][][]{{{1, 2, 3}}};

    @Test
    public void nativeCodecTest() throws Exception {
        Codec json = ServiceLoader.load(Codec.class).get(Codec.JSON);
        String gson = json.encodeString(new CodecTest());
        log.info("{}", gson);
        log.info("{}", json.decodeString(gson, Types.of(CodecTest.class)));
        log.info(new String(json.decodeString(gson, Types.of(CodecTest.class)).list, StandardCharsets.UTF_8));

        Codec jack = ServiceLoader.load(Codec.class).get(Codec.JACKSON);
        String jackson = jack.encodeString(new CodecTest());
        log.info("{}", jackson);
        log.info("{}", jack.decodeString(jackson, Types.of(CodecTest.class)));
        log.info(new String(jack.decodeString(jackson, Types.of(CodecTest.class)).list, StandardCharsets.UTF_8));
    }

    @Test
    public void unicodeTest() throws Exception {
        Map<String, Object> payload = new HashMap<>();
        payload.put("x", "utf8 ^-^");
        payload.put("y", "unicode \uD83D\uDE08");
        payload.put("z", new byte[0]);
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JACKSON);
        log.info("{}", codec.encodeString(payload));
        log.info("{}", codec.decode(codec.encode(payload), Types.MapObject));
    }

    @Test
    public void localDateTimeTest() throws Exception {
        CodecTest test = new CodecTest();
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JACKSON);
        log.info("{}", codec.encodeString(test));

        CodecTest x = codec.decodeString("{\"date\":null,\"time\":null,\"datetime\":\"2011-11-11 11:11:11\"}", Types.of(CodecTest.class));
        log.info("{}", x);

        CodecTest y = codec.decodeString("{\"date\":null,\"time\":null,\"datetime\":\"1645583734520\"}", Types.of(CodecTest.class));
        log.info("{}", y);
        log.info("{}", codec.decode(codec.encode(test), Types.of(CodecTest.class)));

        Map<String, Object> z = new HashMap<>();
        z.put("x", null);
        log.info("{}", new String(codec.encode(z).array()));

        Map<String, Object> k = codec.decode(ByteBuffer.wrap("{\"x\":null}".getBytes(StandardCharsets.UTF_8)), Types.MapObject);
        log.info("{}", String.join(",", k.keySet()));
    }

    @Test
    public void gsonDateTimeTest() throws Exception {
        CodecTest test = new CodecTest();
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        log.info("{}", codec.encodeString(test));
        log.info("{}", codec.decodeString(codec.encodeString(test), Types.of(CodecTest.class)));

        CodecTest x = codec.decodeString("{\"date\":null,\"time\":null,\"datetime\":\"2011-11-11 11:11:11\"}", Types.of(CodecTest.class));
        log.info("{}", x);

        CodecTest y = codec.decodeString("{\"date\":null,\"time\":null,\"timestamp\":\"1645583734520\"}", Types.of(CodecTest.class));
        log.info("{}", y);

        log.info("{}", codec.decode(codec.encode(test), Types.of(CodecTest.class)));
    }

    @Test
    public void exceptionTest() throws Exception {
        RuntimeException e = new CompatibleException("Test x");
        Cause cause = Cause.of(e);
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        Cause clone = codec.decode(codec.encode(cause), Types.of(Cause.class));
        Assert.assertEquals(cause.getBuff(), clone.getBuff());
        log.info("", Cause.of("", "", clone));
        Assert.assertTrue(Cause.of("", "", clone) instanceof CompatibleException);
    }

    @Test
    public void protobufTest() throws Exception {
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.PROTOBUF);
        Inbound parameter = new Inbound();
        parameter.setAttachments(ImmutableMap.of("1", "2"));
        ByteBuffer buffer = codec.encode(parameter);
        Inbound x = codec.decode(buffer, Types.of(Inbound.class));
        log.info("{}", x);
    }

    @Test
    public void protobuf4Test() throws Exception {
        // Init
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.PROTOBUF4);
        Codec jsonCodec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        JCompiler javassistCompiler = ServiceLoader.load(JCompiler.class).get(JCompiler.JAVASSIST);
        Method mockProtoMethod = ProtoTestService.class.getMethods()[0];

        /**************Outbound**************/
        // Encode
        Outbound outbound = new Outbound();
        MeshCode systemError = MeshCode.SYSTEM_ERROR;
        outbound.setCode(systemError.getCode());
        outbound.setMessage(MeshCode.SYSTEM_ERROR.getMessage());
        outbound.setCause(Cause.of(new RuntimeException("mock exception")));
        ByteBuffer outboundEncodeBuffer = codec.encode(outbound);
        log.info("Outbound Encode: {}", new String(outboundEncodeBuffer.array()));
        // Decode
        Outbound outboundDecode = codec.decode(outboundEncodeBuffer, Types.of(Outbound.class));
        log.info("Outbound Decode: {}", new String(jsonCodec.encode(outboundDecode).array()));

        /**************Parameters**************/
        // Encode
        Class<Parameters> mockParamsClazz = javassistCompiler.intype(mockProtoMethod);
        Parameters parameters = mockParamsClazz.getDeclaredConstructor().newInstance();
        ProtobufTest.PsiExecuteRequest.Builder parametersBuilder = ProtobufTest.PsiExecuteRequest.newBuilder();
        parametersBuilder.setAuthorityCode("11111");
        parametersBuilder.setSourcePartnerCode("2222");
        parametersBuilder.setTaskId("3333");
        parametersBuilder.setIndex("4444");
        parametersBuilder.addEncodeData(ByteString.copyFrom("5555".getBytes(StandardCharsets.UTF_8)));
        ProtobufTest.PsiExecuteRequest requestParam = parametersBuilder.build();
        parameters.arguments(new Object[]{requestParam});
        parameters.attachments(new HashMap<>());
        ByteBuffer parametersEncodeBuffer = codec.encode(parameters);
        log.info("Parameters Encode: {}", new String(parametersEncodeBuffer.array()));
        // Decode
        Parameters parametersDecode = codec.decode(parametersEncodeBuffer, Types.of((Type) mockParamsClazz));
        log.info("Parameters Decode: {}", new String(jsonCodec.encode(parametersDecode).array()));

        /**************Returns**************/
        // Encode
        Class<Returns> mockContentClazz = javassistCompiler.retype(mockProtoMethod);
        Returns returns = mockContentClazz.getDeclaredConstructor().newInstance();
        returns.setCode(MeshCode.SUCCESS.getCode());
        returns.setMessage(MeshCode.SUCCESS.getMessage());
        ProtobufTest.PsiExecuteResponse.Builder returnsBuilder = ProtobufTest.PsiExecuteResponse.newBuilder();
        returnsBuilder.addEncoders(ByteString.copyFrom("11111".getBytes(StandardCharsets.UTF_8)));
        returnsBuilder.addMaskData(ByteString.copyFrom("22222".getBytes(StandardCharsets.UTF_8)));
        returnsBuilder.setTaskId("3333");
        returnsBuilder.setStage("4444");
        returns.setContent(returnsBuilder.build());
        ByteBuffer returnsEncodeBuffer = codec.encode(returns);
        log.info("Returns Encode: {}", new String(returnsEncodeBuffer.array()));
        // Decode
        Returns returnsDecode = codec.decode(returnsEncodeBuffer, Types.of((Type) mockContentClazz));
        log.info("Returns Decode: {}", new String(jsonCodec.encode(returnsDecode).array()));
    }

    @Test
    public void internalTypeTest() {
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        InternalType it = new InternalType();
        it.setTime(LocalDateTime.now());
        ByteBuffer buffer = codec.encode(it);
        InternalType.IK nit = codec.decode(buffer, Types.of(InternalType.IK.class));
        log.info("{}", nit.getTime());
    }

    @Data
    public static final class InternalType {
        private LocalDateTime time;

        @Data
        public static final class IK {
            private LocalDateTime time;
        }
    }
}
