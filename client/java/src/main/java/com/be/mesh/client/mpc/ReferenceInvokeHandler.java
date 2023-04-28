/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.cause.*;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.struct.Principal;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.tool.Once;
import com.be.mesh.client.tool.Tool;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.lang.reflect.InvocationHandler;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.UndeclaredThrowableException;
import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@Data
@AllArgsConstructor
public class ReferenceInvokeHandler implements InvocationHandler, Invoker<Object> {

    private final Map<Method, Inspector> inspectors;
    private static final Once<DefaultHandler> dh = Once.with(DefaultInvocationHandler::new);
    protected final Once<Invoker<?>> invoker = Once.with(() -> Filter.composite(this, Filter.CONSUMER));

    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        return Mesh.contextSafe(() -> {
            try {
                Class<?> c = method.getDeclaringClass();
                if (c == Object.class) {
                    return method.invoke(this, args);
                }
                if (method.isDefault()) {
                    return dh.get().invoke(proxy, method, args);
                }
                Execution<Reference> execution = this.referExecution(inspectors.get(method));
                String urn = rewriteURN(execution);
                Mesh.context().rewriteUrn(urn);
                Mesh.context().setAttribute(Mesh.REMOTE, rewriteAddress(urn));

                Parameters parameters = execution.inflect();
                parameters.arguments(args);
                parameters.attachments(new HashMap<>());

                ServiceInvocation invocation = new ServiceInvocation();
                invocation.setProxy(this);
                invocation.setInspector(inspectors.get(method));
                invocation.setParameters(parameters);
                invocation.setExecution(execution);
                invocation.setUrn(URN.from(urn));

                Mesh.context().setAttribute(Mesh.INVOCATION, invocation);

                return invoker.get().invoke(invocation);
            } catch (InvocationTargetException | UndeclaredThrowableException e) {
                throw Tool.destructor(e.getCause(), method.getExceptionTypes());
            } catch (Throwable e) {
                throw Tool.destructor(e, method.getExceptionTypes());
            }
        });
    }

    @Override
    public Object invoke(Invocation invocation) throws Throwable {
        Execution<Reference> execution = this.referExecution(invocation.getInspector());
        Consumer consumer = ServiceLoader.load(Consumer.class).getDefault();
        String address = Tool.anyone(Mesh.context().getAttribute(Mesh.REMOTE), Tool.MESH_ADDRESS.get().any());
        String name = Optional.ofNullable(execution.schema().getCodec()).orElse(MeshFlag.JSON.getName());
        Codec codec = ServiceLoader.load(Codec.class).get(name);
        ByteBuffer buffer = codec.encode(invocation.getParameters());
        InputStream input = new ByteArrayInputStream(buffer.array());
        Future<InputStream> future = consumer.consume(address, Mesh.context().getUrn(), execution, input);
        if (Future.class.isAssignableFrom(invocation.getInspector().getReturnType())) {
            return CompletableFuture.supplyAsync(() -> {
                try {
                    return deserialize(execution, codec, future);
                } catch (RuntimeException | Error e) {
                    throw e;
                } catch (Throwable e) {
                    throw new MeshException(e);
                }
            });
        }
        return deserialize(execution, codec, future);
    }

    private Object deserialize(Execution<Reference> execution, Codec codec, Future<InputStream> future) throws Throwable {
        try (InputStream output = future.get(execution.schema().getTimeout(), TimeUnit.MILLISECONDS)) {
            Returns returns = codec.decode(ByteBuffer.wrap(Tool.readBytes(output)), execution.retype());
            if (null != returns.getCause()) {
                throw Cause.of(returns.getCode(), returns.getMessage(), returns.getCause());
            }
            if (Tool.equals(MeshCode.NOT_FOUND.getCode(), returns.getCode())) {
                throw new NotFoundException(returns.getMessage());
            }
            if (Tool.equals(MeshCode.NO_SERVICE.getCode(), returns.getCode())) {
                throw new NoServiceException(returns.getMessage());
            }
            if (Tool.equals(MeshCode.NO_PROVIDER.getCode(), returns.getCode())) {
                throw new NoProviderException(returns.getMessage());
            }
            if (!Tool.equals(MeshCode.SUCCESS.getCode(), returns.getCode())) {
                throw new MeshException(returns.getMessage());
            }
            if (execution.inspect().getReturnType() == void.class || execution.inspect().getReturnType() == Void.class) {
                return null;
            }
            return returns.getContent();
        } catch (ExecutionException e) {
            throw e.getCause();
        } catch (java.util.concurrent.TimeoutException e) {
            throw new TimeoutException("Invoke %s timeout with %dms", execution.schema().getUrn(), execution.schema().getTimeout());
        }
    }

    protected Execution<Reference> referExecution(Inspector inspector) {
        Eden eden = ServiceLoader.load(Eden.class).getDefault();
        Execution<Reference> execution = eden.refer(inspector.getAnnotation(MPI.class), inspector.getType(), inspector);
        if (null != execution) {
            return execution;
        }
        throw new CompatibleException("Method %s cant be compatible", inspector.getName());
    }

    /**
     * Rewrite the urn by execution context.
     */
    protected String rewriteURN(Execution<Reference> execution) {
        String nodeId = Optional.ofNullable(Mesh.context().getPrincipals().peek()).map(Principal::getNodeId).orElse("");
        String instId = Optional.ofNullable(Mesh.context().getPrincipals().peek()).map(Principal::getInstId).orElse("");
        if (Tool.optional(nodeId) && Tool.optional(instId) && !Mesh.REMOTE_NAME.isPresent() && !Mesh.UNAME.isPresent()) {
            return execution.schema().getUrn();
        }
        URN urn = URN.from(execution.schema().getUrn());
        if (Tool.required(instId)) {
            urn.setNodeId(instId);
        }
        if (Tool.required(nodeId)) {
            urn.setNodeId(nodeId);
        }
        Mesh.UNAME.ifPresent(uname -> urn.setName(urn.getName().replace("${mesh.uname}", uname)));
        Mesh.REMOTE_NAME.ifPresent(name -> urn.setName(urn.getName().replace("${mesh.name}", name)));
        return urn.toString();
    }

    /**
     * Select address if target is direct.
     */
    protected String rewriteAddress(String uns) {
        if (Tool.required(Mesh.context().getAttribute(Mesh.REMOTE))) {
            return Mesh.context().getAttribute(Mesh.REMOTE);
        }
        URN urn = URN.from(uns);
        if (Tool.startWith(urn.getName(), "mesh.")) {
            return Tool.MESH_ADDRESS.get().any();
        }
        if (Tool.required(Tool.MESH_DIRECT.get())) {
            String[] names = Tool.split(Tool.MESH_DIRECT.get(), ",");
            for (String name : names) {
                String[] pair = Tool.split(name, "=");
                if (isDirect(urn, pair)) {
                    return pair[1];
                }
            }
        }
        String address = urn.getFlag().getAddress().replace(".", "");
        if (Tool.isNumeric(address) && Long.parseLong(address) > 0) {
            return String.format("%s:%s", urn.getFlag().getAddress(), urn.getFlag().getPort());
        }
        return Tool.MESH_ADDRESS.get().any();
    }

    private boolean isDirect(URN urn, String[] pair) {
        if (pair.length < 2 || Tool.optional(pair[1])) {
            return false;
        }
        if (!Tool.contains(pair[0], "@")) {
            return Tool.startWith(urn.getName(), pair[0]);
        }
        String[] nn = pair[0].split("@");
        if (nn.length < 2 || Tool.optional(nn[1])) {
            return false;
        }
        return Tool.equals(Tool.toLowerCase(urn.getNodeId()), Tool.toLowerCase(nn[1])) && Tool.startWith(urn.getName(), nn[0]);
    }

    public boolean isGT16() {
        String[] vpr = System.getProperty("java.version").split("\\.");
        int discard = Integer.parseInt(vpr[0]);
        if (discard == 1) {
            return 16 <= Integer.parseInt(vpr[1]);
        } else {
            return 16 <= discard;
        }
    }

}
