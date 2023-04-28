/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.tool.Tool;
import lombok.EqualsAndHashCode;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.UndeclaredThrowableException;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@EqualsAndHashCode(callSuper = true)
public class GenericInvokeHandler extends ReferenceInvokeHandler {

    private final Context.Key<GenericExecution> executionKey = new Context.Key<>("mesh.generic.execution", Types.of(GenericExecution.class));

    public GenericInvokeHandler() {
        super(Collections.emptyMap());
    }

    public Object invoke(String urn, Map<String, Object> arguments) throws Throwable {
        return Mesh.contextSafe(() -> {
            try {
                URN unx = URN.from(urn);
                GenericExecution execution = new GenericExecution(unx);
                Mesh.context().setAttribute(executionKey, execution);
                Mesh.context().rewriteUrn(rewriteURN(execution));
                Mesh.context().setAttribute(Mesh.REMOTE, rewriteAddress(urn));

                GenericParameters parameters = new GenericParameters();
                parameters.attachments(new HashMap<>());
                parameters.putAll(arguments);

                ServiceInvocation invocation = new ServiceInvocation();
                invocation.setProxy(this);
                invocation.setInspector(execution.inspect());
                invocation.setParameters(parameters);
                invocation.setExecution(execution);
                invocation.setUrn(URN.from(urn));

                Mesh.context().setAttribute(Mesh.INVOCATION, invocation);

                return invoker.get().invoke(invocation);
            } catch (InvocationTargetException | UndeclaredThrowableException e) {
                throw Tool.destructor(e.getCause());
            } catch (Throwable e) {
                throw Tool.destructor(e);
            }
        });
    }

    @Override
    protected Execution<Reference> referExecution(Inspector inspector) {
        return Mesh.context().getAttribute(executionKey);
    }
}
