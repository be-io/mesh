/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Keys {

    @index(0)
    public cno: string = ""
    @index(5)
    public pno: string = "";
    @index(10)
    public kno: string = "";
    @index(15)
    public kind: string = "";
    @index(20)
    public csr: string = "";
    @index(25)
    public key: string = "";
    @index(30)
    public status: number = 0;
    @index(35)
    public version: string = "";
}

export class KeyCsr {

    @index(0)
    public cno: string = "";
    @index(5)
    public pno: string = "";
    @index(10)
    public domain: string = "";
    @index(15)
    public subject: string = "";
    @index(20)
    public length: number = 0;
    @index(25)
    public expire_at: number = 0;
    @index(30)
    public mail: string = "";
    @index(35)
    public is_ca: boolean = false;
    @index(40)
    public ca_cert: string = "";
    @index(45)
    public ca_key: string = "";
    @index(50)
    public ips: string[] = [];

}