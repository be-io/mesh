/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {dict, Dict, list, List} from "@/macro/dsa";


export type  Type<T> = Function & { prototype: T }
export type  Class<T> = Type<T> & { new(...args: any[]): T }

class Species {

    public static readonly FUNCTION = 1;
    public static readonly ABSTRACT = 2;
    public readonly macro: Type<any>;
    public readonly name: string;
    public readonly kind: Type<any>;
    public readonly traits: Type<any>[];
    public readonly metadata: any;
    public readonly modifier: number;

    constructor(macro: Type<any>, name: string, kind: Type<any>, traits: Type<any>[], metadata: any, modifier: number) {
        this.macro = macro;
        this.name = name;
        this.kind = kind;
        this.traits = traits;
        this.metadata = metadata;
        this.modifier = modifier;
    }

    isfunction(): boolean {
        return (this.modifier & Species.FUNCTION) == Species.FUNCTION;
    }

    isabstract(): boolean {
        return (this.modifier & Species.ABSTRACT) == Species.ABSTRACT;
    }

}

class ark {

    // key/value
    private tinder: Dict<string, List<Species>> = new dict();
    // parent/target/macro/metadata
    private annotations: Dict<any, Dict<any, Dict<any, any>>> = new dict();

    private ak(macro: Type<any>, name: string, kind: Type<any>): string {
        return `${macro.name}-${name}-${kind.name}`
    }

    /**
     * Annotate handler.
     * @param target
     * @param parent
     * @param macro
     * @param metadata
     */
    annotate(macro: any, target: any, parent: any, metadata: any): void {
        this.annotations.computeIfy(parent, k => {
            return new dict<any, Dict<any, any>>()
        }).computeIfy(target, k => {
            return new dict<any, any>()
        }).set(macro, metadata);
    }

    /**
     * Inspect the metadata of target object.
     * @param parent
     * @param target
     * @param macro
     */
    metadata<T>(macro: Type<T>, parent: any, target: any): T | undefined {
        return this.annotations.get(parent)?.get(target)?.get(macro);
    }

    /**
     * Register a annotated object as a metadata.
     * @param macro annotated decorator
     * @param name object name
     * @param kind object type
     * @param metadata decorated object
     */
    register(macro: Type<any>, name: string, kind: Type<any>, metadata: any): void {
        const traits: Type<any>[] = [];
        let sc = kind.prototype;
        while (sc && sc.constructor) {
            if (kind != sc.constructor && this.metadata(macro, sc.constructor, sc.constructor)) {
                traits.push(sc.constructor);
            }
            sc = Object.getPrototypeOf(sc);
        }
        const isfunction = !kind.constructor ? Species.FUNCTION : 0
        const isabstract = traits.length < 1 ? Species.ABSTRACT : 0;
        if (traits.length < 1) {
            traits.push(kind);
        }
        for (let trait of traits) {
            const key = this.ak(macro, name, trait);
            this.tinder.get(key)?.forEach((v) => {
                if (v.kind == kind) {
                    console.log(`Object of ${trait.name} named ${name} has been register already.`)
                }
            })
            const species = new Species(macro, name, kind, traits, metadata, isfunction | isabstract);
            const pvs = this.tinder.computeIfy(key, k => new list());
            pvs.push(species);
            pvs.sort((p, n) => {
                    if (p.isabstract() || p.isfunction()) {
                        return 1;
                    }
                    if (n.isabstract() || p.isfunction()) {
                        return -1;
                    }
                    if (p.metadata.priority && n.metadata.priority) {
                        return p.metadata.priority - n.metadata.priority;
                    }
                    return 0;
                }
            )
        }
    }

    /**
     * Unregister the class metadata.
     * @param macro
     * @param kind
     */
    unregister(macro: Type<any>, kind: Type<any>): void {
        const matches: string[] = [];
        this.tinder.forEach((vs, k) => vs.forEach(v => {
            if (v.macro == macro && v.kind == kind) {
                matches.push(k);
            }
        }))
        for (let match of matches) {
            this.tinder.delete(match);
        }
    }

    /**
     * Export all classes.
     * @param macro
     */
    export(macro: any): Dict<Type<any>, Dict<string, List<Class<any>>>> {
        const matches = new dict<Type<any>, Dict<string, List<Class<any>>>>();
        for (let species of this.tinder.values()) {
            for (let spec of species) {
                if (spec.isfunction() || spec.macro != macro) {
                    continue
                }
                // @ts-ignore
                const kind: Class<any> = spec.kind;
                for (let trait of spec.traits) {
                    matches.computeIfy(trait, k => new dict()).computeIfy(spec.name, k => new list()).push(kind);
                }
            }
        }
        return matches
    }

    /**
     * Inspect the subclass of kind annotated with macro.
     * @param macro
     * @param kind
     */
    providers<T>(macro: any, kind: Type<T>): List<Class<T>> {
        // @ts-ignore
        return this.export(macro).computeIfy(kind, k => new dict()).map((k, v) => v).flatMap(v => v).filter(v => v != kind);
    }

    /**
     * Inspect the subclass of kind annotated with macro for name.
     * @param macro
     * @param kind
     * @param name
     */
    provider<T>(macro: Type<any>, kind: Type<T>, name: string): List<Class<T>> {
        const key = this.ak(macro, name, kind);
        const pvs: List<Class<T>> = new list();
        this.tinder.computeIfy(key, k => new list()).forEach((pv) => {
            if (!pv.isfunction() && !pv.isabstract()) {
                // @ts-ignore
                pvs.push(pv.kind);
            }
        });
        return pvs;
    }

    inspect<T>(kind: Type<T>): Dict<string, any> {
        const dic: Dict<string, any> = new dict();
        Object.entries(Object.getOwnPropertyDescriptors(kind.prototype)).filter(x => {
            return x[0] !== 'constructor' && typeof x[1].value === 'function';
        }).forEach(x => {
            dic.set(x[0], x[1].value);
        })
        return dic;
    }

}

const Ark = new ark();

export {
    Ark
}