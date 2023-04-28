/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export const environVersion = "1.0.0";

export class Environ {

    @index(0)
    public version: string = "";
    @index(5)
    public node_id: string = "";
    @index(10)
    public inst_id: string = "";
    @index(15)
    public inst_name: string = "";
    @index(20)
    public root_crt: string = "";
    @index(25)
    public root_key: string = "";
    @index(30)
    public node_crt: string = "";

}

export class Lattice {
    @index(0)
    public zone: string = "";
    @index(5)
    public cluster: string = "";
    @index(10)
    public cell: string = "";
    @index(15)
    public group: string = "";
    @index(20)
    public address: string = "";
}