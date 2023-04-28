/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Principal {
    @index(0)
    public node_id: string = "";
    @index(5)
    public inst_id: string = "";

    public toString(): string {
        return `{"node_id":"${this.node_id}","inst_id":"${this.inst_id}"`;
    }
}

export class Location extends Principal {
    @index(10)
    public ip: string = "";
    @index(15)
    public port: string = "";
    @index(20)
    public host: string = "";
    @index(25)
    public name: string = "";

    public toString(): string {
        return `{"node_id":"${this.node_id}","inst_id":"${this.inst_id}","ip":"${this.ip}","host":"${this.host}","port":"${this.port}","name":"${this.name}"}`;
    }
}