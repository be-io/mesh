/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Versions {

    /**
     * Platform version.
     */
    @index(0)
    public version: string = "";
    /**
     * Network modules info.
     */
    @index(5)
    public infos!: Map<string, string>;

    /**
     * Is the version less than the appointed version.
     *
     * @param v compared version
     * @return true less than
     */
    public compare(v: string): number {
        if (!this.version) {
            return -1;
        }
        if (!v) {
            return 1;
        }
        const rvs = this.version.split(".");
        const vs = v.split(".");
        for (let index = 0; index < rvs.length; index++) {
            if (vs.length <= index) {
                return 0;
            }
            if (vs[index] == "*") {
                continue
            }
            const c = rvs[index].localeCompare(vs[index]);
            if (0 != c) {
                return c;
            }
        }
        return 0;
    }

    public compareTo(ver: Versions): number {
        if (!ver) {
            return 1;
        }
        return this.compare(ver.version);
    }

    public anyVersions(sets: string): Versions {
        const anyInfos = this.infos || new Map<string, string>();
        const vk = `${sets}.version`;
        const ak = `${sets}.arch`;
        const ok = `${sets}.os`;
        const ck = `${sets}.commit_id`;
        const myInfos = new Map<string, string>();
        myInfos.set(vk, anyInfos.get(vk) || "");
        myInfos.set(ak, anyInfos.get(ak) || "");
        myInfos.set(ok, anyInfos.get(ok) || "");
        myInfos.set(ck, anyInfos.get(ck) || "");
        const versions = new Versions();
        versions.version = anyInfos.get(vk) || "";
        versions.infos = myInfos;
        return versions;
    }
}