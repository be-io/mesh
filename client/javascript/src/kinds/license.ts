/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class License {

    @index(0)
    public version: string = "";
    @index(5)
    public level: number = 0;
    @index(10)
    public name: string = "";
    @index(15)
    public create_by: string = "";
    @index(20)
    public create_at: number = 0;
    @index(25)
    public active_at: number = 0;
    @index(30)
    public factors: string[] = [];
    @index(35)
    public signature: string = "";
    @index(40)
    public node_id: string = "";
    @index(45)
    public inst_id: string = "";
    @index(50)
    public server: string = "";
    @index(55)
    public crt: string = "";
    @index(60)
    public group: string[] = [];
    @index(65)
    public replicas: number = 0;
    @index(70)
    public max_cooperators: number = 0;
    @index(75)
    public max_tenants: number = 0;
    @index(80)
    public max_users: number = 0;
    @index(85)
    public max_mills: number = 0;
    @index(90)
    public white_urns: string[] = [];
    @index(95)
    public black_urns: string[] = [];
    @index(100)
    public super_urns: string[] = [];

}
