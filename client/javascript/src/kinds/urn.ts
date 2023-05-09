/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import Tool from "@/tool/tool";

export class URN {
    public static CN = "trustbe";
    public static MESH_DOMAIN = this.CN + ".cn";
    public static LocalNodeId = "LX0000000000000";
    public static LocalInstId = "JG0000000000000000";
    public domain: string = "";
    public nodeId: string = "";
    public flag: URNFlag = new URNFlag();
    public name: string = "";

    public static from(urn: string): URN {
        const name = new URN();
        if (!urn) {
            console.log(`Unresolved urn ${urn}`)
            return name;
        }
        const names = name.asArray(urn);
        if (names.length < 5) {
            console.log(`Unresolved urn ${urn}`)
            return name;
        }
        names.reverse();
        name.domain = `${names.splice(1, 1).join("")}.${names.splice(0, 1).join("")}`;
        name.nodeId = names.splice(0, 1).join("");
        name.flag = URNFlag.from(names.splice(0, 1).join(""));
        name.name = names.join(".");
        return name;
    }

    public static urn(name: string): string {
        const urn = new URN();
        urn.nodeId = this.LocalNodeId;
        urn.flag = URNFlag.from("0001000000000000000000000000000000000");
        urn.name = name;
        urn.domain = this.MESH_DOMAIN;
        return urn.toString();
    }

    toString(): string {
        const urn = [];
        const names = this.asArray(this.name);
        names.reverse();
        names.filter(x => x).forEach(x => urn.push(x));
        urn.push(this.flag.toString());
        urn.push(this.nodeId);
        urn.push(Tool.anyone(this.domain, URN.MESH_DOMAIN));
        return urn.join(".");
    }

    private asArray(text: string): string[] {
        const pairs = text.split(".");
        const names = [];
        const buffer = [];
        for (let pair of pairs) {
            if (pair.startsWith("${")) {
                buffer.push(pair);
                buffer.push('.');
                continue
            }
            if (pair.endsWith("}")) {
                buffer.push(pair);
                names.push(buffer.join(""))
                buffer.splice(0, buffer.length);
                continue
            }
            names.push(pair)
        }
        return names;
    }

    matchName(name: string): boolean {
        if (!name) {
            return false;
        }
        return name.endsWith(".*") && name.startsWith(name.replace(".*", ""));
    }

    resetFlag(proto: string, codec: string): URN {
        this.flag.proto = proto;
        this.flag.codec = codec;
        return this;
    }
}

export class URNFlag {

    public v: string = "00";
    public proto: string = "00";
    public codec: string = "00";
    public version: string = "000000";
    public zone: string = "00";
    public cluster: string = "00";
    public cell: string = "00";
    public group: string = "00";
    public address: string = "000000000000";
    public port: string = "00000";

    public static from(value: string): URNFlag {
        if (!value) {
            return new URNFlag();
        }
        const flag = new URNFlag();
        flag.v = Tool.substring(value, 0, 2);
        flag.proto = Tool.substring(value, 2, 2);
        flag.codec = Tool.substring(value, 4, 2);
        flag.version = this.reduce(Tool.substring(value, 6, 6), 2);
        flag.zone = Tool.substring(value, 12, 2);
        flag.cluster = Tool.substring(value, 14, 2);
        flag.cell = Tool.substring(value, 16, 2);
        flag.group = Tool.substring(value, 18, 2);
        flag.address = this.reduce(Tool.substring(value, 20, 12), 3);
        flag.port = this.reduce(Tool.substring(value, 32, 5), 5);
        return flag;
    }

    toString(): string {
        const v = this.padding(Tool.anyone(this.v, ""), 2);
        const proto = this.padding(Tool.anyone(this.proto, ""), 2);
        const codec = this.padding(Tool.anyone(this.codec, ""), 2);
        const version = this.paddingChain(Tool.anyone(this.version, ""), 2, 3);
        const zone = this.padding(Tool.anyone(this.zone, ""), 2);
        const cluster = this.padding(Tool.anyone(this.cluster, ""), 2);
        const cell = this.padding(Tool.anyone(this.cell, ""), 2);
        const group = this.padding(Tool.anyone(this.group, ""), 2);
        const address = this.paddingChain(Tool.anyone(this.address, ""), 3, 4);
        const port = this.padding(Tool.anyone(this.port, ""), 5);
        return `${v}${proto}${codec}${version}${zone}${cluster}${cell}${group}${address}${port}`;
    }

    private padding(v: string, length: number): string {
        const value = v.replace(".", "");
        if (value.length == length) {
            return value;
        }
        if (value.length < length) {
            return Tool.repeat('0', length - v.length) + value;
        }
        return v.substring(0, length);
    }

    private paddingChain(v: string, length: number, size: number): string {
        const array = v.split("[.*]");
        const frags = array.map(x => parseInt(x) ? x : "");
        if (frags.length == size) {
            return frags.map(x => this.padding(x, length)).join("")
        }
        if (frags.length < size) {
            for (let index = 0; index < size - array.length; index++) {
                frags.push("");
            }
            return frags.map(x => this.padding(x, length)).join("")
        }
        return frags.slice(0, size).map(x => this.padding(x, length)).join("")
    }

    private static reduce(v: string, length: number): string {
        if (!v) {
            return v;
        }
        const bu: string[] = [];
        for (let index = 0; index < v.length; index += length) {
            let hasNoneZero: boolean = true;
            for (let offset = index; offset < index + length; offset++) {
                if (v.charAt(offset) != '0') {
                    hasNoneZero = false;
                }
                if (!hasNoneZero) {
                    bu.push(v.charAt(offset));
                }
            }
            if (hasNoneZero) {
                bu.push('0');
            }
            bu.push('.');
        }
        return bu.slice(0, bu.length - 1).join("");
    }
}