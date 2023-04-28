package com.be.mesh.client.codec;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.codec.proto.Protobuf4;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.tool.Tool;
import com.google.protobuf.ByteString;
import com.google.protobuf.Message;
import org.springframework.util.CollectionUtils;

import java.lang.reflect.Method;
import java.lang.reflect.Type;
import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

/**
 * @author nanfeng
 * @date 2023/2/14 2:44 PM
 */
@SPI(Codec.PROTOBUF4)
public class Protobuf4Codec implements Codec {

    @Override
    public ByteBuffer encode(Object value) {
        return this.encode0(value, object -> {
            Class<?> clazz = object.getClass();
            if (Parameters.class.isAssignableFrom(clazz)) {
                return encodeParameters((Parameters) object);
            }
            if (Returns.class.isAssignableFrom(clazz)) {
                return encodeReturns((Returns) object);
            }
            throw new MeshException(MeshCode.CRYPT_CODEC_ERROR, new RuntimeException("illegal proto codec type"));
        });
    }

    private ByteBuffer encodeParameters(Parameters parameters) {
        Protobuf4.InBound.Builder builder = Protobuf4.InBound.newBuilder();
        builder.putAllAttachments(parameters.attachments());
        Map<Integer, ByteString> arguments = new HashMap<>();
        if (parameters.arguments() != null) {
            for (int idx = 0; idx < parameters.arguments().length; idx++) {
                Object argument = parameters.arguments()[idx];
                if (argument instanceof Message) {
                    arguments.put(idx, ByteString.copyFrom(((Message) argument).toByteArray()));
                }
            }
        }
        builder.putAllArguments(arguments);
        return ByteBuffer.wrap(builder.build().toByteArray());
    }

    private ByteBuffer encodeReturns(Returns returns) {
        Protobuf4.OutBound.Builder builder = Protobuf4.OutBound.newBuilder();
        builder.setCode(Optional.ofNullable(returns.getCode()).orElse(""));
        builder.setMessage(Optional.ofNullable(returns.getMessage()).orElse(""));
        Cause ca = returns.getCause();
        if (ca != null) {
            Protobuf4.Cause.Builder causer = Protobuf4.Cause.newBuilder();
            causer.setName(Optional.ofNullable(ca.getName()).orElse(""));
            causer.setPos(Optional.ofNullable(ca.getPos()).orElse(""));
            causer.setText(Optional.ofNullable(ca.getText()).orElse(""));
            causer.setBuff(ByteString.copyFrom(Optional.ofNullable(ca.getBuff()).orElse(new byte[]{})));
            builder.setCause(causer.build());
        }
        if (returns.getContent() instanceof Message) {
            builder.setContent(ByteString.copyFrom(((Message) returns.getContent()).toByteArray()));
        }
        return ByteBuffer.wrap(builder.build().toByteArray());
    }


    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        return this.decode0(buffer, type, (x, y) -> {
            try {
                Class<?> clazz = (Class<?>)y.getRawType();
                if (Parameters.class.isAssignableFrom(clazz)) {
                    return decodeParameters(x, y);
                }
                if (Returns.class.isAssignableFrom(clazz)) {
                    return decodeReturns(x, y);
                }
                throw new MeshException(MeshCode.CRYPT_CODEC_ERROR, new RuntimeException("illegal proto codec type"));
            } catch (Exception e) {
                throw new MeshException(e);
            }
        });
    }

    @SuppressWarnings("unchecked")
    private <T> T decodeParameters(ByteBuffer buffer, Types<T> type) throws Exception {
        Protobuf4.InBound inBound = Protobuf4.InBound.parseFrom(buffer);
        Class<Parameters> parametersClazz = (Class<Parameters>) type.getRawType();
        Parameters parameters = Tool.newInstance(parametersClazz);
        parameters.attachments(inBound.getAttachmentsMap());
        Map<Integer, Type> types = parameters.argumentTypes();
        Object[] arguments = new Object[types.size()];
        types.forEach((idx, kind) -> {
            if (Types.of(Message.class).isAssignableFrom(kind)) {
                if (kind instanceof Class<?>) {
                    try {
                        arguments[idx] = buildMessage((Class<?>) kind, inBound.getArgumentsMap().get(idx));
                    } catch (Exception e) {
                        throw new MeshException(MeshCode.SYSTEM_ERROR, e);
                    }
                }

            }
        });
        parameters.arguments(arguments);
        return (T) parameters;
    }

    @SuppressWarnings("unchecked")
    private <T> T decodeReturns(ByteBuffer buffer, Types<T> type) throws Exception {
        Protobuf4.OutBound outBound = Protobuf4.OutBound.parseFrom(buffer);
        Class<Returns> returnsClazz = (Class<Returns>) type.getRawType();
        Returns returns = Tool.newInstance(returnsClazz);
        returns.setCode(outBound.getCode());
        returns.setMessage(outBound.getMessage());
        Method getContentMethod = returnsClazz.getMethod("getContent");
        Class<?> contentType = (Class<?>) getContentMethod.getGenericReturnType();
        returns.setContent(buildMessage(contentType, outBound.getContent()));
        return (T) returns;
    }

    protected Message buildMessage(Class<?> clazz, ByteString data) throws Exception {
        if (!Message.class.isAssignableFrom(clazz)) {
            return null;
        }
        Method newBuilder = clazz.getMethod("newBuilder");
        Message.Builder builder = (Message.Builder) newBuilder.invoke(null);
        builder.mergeFrom(data);
        return builder.build();
    }
}
