/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.struct.Outbound;
import com.be.mesh.client.tool.Tool;
import io.grpc.*;
import io.grpc.stub.ServerCalls;
import io.grpc.stub.StreamObserver;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.nio.ByteBuffer;

/**
 * @author coyzeng@gmail.com
 */
public class GrpcBindableService implements BindableService {

    @Override
    public ServerServiceDefinition bindService() {
        MethodDescriptor.Marshaller<InputStream> im = new GrpcMarshaller();
        MethodDescriptor.Marshaller<InputStream> om = new GrpcMarshaller();
        MethodDescriptor.Builder<InputStream, InputStream> descriptor = MethodDescriptor.newBuilder(im, om);
        descriptor.setType(MethodDescriptor.MethodType.UNARY);
        descriptor.setFullMethodName(GrpcUtils.MESH_INVOKE_METHOD);
        ServerCallHandler<InputStream, InputStream> handler = ServerCalls.asyncUnaryCall(new UnaryStreamHandler());
        ServerMethodDefinition<InputStream, InputStream> definition = ServerMethodDefinition.create(descriptor.build(), handler);
        return ServerServiceDefinition.builder(GrpcUtils.MESH_SERVICE_NAME).addMethod(definition).build();
    }

    @Slf4j
    @AllArgsConstructor
    private static class AsyncInputObserver implements StreamObserver<InputStream> {

        private final StreamObserver<InputStream> observer;

        @Override
        public void onNext(InputStream input) {
            Mesh.contextSafeCaught(() -> {
                try (InputStream ref = input) {
                    GrpcContextKey.setContext();
                    Transporter transporter = ServiceLoader.load(Transporter.class).get(Transporter.PROVIDER);
                    ByteBuffer buffer = ByteBuffer.wrap(Tool.readBytes(ref));
                    ByteBuffer outbound = transporter.transport(GrpcContextKey.URN.get(), buffer);
                    this.observer.onNext(new ByteArrayInputStream(outbound.array()));
                    this.observer.onCompleted();
                } catch (Throwable e) {
                    log.error(String.format("[%s#%s] Invoke service %s with error.", GrpcContextKey.CTX_TRACE_ID.get(), GrpcContextKey.CTX_SPAN_ID.get(), GrpcContextKey.URN.get()), e);
                    Codec codec = ServiceLoader.load(Codec.class).getDefault();
                    Outbound outbound = new Outbound();
                    outbound.setCode(MeshCode.SYSTEM_ERROR.getCode());
                    outbound.setMessage(e.getMessage());
                    outbound.setCause(Cause.of(e));
                    this.observer.onNext(new ByteArrayInputStream(codec.encode(outbound).array()));
                    this.observer.onCompleted();
                }
            });
        }

        @Override
        public void onError(Throwable e) {
            log.error("Listen inbound with error.", e);
            this.observer.onError(e);
        }

        @Override
        public void onCompleted() {
            this.observer.onCompleted();
        }

    }

    private static final class BidiStreamingHandler implements ServerCalls.BidiStreamingMethod<InputStream, InputStream> {

        @Override
        public StreamObserver<InputStream> invoke(StreamObserver<InputStream> observer) {
            return new AsyncInputObserver(observer);
        }
    }

    private static final class UnaryStreamHandler implements ServerCalls.UnaryMethod<InputStream, InputStream> {

        @Override
        public void invoke(InputStream inputStream, StreamObserver<InputStream> observer) {
            new AsyncInputObserver(observer).onNext(inputStream);
        }
    }

}
