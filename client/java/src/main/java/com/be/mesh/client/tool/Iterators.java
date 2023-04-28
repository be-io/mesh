/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import java.io.Serializable;
import java.util.*;
import java.util.function.Function;
import java.util.function.UnaryOperator;

/**
 * @author coyzeng@gmail.com
 */
public class Iterators<V> implements Iterator<V>, Serializable {

    private final Function<V, List<V>> next;
    private final Deque<V> queue = new ArrayDeque<>();

    public Iterators(V value, UnaryOperator<V> next) {
        this.next = v -> Collections.singletonList(next.apply(value));
        this.queue.push(value);
    }

    public Iterators(V value, Function<V, List<V>> next) {
        this.next = next;
        this.queue.push(value);
    }

    @Override
    public boolean hasNext() {
        return null != queue.peek();
    }

    @Override
    public V next() {
        V value = queue.pop();
        List<V> children = next.apply(value);
        if (null == children || children.isEmpty()) {
            return value;
        }
        children.forEach(queue::push);
        return value;
    }

    @Override
    public void remove() {
        queue.poll();
    }
}
