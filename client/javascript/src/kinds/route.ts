/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro";

export class Route {
    @index(0)
    public node_id: string = "";
    @index(5)
    public inst_id: string = "";
    @index(10)
    public name: string = "";
    @index(15)
    public inst_name: string = "";
    @index(20)
    public address: string = "";
    @index(25)
    public describe: string = "";
    @index(30)
    public host_root: string = "";
    @index(35)
    public host_key: string = "";
    @index(40)
    public host_crt: string = "";
    @index(45)
    public guest_root: string = "";
    @index(50)
    public guest_key: string = "";
    @index(55)
    public guest_crt: string = "";
    @index(60)
    public status: number = 0;
    @index(65)
    public version: number = 0;
    @index(70)
    public auth_code: string = "";
    @index(75)
    public expire_at: number = 0;
    @index(80)
    public extra: string = "";
    @index(85)
    public create_at: string = "";
    @index(90)
    public create_by: string = "";
    @index(95)
    public update_at: string = "";
    @index(100)
    public update_by: string = "";
    @index(105)
    public group: string = "";
    @index(110)
    public upstream: string = "";
    @index(115)
    public downstream: string = "";
    @index(120)
    public static_ip: string = "";
    @index(125)
    public proxy: string = "";
    @index(130)
    public concurrency: number = 0;
}