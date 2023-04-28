/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Reference {

    @index(0)
    public urn: string = "";
    @index(5)
    public namespace: string = "";
    @index(10)
    public name: string = "";
    @index(15)
    public version: string = "";
    @index(20)
    public proto: string = "";
    @index(25)
    public codec: string = "";
    @index(30)
    public flags: number = 0;
    @index(35)
    public timeout: number = 0;
    @index(40)
    public retries: number = 0;
    @index(45)
    public node: string = "";
    @index(50)
    public inst: string = "";
    @index(55)
    public zone: string = "";
    @index(60)
    public cluster: string = "";
    @index(65)
    public cell: string = "";
    @index(70)
    public group: string = "";
    @index(75)
    public address: string = "";

}