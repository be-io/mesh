/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.tool.LambdaRn;
import com.be.mesh.client.tool.LambdaSp;
import com.be.mesh.client.tool.SafeAsync;
import lombok.extern.slf4j.Slf4j;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.Optional;
import java.util.concurrent.Callable;

/**
 * Mesh runtime context.
 *
 * <pre>
 *     static {
 *         Set<String> fs = new HashSet<>(1);
 *         fs.add("mtx");
 *         Set<String> ms = new HashSet<>(1);
 *         ms.add("push");
 *         ms.add("pop");
 *         Reflection.registerFieldsToFilter(Mesh.class, fs);
 *         Reflection.registerMethodsToFilter(Mesh.class, ms);
 *     }
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@Slf4j(topic = "mesh-context")
public final class Mesh {

    private Mesh() {
    }

    private static final ThreadLocal<Deque<Context>> mtx = new ThreadLocal<>();

    public static Context context() {
        return Optional.ofNullable(mtx.get()).map(Deque::peek).orElseGet(MpcContext::create);
    }

    public static void release() {
        mtx.remove();
    }

    public static void reset(Context ctx) {
        Deque<Context> contexts = new ArrayDeque<>(3);
        contexts.push(ctx);
        mtx.set(contexts);
    }

    public static boolean isEmpty() {
        return null == mtx.get();
    }

    private static void push(Context context) {
        Deque<Context> contexts = mtx.get();
        if (null == contexts) {
            return;
        }
        contexts.push(context);
    }

    private static void pop() {
        Deque<Context> contexts = mtx.get();
        if (null == contexts || contexts.isEmpty()) {
            return;
        }
        contexts.pop();
    }

    /**
     * Execute with context safety.
     */
    public static void contextSafeUncheck(LambdaRn sp) {
        contextSafeUncheck(() -> {
            sp.execute();
            return true;
        });
    }

    /**
     * Execute with context safety.
     */
    public static <T> T contextSafeUncheck(LambdaSp<T> sp) {
        try {
            return contextSafe(sp);
        } catch (RuntimeException | Error e) {
            throw e;
        } catch (Throwable e) {
            throw new MeshException(e);
        }
    }

    /**
     * Context safe no cause.
     */
    public static void contextSafeCaught(LambdaRn sp) {
        try {
            contextSafe(sp);
        } catch (Throwable e) {
            log.info("", e);
        }
    }

    /**
     * Execute with context safety.
     */
    public static void contextSafe(LambdaRn sp) throws Throwable {
        contextSafe(() -> {
            sp.execute();
            return true;
        });
    }

    /**
     * Execute with context safety.
     */
    public static <T> T contextSafe(LambdaSp<T> sp) throws Throwable {
        boolean contextDisable = Mesh.isEmpty();
        try {
            if (contextDisable) {
                Mesh.reset(MpcContext.create());
            } else {
                Mesh.push(Mesh.context().resume());
            }
            return sp.execute();
        } finally {
            if (contextDisable) {
                Mesh.release();
            } else {
                Mesh.pop();
            }
        }
    }

    /**
     * Construct wrap of context safety.
     */
    public static Runnable threadSafe(Runnable routine) {
        return new SafeAsync<>(routine);
    }

    /**
     * Construct wrap of context safety.
     */
    public static <T> Callable<T> threadSafe(Callable<T> routine) {
        return new SafeAsync<>(routine);
    }

    /**
     * Mesh invoke mpi name attributes.
     */
    public static final Context.Key<String> UNAME = new Context.Key<>("mesh.mpc.uname", Types.of(String.class));

    /**
     * Mesh invocation attributes.
     */
    public static final Context.Key<Invocation> INVOCATION = new Context.Key<>("mesh.invocation", Types.of(Invocation.class));
    /**
     * Mesh mpc remote address.
     */
    public static final Context.Key<String> REMOTE = new Context.Key<>("mesh.mpc.address", Types.of(String.class));
    /**
     * Remote app name.
     */
    public static final Context.Key<String> REMOTE_NAME = new Context.Key<>("mesh.mpc.remote.name", Types.of(String.class));
}
