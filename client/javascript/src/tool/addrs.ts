/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {list, List} from "@/macro";

export class Addrs {

    private readonly servers: List<StatefulServer> = new list();
    private readonly availableAddrs: List<string> = new list();

    constructor(addrs: string) {
        for (let server of (addrs || "").split(",")) {
            if (server && server.trim().length > 0 && !this.availableAddrs.includes(server)) {
                this.servers.push(new StatefulServer(true, server));
                this.availableAddrs.push(server);
            }
        }
    }

    any(): string {
        const addrs = this.availableAddrs;
        if (!addrs || addrs.length < 1) {
            return "";
        }
        return addrs[Math.floor(Math.random() * addrs.length)]
    }
}

export class StatefulServer {
    private available: boolean;
    private address: string;

    constructor(available: boolean, address: string) {
        this.available = available;
        this.address = address;
    }
}