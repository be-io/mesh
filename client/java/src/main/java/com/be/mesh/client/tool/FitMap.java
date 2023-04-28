/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import java.util.Collection;
import java.util.Map;
import java.util.Set;

/**
 * @author coyzeng@gmail.com
 */
public class FitMap<K, V> implements Map<K, V> {

    private final Map<K, V> dict;

    public FitMap(Map<K, V> fixed) {
        this.dict = fixed;
    }

    @Override
    public int size() {
        return this.dict.size();
    }

    @Override
    public boolean isEmpty() {
        return this.dict.isEmpty();
    }

    @Override
    public boolean containsKey(Object key) {
        return this.dict.containsKey(key);
    }

    @Override
    public boolean containsValue(Object value) {
        return this.dict.containsValue(value);
    }

    @Override
    public V get(Object key) {
        return this.dict.get(key);
    }

    @Override
    public V put(K key, V value) {
        return this.dict.put(key, value);
    }

    @Override
    public V remove(Object key) {
        return this.dict.remove(key);
    }

    @Override
    public void putAll(Map<? extends K, ? extends V> m) {
        this.dict.putAll(m);
    }

    @Override
    public void clear() {
        this.dict.clear();
    }

    @Override
    public Set<K> keySet() {
        return this.dict.keySet();
    }

    @Override
    public Collection<V> values() {
        return this.dict.values();
    }

    @Override
    public Set<Entry<K, V>> entrySet() {
        return this.dict.entrySet();
    }
}
