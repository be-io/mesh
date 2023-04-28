/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import lombok.extern.slf4j.Slf4j;

import java.util.Optional;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Supplier;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public final class Once<T> {

    private Supplier<T> producer;
    private volatile T value;

    public T get(Supplier<T> producer) {
        this.producer = producer;
        return get();
    }

    public boolean isPresent() {
        return null != get();
    }

    public void ifPresent(Consumer<T> consumer) {
        if (isPresent()) {
            consumer.accept(get());
        }
    }

    public <V> Optional<V> map(Function<T, V> consumer) {
        if (isPresent()) {
            return Optional.ofNullable(consumer.apply(get()));
        }
        return Optional.empty();
    }

    public boolean isPresentWithoutGet() {
        return null != this.value;
    }

    public T get() {
        if (null != value) {
            return value;
        }
        if (null != producer) {
            this.value = producer.get();
        }
        return this.value;
    }

    public void put(T value) {
        this.value = value;
    }

    public void release() {
        this.value = null;
    }

    public static <T> Once<T> with(Supplier<T> producer) {
        Once<T> once = new Once<>();
        once.producer = producer;
        return once;
    }

    public static <T> Once<T> of(T value) {
        return with(() -> value);
    }

    public <V> Once<V> any(Function<T, V> m) {
        if (isPresent()) {
            return Once.of(m.apply(get()));
        }
        return Once.empty();
    }

    public static <T> Once<T> empty() {
        return new Once<>();
    }

}
