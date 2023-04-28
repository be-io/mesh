/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.prsim.Codeable;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.struct.Outbound;
import com.be.mesh.client.struct.Service;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;

import java.nio.ByteBuffer;
import java.util.Optional;

import static com.be.mesh.client.mpc.Transporter.PROVIDER;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI(PROVIDER)
public class ProviderTransporter implements Transporter {

    @Override
    public ByteBuffer transport(String uname, ByteBuffer buffer) throws Throwable {
        return Mesh.contextSafe(() -> {
            try {
                Mesh.context().rewriteUrn(uname);
                URN urn = URN.from(uname);
                Eden eden = ServiceLoader.load(Eden.class).getDefault();
                Execution<Service> execution = eden.infer(uname);
                if (null == execution) {
                    Outbound outbound = new Outbound();
                    outbound.setCode(MeshCode.NO_SERVICE.getCode());
                    outbound.setMessage(String.format("No mpi named %s.", urn.getName()));
                    return ServiceLoader.load(Codec.class).getDefault().encode(outbound);
                }
                String name = Optional.of(urn).map(URN::getFlag).map(URNFlag::getCodec).map(MeshFlag::ofCodec).orElse(MeshFlag.JSON).getName();
                Codec codec = ServiceLoader.load(Codec.class).get(Tool.anyone(execution.schema().getCodec(), name));
                return service(urn, codec, execution, buffer);
            } catch (Throwable e) {
                log.error(String.format("[%s#%s] Invoke service %s with error.", Mesh.context().getTraceId(), Mesh.context().getSpanId(), Mesh.context().getUrn()), e);
                Codec codec = ServiceLoader.load(Codec.class).getDefault();
                Outbound outbound = new Outbound();
                outbound.setCode(MeshCode.SYSTEM_ERROR.getCode());
                outbound.setMessage(e.getMessage());
                outbound.setCause(Cause.of(e));
                return codec.encode(outbound);
            }
        });
    }

    /**
     * Construct the {@link Invocation} to invoke service.
     *
     * @param urn    Uniform resource domain name.
     * @param codec  Payload codec.
     * @param buffer Input stream.
     */
    private ByteBuffer service(URN urn, Codec codec, Execution<Service> execution, ByteBuffer buffer) {
        try {
            if (null == execution) {
                Outbound outbound = new Outbound();
                outbound.setCode(MeshCode.NOT_FOUND.getCode());
                outbound.setMessage(String.format("No mpi named %s.", urn.toString()));
                return codec.encode(outbound);
            }
            Parameters parameters = codec.decode(buffer, execution.intype());

            ServiceInvocation invocation = new ServiceInvocation();
            invocation.setProxy(execution);
            invocation.setInspector(execution.inspect());
            invocation.setParameters(parameters);
            invocation.setExecution(execution);
            invocation.setUrn(urn);

            Object result = execution.invoke(invocation);
            Returns returns = execution.reflect();
            returns.setCode(MeshCode.SUCCESS.getCode());
            returns.setContent(result);
            return codec.encode(returns);
        } catch (Throwable e) {
            log.error(String.format("Invoke service %s with error.", urn), e);
            Outbound outbound = new Outbound();
            if (e instanceof Codeable) {
                outbound.setCode(((Codeable) e).getCode());
            } else {
                outbound.setCode(MeshCode.SYSTEM_ERROR.getCode());
            }
            outbound.setMessage(e.getMessage());
            outbound.setCause(Cause.of(e));
            return codec.encode(outbound);
        }
    }
}
