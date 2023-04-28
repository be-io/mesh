/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.tool.Tool;

import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@SPI(pattern = "*")
public interface Filter {

    /**
     * Provider flag.
     */
    String PROVIDER = "PROVIDER";

    /**
     * Consumer flag.
     */
    String CONSUMER = "CONSUMER";

    /**
     * Invoke the next filter.
     *
     * @param invoker    Service invoker.
     * @param invocation Service invocation
     * @return Invoke result
     * @throws Throwable cause
     */
    Object invoke(Invoker<?> invoker, Invocation invocation) throws Throwable;

    /**
     * Composite the filter spi providers as a invoker.
     *
     * @param invoker Delegate invoker.
     * @param pattern Filter pattern.
     * @return Composite invoker.
     */
    static <T> Invoker<T> composite(Invoker<T> invoker, String pattern) {
        List<Filter> filters = ServiceLoader.load(Filter.class).list().stream().filter(filter -> Optional.ofNullable(filter.getClass().getAnnotation(SPI.class)).map(x -> Tool.equals(x.pattern(), pattern)).orElse(true)).sorted((x, y) -> {
            int xp = Optional.ofNullable(x.getClass().getAnnotation(SPI.class)).map(SPI::priority).orElse(Integer.MIN_VALUE);
            int yp = Optional.ofNullable(y.getClass().getAnnotation(SPI.class)).map(SPI::priority).orElse(Integer.MIN_VALUE);
            return Integer.compare(yp, xp);
        }).collect(Collectors.toList());
        return composite(invoker, filters);
    }

    /**
     * Composite the filter spi providers as an invoker.
     *
     * @param invoker Delegate invoker.
     * @param filters Filter list.
     * @return Composite invoker.
     */
    static <T> Invoker<T> composite(Invoker<T> invoker, List<Filter> filters) {
        Invoker<T> last = invoker;
        for (int i = filters.size() - 1; i >= 0; i--) {
            final Filter filter = filters.get(i);
            final Invoker<?> next = last;
            last = invocation -> {
                SPI spi = filter.getClass().getAnnotation(SPI.class);
                if (null == spi || spi.exclude().length < 1) {
                    return filter.invoke(next, invocation);
                }
                for (String exclude : spi.exclude()) {
                    if (Tool.startWith(invocation.getUrn().getName(), exclude)) {
                        return next.invoke(invocation);
                    }
                }
                return filter.invoke(next, invocation);
            };
        }
        return last;
    }

}
