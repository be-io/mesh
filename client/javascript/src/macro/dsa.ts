/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

export interface Dict<K, V> extends Map<K, V> {

    computeIfy(key: K, value: (key: K) => V): V;

    groupBy<GK>(group: (key: K, value: V) => GK): Dict<GK, V[]>;

    map<T>(m: (key: K, value: V) => T): List<T>;

    groupKV<KK, VV>(key: (key: K) => KK, value: (key: K, value: V) => VV): Dict<KK, VV>;
}

export class dict<K, V> extends Map<K, V> implements Dict<K, V> {

    computeIfy(key: K, value: (key: K) => V): V {
        let v = this.get(key);
        if (!v) {
            v = value(key);
            this.set(key, v);
        }
        return v;
    }

    groupBy<GK>(group: (key: K, value: V) => GK): Dict<GK, V[]> {
        const groups: Dict<GK, V[]> = new dict();
        this.forEach((v, k) => groups.computeIfy(group(k, v), k => []).push(v));
        return groups;
    }

    map<T>(m: (key: K, value: V) => T): List<T> {
        const vs: List<T> = new list();
        this.forEach((v, k) => vs.push(m(k, v)));
        return vs;
    }

    groupKV<KK, VV>(key: (key: K) => KK, value: (key: K, value: V) => VV): Dict<KK, VV> {
        const groups: Dict<KK, VV> = new dict();
        this.forEach((v, k) => groups.set(key(k), value(k, v)));
        return groups;
    }

}

export interface List<T> extends Array<T> {

    groupBy<K>(group: (value: T) => K): Dict<K, List<T>>;

    groupKV<K, V>(group: (value: T) => K, value: (key: K, value: T[]) => V): Dict<K, V>;

}

export class list<T> extends Array implements List<T> {

    groupBy<K>(group: (value: T) => K): Dict<K, List<T>> {
        const groups: Dict<K, List<T>> = new dict();
        this.forEach((v) => groups.computeIfy(group(v), k => new list()).push(v));
        return groups;
    }

    groupKV<K, V>(group: (value: T) => K, value: (key: K, value: T[]) => V): Dict<K, V> {
        const groups: Dict<K, T[]> = new dict();
        this.forEach((v) => groups.computeIfy(group(v), k => []).push(v));
        const kvs: Dict<K, V> = new dict();
        groups.forEach((v, k) => kvs.set(k, value(k, v)));
        return kvs;
    }

}